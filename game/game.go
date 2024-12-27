package game

import (
	"errors"
	"sync"
)

type GameState string

const (
	Waiting  GameState = "waiting"
	Playing  GameState = "playing"
	Finished GameState = "finished"
)

type Game struct {
	ID      string             `json:"id"`
	Players map[string]*Player `json:"players"`
	State   GameState          `json:"state"`
	mu      sync.Mutex
}

func NewGame(id string) *Game {
	return &Game{
		ID:      id,
		Players: make(map[string]*Player),
		State:   Waiting,
	}
}

func (g *Game) AddPlayer(player *Player) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.State != Waiting {
		return errors.New("game has already started")
	}

	g.Players[player.ID] = player
	return nil
}

func (g *Game) SubmitNumber(playerID string, number uint64) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	player, exists := g.Players[playerID]
	if !exists {
		return errors.New("player not found")
	}

	if player.Lives <= 0 {
		return errors.New("player is eliminated")
	}

	num := number
	player.Number = &num
	player.HasSubmitted = true
	return nil
}

func (g *Game) AllPlayersSubmitted() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, p := range g.Players {
		if p.Lives > 0 && !p.HasSubmitted {
			return false
		}
	}
	return true
}

func (g *Game) EvaluateRound() {
	g.mu.Lock()
	defer g.mu.Unlock()

	var min, max uint64
	var minPlayer, maxPlayer *Player
	numberCount := make(map[uint64][]string)

	// First pass: find duplicates and collect numbers
	for _, p := range g.Players {
		if p.Lives <= 0 || p.Number == nil {
			continue
		}
		numberCount[*p.Number] = append(numberCount[*p.Number], p.ID)
	}

	// Handle duplicates (mismo)
	for num, players := range numberCount {
		if len(players) > 1 {
			for _, playerID := range players {
				g.Players[playerID].Lives--
			}
			continue
		}
		println(num)

		// Track min and max for single numbers
		if len(players) == 1 {
			player := g.Players[players[0]]
			if minPlayer == nil || *player.Number < min {
				min = *player.Number
				minPlayer = player
			}
			if maxPlayer == nil || *player.Number > max {
				max = *player.Number
				maxPlayer = player
			}
		}
	}

	// Reduce lives for min and max players
	if minPlayer != nil {
		minPlayer.Lives--
	}
	if maxPlayer != nil {
		maxPlayer.Lives--
	}

	// Reset for next round
	for _, p := range g.Players {
		p.Number = nil
		p.HasSubmitted = false
	}

	// Check if game is finished
	activePlayers := 0
	for _, p := range g.Players {
		if p.Lives > 0 {
			activePlayers++
		}
	}
	if activePlayers <= 1 {
		g.State = Finished
	}
}
