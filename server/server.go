package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Player struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Lives          int    `json:"lives"`
	Number         uint64 `json:"number"`
	Submitted      bool   `json:"submitted"`
	Eliminated     bool   `json:"eliminated"`
	IsHost         bool   `json:"isHost"`
	SubmissionTime time.Time
}

type Game struct {
	ID             string            `json:"id"`
	Players        map[string]*Player `json:"players"`
	HasStarted     bool              `json:"hasStarted"`
	RoundCompleted bool              `json:"roundCompleted"`
	Mutex          sync.Mutex        `json:"-"`
}

var (
	games   = make(map[string]*Game)
	gamesMu sync.Mutex
)

func main() {

    // Serve static files from "./static" directory.
    http.Handle("/", http.FileServer(http.Dir("./static")))

    http.HandleFunc("/createGame", handleCreateGame)
    http.HandleFunc("/joinGame", handleJoinGame)
    http.HandleFunc("/submitNumber", handleSubmitNumber)
    http.HandleFunc("/nextRound", handleNextRound)
    http.HandleFunc("/gameState", handleGameState)

    log.Println("Starting server on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// serveIndex serves our embedded HTML/JS/CSS
func serveIndex(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}

// handleCreateGame handles the creation of a new game and returns the game ID
// Expecting a POST with form data: hostName
func handleCreateGame(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		return
	}
	hostName := strings.TrimSpace(r.FormValue("hostName"))
	if hostName == "" {
		http.Error(w, "Host name cannot be empty", http.StatusBadRequest)
		return
	}

	gameID := generateGameID()
	playerID := generatePlayerID()

	g := &Game{
		ID:      gameID,
		Players: make(map[string]*Player),
	}
	g.Players[playerID] = &Player{
		ID:         playerID,
		Name:       hostName,
		Lives:      7,
		IsHost:     true,
		Submitted:  false,
		Eliminated: false,
	}

	gamesMu.Lock()
	games[gameID] = g
	gamesMu.Unlock()

	resp := map[string]string{
		"gameID":   gameID,
		"playerID": playerID,
	}
	writeJSON(w, resp)
}

// handleJoinGame handles joining a game
// Expecting a POST with form data: gameID, playerName
func handleJoinGame(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		return
	}
	gameID := strings.TrimSpace(r.FormValue("gameID"))
	playerName := strings.TrimSpace(r.FormValue("playerName"))

	gamesMu.Lock()
	g, ok := games[gameID]
	gamesMu.Unlock()
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	if g.HasStarted {
		http.Error(w, "Game already started", http.StatusForbidden)
		return
	}

	playerID := generatePlayerID()
	g.Players[playerID] = &Player{
		ID:         playerID,
		Name:       playerName,
		Lives:      7,
		Submitted:  false,
		Eliminated: false,
	}

	resp := map[string]string{
		"gameID":   gameID,
		"playerID": playerID,
	}
	writeJSON(w, resp)
}

// handleSubmitNumber handles a player's number submission
// Expecting a POST with form data: gameID, playerID, number
func handleSubmitNumber(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		return
	}
	gameID := r.FormValue("gameID")
	playerID := r.FormValue("playerID")
	numberStr := r.FormValue("number")

	num, err := strconv.ParseUint(numberStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	gamesMu.Lock()
	g, ok := games[gameID]
	gamesMu.Unlock()
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	p, ok := g.Players[playerID]
	if !ok {
		http.Error(w, "Player not found in this game", http.StatusNotFound)
		return
	}
	if p.Eliminated {
		http.Error(w, "Player is eliminated", http.StatusForbidden)
		return
	}

	// If the game hasn't started, check if host can start it
	if !g.HasStarted {
		// Only the host can start the game if enough players
		if !p.IsHost {
			http.Error(w, "Only host can start the game", http.StatusForbidden)
			return
		}
		if len(g.Players) < 3 {
			http.Error(w, "Minimum 3 players required to start the game", http.StatusForbidden)
			return
		}
		g.HasStarted = true
	}

	// Mark the player's submission
	p.Submitted = true
	p.Number = num
	p.SubmissionTime = time.Now()

	// Check if all have submitted => do the round resolution
	allSubmitted := true
	for _, pl := range g.Players {
		if !pl.Eliminated && !pl.Submitted {
			allSubmitted = false
			break
		}
	}
	if allSubmitted {
		resolveRound(g)
		g.RoundCompleted = true
	}

	resp := map[string]string{"status": "ok"}
	writeJSON(w, resp)
}

// handleNextRound can only be triggered by the host, after a round is resolved
// Expecting a POST with form data: gameID, hostID
func handleNextRound(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		return
	}
	gameID := r.FormValue("gameID")
	hostID := r.FormValue("hostID")

	gamesMu.Lock()
	g, ok := games[gameID]
	gamesMu.Unlock()
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	host, ok := g.Players[hostID]
	if !ok || !host.IsHost {
		http.Error(w, "Only the host can start the next round", http.StatusForbidden)
		return
	}
	if !g.RoundCompleted {
		http.Error(w, "Round not completed yet", http.StatusForbidden)
		return
	}

	// Check if game already ended
	if checkGameOver(g) {
		http.Error(w, "Game is already over", http.StatusForbidden)
		return
	}

	// Reset submission states for the next round
	for _, pl := range g.Players {
		if pl.Eliminated {
			continue
		}
		pl.Submitted = false
		pl.Number = 0
	}
	g.RoundCompleted = false

	resp := map[string]string{"status": "ok"}
	writeJSON(w, resp)
}

// handleGameState returns the current state of the game as JSON
// Expecting: GET or POST with either query or form: gameID, playerID
func handleGameState(w http.ResponseWriter, r *http.Request) {
	gameID := r.FormValue("gameID")
	playerID := r.FormValue("playerID")

	gamesMu.Lock()
	g, ok := games[gameID]
	gamesMu.Unlock()
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	state := struct {
		GameID         string            `json:"gameID"`
		HasStarted     bool              `json:"hasStarted"`
		RoundCompleted bool              `json:"roundCompleted"`
		Players        map[string]*Player `json:"players"`
		IsHost         bool              `json:"isHost"`
		GameOver       bool              `json:"gameOver"`
		Winner         string            `json:"winner"`
	}{
		GameID:         g.ID,
		HasStarted:     g.HasStarted,
		RoundCompleted: g.RoundCompleted,
		Players:        g.Players,
		IsHost:         false,
		GameOver:       false,
		Winner:         "",
	}

	if p, found := g.Players[playerID]; found && p.IsHost {
		state.IsHost = true
	}

	if checkGameOver(g) {
		state.GameOver = true
		state.Winner = getWinner(g)
	}

	writeJSON(w, state)
}

// -----------------------------------
// Game Logic
// -----------------------------------

func resolveRound(g *Game) {
	// Identify min, max
	var minNum, maxNum uint64
	var hasMin, hasMax bool

	// Keep track of all submissions in a map to check for duplicates
	submissions := make(map[uint64][]*Player)

	for _, p := range g.Players {
		if p.Eliminated {
			continue
		}
		num := p.Number
		if !hasMin || num < minNum {
			minNum = num
			hasMin = true
		}
		if !hasMax || num > maxNum {
			maxNum = num
			hasMax = true
		}
		submissions[num] = append(submissions[num], p)
	}

	// Handle min, max life loss
	for _, p := range g.Players {
		if p.Eliminated {
			continue
		}
		if p.Number == minNum {
			p.Lives--
			if p.Lives <= 0 {
				p.Eliminated = true
			}
		}
		if p.Number == maxNum {
			p.Lives--
			if p.Lives <= 0 {
				p.Eliminated = true
			}
		}
	}

	// Handle "mismo" => if a number was entered by exactly 2 players
	for num, players := range submissions {
		// If exactly 2 players have the same number => both eliminated
		if len(players) == 2 {
			// For the problem statement: "If 2 players have the same number, it's a 'mismo' situation"
			// The statement says: "they are both eliminated from the game."
			// So let's do exactly that:
			for _, p := range players {
				p.Eliminated = true
			}
		}
	}

	// Check if we still have more than one player
	// If only 1 remains, game is over
	checkGameOver(g)
}

func checkGameOver(g *Game) bool {
	activeCount := 0
	for _, p := range g.Players {
		if !p.Eliminated {
			activeCount++
		}
	}
	return activeCount <= 1
}

func getWinner(g *Game) string {
	for _, p := range g.Players {
		if !p.Eliminated {
			return p.Name
		}
	}
	return ""
}

// -----------------------------------
// Utility
// -----------------------------------

func generateGameID() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func generatePlayerID() string {
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}


<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Mismo</title>
  <style>
    body {
      font-family: 'Helvetica Neue', sans-serif;
      background: #fefefe;
      color: #333;
      text-align: center;
      padding: 20px;
      margin: 0;
    }
    button {
      background: #00C853;
      border: none;
      padding: 10px 20px;
      color: #fff;
      font-size: 1em;
      cursor: pointer;
      margin: 5px;
      border-radius: 4px;
    }
    button:hover {
      background: #00B44A;
    }
    input[type="text"], input[type="number"] {
      padding: 6px;
      font-size: 1em;
      border: 1px solid #ccc;
      border-radius: 4px;
      margin: 4px;
    }
    #gameScreen, #waitingRoom, #submitSection, #roundResults, #gameOverSection {
      display: none;
    }
    .playerList {
      margin: 10px 0;
    }
    .player {
      display: inline-block;
      margin: 5px;
      padding: 5px 10px;
      border: 1px solid #ccc;
      border-radius: 6px;
    }
    .eliminated {
      background-color: #ffebee;
      color: #b71c1c;
      text-decoration: line-through;
    }
    .lives {
      font-weight: bold;
    }
  </style>
</head>
<body>
  <h1 style="color:#3f51b5;">Mismo</h1>
  <div id="landingPage">
    <button id="startGameBtn">Start a Game</button>
    <button id="joinGameBtn">Join a Game</button>
  </div>
  <div id="startGameForm">
    <h2>Start a New Game</h2>
    <input type="text" id="hostName" placeholder="Your Name" />
    <button id="createGameBtn">Create Game</button>
  </div>
  <div id="joinGameForm">
    <h2>Join a Game</h2>
    <input type="text" id="joinGameID" placeholder="Game ID" />
    <input type="text" id="joinName" placeholder="Your Name" />
    <button id="joinBtn">Join</button>
  </div>
  <div id="waitingRoom">
    <h2>Waiting Room</h2>
    <p>Game ID: <span id="displayGameID"></span></p>
    <div class="playerList" id="waitingPlayers"></div>
    <p id="waitingMessage">Waiting for at least 3 players...</p>
    <button id="startGameNowBtn" style="background:#ff9800;">Start Game</button>
  </div>
  <div id="gameScreen">
    <h2>Game ID: <span id="gameIDSpan"></span></h2>
    <h3>Players</h3>
    <div class="playerList" id="playerList"></div>
    <div id="submitSection">
      <h3>Submit Your Number</h3>
      <input type="number" id="numberInput" />
      <button id="submitNumberBtn">Submit</button>
    </div>
    <div id="roundResults">
      <h3>Round Results</h3>
      <p id="resultsText"></p>
    </div>
    <button id="nextRoundBtn" style="background:#673ab7;">Next Round</button>
  </div>
  <div id="gameOverSection">
    <h2>Game Over!</h2>
    <p>Winner: <span id="winnerName"></span></p>
  </div>

  <script>
    // Simple state
    let currentGameID = null;
    let currentPlayerID = null;
    let isHost = false;
    let gameOver = false;

    const landingPage = document.getElementById('landingPage');
    const startGameForm = document.getElementById('startGameForm');
    const joinGameForm = document.getElementById('joinGameForm');
    const waitingRoom = document.getElementById('waitingRoom');
    const gameScreen = document.getElementById('gameScreen');
    const gameOverSection = document.getElementById('gameOverSection');

    const hostNameInput = document.getElementById('hostName');
    const joinGameIDInput = document.getElementById('joinGameID');
    const joinNameInput = document.getElementById('joinName');

    const displayGameID = document.getElementById('displayGameID');
    const waitingPlayers = document.getElementById('waitingPlayers');
    const waitingMessage = document.getElementById('waitingMessage');
    const startGameNowBtn = document.getElementById('startGameNowBtn');

    const gameIDSpan = document.getElementById('gameIDSpan');
    const playerList = document.getElementById('playerList');
    const submitSection = document.getElementById('submitSection');
    const roundResults = document.getElementById('roundResults');
    const resultsText = document.getElementById('resultsText');
    const nextRoundBtn = document.getElementById('nextRoundBtn');

    const winnerName = document.getElementById('winnerName');

    // Landing page button handlers
    document.getElementById('startGameBtn').addEventListener('click', () => {
      landingPage.style.display = 'none';
      startGameForm.style.display = 'block';
    });
    document.getElementById('joinGameBtn').addEventListener('click', () => {
      landingPage.style.display = 'none';
      joinGameForm.style.display = 'block';
    });

    // Create a new game
    document.getElementById('createGameBtn').addEventListener('click', async () => {
      const hostName = hostNameInput.value.trim();
      if(!hostName) {
        alert('Please enter your name');
        return;
      }
      try {
        const formData = new FormData();
        formData.append('hostName', hostName);
        const resp = await fetch('/createGame', {
          method: 'POST',
          body: formData
        });
        if(!resp.ok) {
          const t = await resp.text();
          alert("Error creating game: " + t);
          return;
        }
        const data = await resp.json();
        currentGameID = data.gameID;
        currentPlayerID = data.playerID;
        isHost = true;
        startGameForm.style.display = 'none';
        waitingRoom.style.display = 'block';
        displayGameID.textContent = currentGameID;
        pollGameState();
      } catch(err) {
        alert("Error: " + err);
      }
    });

    // Join an existing game
    document.getElementById('joinBtn').addEventListener('click', async () => {
      const gameID = joinGameIDInput.value.trim();
      const playerName = joinNameInput.value.trim();
      if(!gameID || !playerName) {
        alert('Please enter both game ID and your name');
        return;
      }
      try {
        const formData = new FormData();
        formData.append('gameID', gameID);
        formData.append('playerName', playerName);
        const resp = await fetch('/joinGame', {
          method: 'POST',
          body: formData
        });
        if(!resp.ok) {
          const t = await resp.text();
          alert("Error joining game: " + t);
          return;
        }
        const data = await resp.json();
        currentGameID = data.gameID;
        currentPlayerID = data.playerID;
        isHost = false;
        joinGameForm.style.display = 'none';
        waitingRoom.style.display = 'block';
        displayGameID.textContent = currentGameID;
        pollGameState();
      } catch(err) {
        alert("Error: " + err);
      }
    });

    // Host: start the game by submitting a dummy number (this triggers game start if we have 3+)
    startGameNowBtn.addEventListener('click', async () => {
      if(!currentGameID || !currentPlayerID) return;
      try {
        // Submit a number to forcibly start the game
        const formData = new FormData();
        formData.append('gameID', currentGameID);
        formData.append('playerID', currentPlayerID);
        formData.append('number', '0');
        const resp = await fetch('/submitNumber', {
          method: 'POST',
          body: formData
        });
        if(!resp.ok) {
          const t = await resp.text();
          alert("Error starting game: " + t);
          return;
        }
        // If success, waitingRoom will remain but we expect the game state to say "HasStarted=true"
      } catch(err) {
        alert("Error: " + err);
      }
    });

    // Submit a number
    document.getElementById('submitNumberBtn').addEventListener('click', async () => {
      const numInput = document.getElementById('numberInput');
      const val = numInput.value.trim();
      if(val === '') {
        alert('Please enter a number');
        return;
      }
      try {
        const formData = new FormData();
        formData.append('gameID', currentGameID);
        formData.append('playerID', currentPlayerID);
        formData.append('number', val);
        const resp = await fetch('/submitNumber', {
          method: 'POST',
          body: formData
        });
        if(!resp.ok) {
          const t = await resp.text();
          alert("Error: " + t);
          return;
        }
        roundResults.style.display = 'none';
        submitSection.style.display = 'none';
      } catch(err) {
        alert("Error: " + err);
      }
    });

    // Next round
    nextRoundBtn.addEventListener('click', async () => {
      if(!currentGameID || !currentPlayerID) return;
      try {
        const formData = new FormData();
        formData.append('gameID', currentGameID);
        formData.append('hostID', currentPlayerID);
        const resp = await fetch('/nextRound', {
          method: 'POST',
          body: formData
        });
        if(!resp.ok) {
          const t = await resp.text();
          alert("Error: " + t);
          return;
        }
        roundResults.style.display = 'none';
        submitSection.style.display = 'block';
        pollGameState();
      } catch(err) {
        alert("Error: " + err);
      }
    });

    // Poll game state
    async function pollGameState() {
      if(!currentGameID || !currentPlayerID || gameOver) return;
      try {
        const formData = new FormData();
        formData.append('gameID', currentGameID);
        formData.append('playerID', currentPlayerID);
        const resp = await fetch('/gameState', {
          method: 'POST',
          body: formData
        });
        if(!resp.ok) {
          setTimeout(pollGameState, 2000);
          return;
        }
        const data = await resp.json();
        // Update UI based on game state
        renderGameState(data);
        if(!gameOver) {
          setTimeout(pollGameState, 2000); // poll again
        }
      } catch(err) {
        console.log("poll error:", err);
        setTimeout(pollGameState, 2000);
      }
    }

    function renderGameState(state) {
      const { gameID, hasStarted, roundCompleted, players, isHost: hostFlag, gameOver: over, winner } = state;
      isHost = hostFlag;
      currentGameID = gameID;

      // If game over, show final
      if(over) {
        gameOver = true;
        waitingRoom.style.display = 'none';
        gameScreen.style.display = 'none';
        gameOverSection.style.display = 'block';
        winnerName.textContent = winner || 'No one';
        return;
      }

      // Update waiting room
      if(!hasStarted) {
        waitingRoom.style.display = 'block';
        gameScreen.style.display = 'none';
        // Show current players
        let html = '';
        for(const pid in players) {
          const pl = players[pid];
          html += `<div class="player">${pl.name} (${pl.lives} lives)</div>`;
        }
        waitingPlayers.innerHTML = html;
        // If host and 3+ players, show a start game button
        if(isHost && Object.keys(players).length >= 3) {
          startGameNowBtn.style.display = 'inline-block';
        } else {
          startGameNowBtn.style.display = 'none';
        }
        return;
      }

      // If hasStarted => game screen
      waitingRoom.style.display = 'none';
      gameScreen.style.display = 'block';
      gameIDSpan.textContent = gameID;

      // Build player list
      let playersHTML = '';
      let allSubmitted = true;
      for(const pid in players) {
        const pl = players[pid];
        let classes = "player";
        if(pl.eliminated) {
          classes += " eliminated";
        }
        playersHTML += `<div class="${classes}">
            <div>${pl.name}</div>
            <div class="lives">${pl.lives} ${pl.lives === 1 ? 'life' : 'lives'}</div>
            <div>${pl.submitted ? 'Submitted' : 'Thinking'}</div>
          </div>`;
        if(!pl.eliminated && !pl.submitted) {
          allSubmitted = false;
        }
      }
      playerList.innerHTML = playersHTML;

      // Show/Hide submit section based on whether this player is eliminated
      if(players[currentPlayerID].eliminated) {
        submitSection.style.display = 'none';
      } else {
        // If not submitted, show the input
        if(!players[currentPlayerID].submitted && !roundCompleted) {
          submitSection.style.display = 'block';
        } else {
          submitSection.style.display = 'none';
        }
      }

      // If the round is completed, show results
      if(roundCompleted) {
        roundResults.style.display = 'block';
        resultsText.textContent = "Round has ended. Check changes in lives!";
      } else {
        roundResults.style.display = 'none';
      }

      // Show next round button only to host if round completed
      if(isHost && roundCompleted) {
        nextRoundBtn.style.display = 'inline-block';
      } else {
        nextRoundBtn.style.display = 'none';
      }
    }
  </script>
</body>
</html>
