// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Player represents a game participant.
type Player struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Lives     int             `json:"lives"`
	Number    int             `json:"number"`
	HasPlayed bool            `json:"hasPlayed"`
	IsHost    bool            `json:"isHost"`
	Conn      *websocket.Conn `json:"-"`
	mu        sync.Mutex      `json:"-"`
}

// Game represents an instance of the game.
type Game struct {
	ID      string             `json:"id"`
	Players map[string]*Player `json:"players"`
	State   string             `json:"state"` // "waiting", "playing", "roundEnd", "gameOver"
	Round   int                `json:"round"`
	mu      sync.Mutex         `json:"-"`
}

var (
	games    = make(map[string]*Game)
	gamesMu  sync.Mutex
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// createGame initializes a new game with a unique ID.
func createGame() *Game {
	return &Game{
		ID:      uuid.New().String()[:6],
		Players: make(map[string]*Player),
		State:   "waiting",
		Round:   1,
	}
}

// broadcast sends the current game state to all connected players.
func (g *Game) broadcast() {
	g.mu.Lock()
	state := map[string]interface{}{
		"id":      g.ID,
		"state":   g.State,
		"round":   g.Round,
		"players": g.Players,
	}
	players := make([]*Player, 0, len(g.Players))
	for _, p := range g.Players {
		players = append(players, p)
	}
	g.mu.Unlock()

	for _, p := range players {
		p.mu.Lock()
		if p.Conn != nil {
			err := p.Conn.WriteJSON(state)
			if err != nil {
				log.Printf("Error broadcasting to player %s: %v", p.ID, err)
				p.Conn.Close()
				p.Conn = nil
			}
		}
		p.mu.Unlock()
	}
}

// addPlayer adds a new player to the game.
func (g *Game) addPlayer(name string, conn *websocket.Conn) *Player {
	g.mu.Lock()
	defer g.mu.Unlock()

	isHost := len(g.Players) == 0
	player := &Player{
		ID:     uuid.New().String(),
		Name:   name,
		Lives:  7,
		Conn:   conn,
		IsHost: isHost,
	}
	g.Players[player.ID] = player
	return player
}

// removePlayer removes a player from the game.
func (g *Game) removePlayer(playerID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.Players, playerID)
}

// submitNumber processes a player's number submission.
func (g *Game) submitNumber(playerID string, number int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.State != "playing" {
		return
	}

	if player, exists := g.Players[playerID]; exists && !player.HasPlayed && player.Lives > 0 {
		player.Number = number
		player.HasPlayed = true

		// Check if all active players have submitted their numbers.
		allPlayed := true
		for _, p := range g.Players {
			if p.Lives > 0 && !p.HasPlayed {
				allPlayed = false
				break
			}
		}

		if allPlayed {
			g.resolveRound()
			g.State = "roundEnd"
		}
	}
}

// resolveRound applies game rules after all players have submitted their numbers.
func (g *Game) resolveRound() {
	// Collect all submitted numbers.
	numbers := make([]int, 0)
	counts := make(map[int]int)
	playerNumbers := make(map[string]int) // playerID -> number

	for _, p := range g.Players {
		if p.Lives > 0 && p.HasPlayed {
			numbers = append(numbers, p.Number)
			counts[p.Number]++
			playerNumbers[p.ID] = p.Number
		}
	}

	if len(numbers) == 0 {
		return
	}

	// Determine min and max.
	min := numbers[0]
	max := numbers[0]
	for _, num := range numbers {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}

	// Identify players with min and max numbers.
	minPlayers := []string{}
	maxPlayers := []string{}
	for pid, num := range playerNumbers {
		if num == min {
			minPlayers = append(minPlayers, pid)
		}
		if num == max {
			maxPlayers = append(maxPlayers, pid)
		}
	}

	// Reduce lives for players with min and max numbers.
	for _, pid := range minPlayers {
		if player, exists := g.Players[pid]; exists && player.Lives > 0 {
			player.Lives--
			if player.Lives <= 0 {
				// Player is eliminated.
				player.Lives = 0
			}
		}
	}

	for _, pid := range maxPlayers {
		if player, exists := g.Players[pid]; exists && player.Lives > 0 {
			player.Lives--
			if player.Lives <= 0 {
				// Player is eliminated.
				player.Lives = 0
			}
		}
	}

	// Handle "Mismo" situations: exactly two players with the same number.
	for num, cnt := range counts {
		if cnt == 2 {
			for pid, pnum := range playerNumbers {
				if pnum == num {
					if player, exists := g.Players[pid]; exists && player.Lives > 0 {
						player.Lives = 0 // Eliminated
					}
				}
			}
		}
	}

	// Check if the game is over.
	activePlayers := 0
	for _, p := range g.Players {
		if p.Lives > 0 {
			activePlayers++
		}
	}
	if activePlayers <= 1 {
		g.State = "gameOver"
	}
}

// startNextRound prepares the game for the next round.
func (g *Game) startNextRound() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.State != "roundEnd" && g.State != "gameOver" {
		return
	}

	// Reset players' submission status.
	for _, p := range g.Players {
		if p.Lives > 0 {
			p.HasPlayed = false
			p.Number = 0
		}
	}

	if g.State != "gameOver" {
		g.Round++
		g.State = "playing"
	}
}

// createGameHandler handles the creation of a new game.
func createGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		return
	}

	game := createGame()

	gamesMu.Lock()
	games[game.ID] = game
	gamesMu.Unlock()

	response := map[string]string{
		"gameId": game.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// wsHandler manages WebSocket connections for a specific game.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract game ID from the URL.
	gameID := r.URL.Path[len("/ws/game/"):]

	gamesMu.Lock()
	game, exists := games[gameID]
	gamesMu.Unlock()

	if !exists {
		http.Error(w, "Game not found.", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade Error: %v", err)
		return
	}
	defer conn.Close()

	var player *Player

	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("WebSocket Read Error: %v", err)
			if player != nil {
				game.removePlayer(player.ID)
				game.broadcast()
			}
			break
		}

		switch msg["type"] {
		case "join":
			name, ok := msg["name"].(string)
			if !ok || name == "" {
				conn.WriteJSON(map[string]string{"error": "Invalid name."})
				continue
			}
			player = game.addPlayer(name, conn)
			game.broadcast()

		case "start":
			if player == nil || !player.IsHost {
				conn.WriteJSON(map[string]string{"error": "Only host can start the game."})
				continue
			}
			game.mu.Lock()
			if game.State != "waiting" {
				game.mu.Unlock()
				conn.WriteJSON(map[string]string{"error": "Game already started."})
				continue
			}
			if len(game.Players) < 3 {
				game.mu.Unlock()
				conn.WriteJSON(map[string]string{"error": "At least 3 players required to start."})
				continue
			}
			game.State = "playing"
			game.mu.Unlock()
			game.broadcast()

		case "number":
			if player == nil || game.State != "playing" || player.Lives <= 0 {
				conn.WriteJSON(map[string]string{"error": "Cannot submit number at this time."})
				continue
			}
			// Parse the number.
			numberFloat, ok := msg["number"].(float64)
			if !ok {
				conn.WriteJSON(map[string]string{"error": "Invalid number."})
				continue
			}
			number := int(numberFloat)
			if number < 0 {
				conn.WriteJSON(map[string]string{"error": "Number cannot be negative."})
				continue
			}
			game.submitNumber(player.ID, number)
			game.broadcast()

		case "nextRound":
			if player == nil || !player.IsHost {
				conn.WriteJSON(map[string]string{"error": "Only host can start the next round."})
				continue
			}
			game.startNextRound()
			game.broadcast()

		default:
			conn.WriteJSON(map[string]string{"error": "Unknown message type."})
		}
	}
}

func main() {
	// Serve static files from the "static" directory.
	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	// API Endpoints
	http.HandleFunc("/create-game", createGameHandler)
	http.HandleFunc("/ws/game/", wsHandler)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
