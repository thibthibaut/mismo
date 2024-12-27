package game

type Player struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Lives        int     `json:"lives"`
	Number       *uint64 `json:"number"`
	IsHost       bool    `json:"isHost"`
	HasSubmitted bool    `json:"hasSubmitted"`
}

func NewPlayer(id, name string, isHost bool) *Player {
	return &Player{
		ID:           id,
		Name:         name,
		Lives:        7,
		IsHost:       isHost,
		Number:       nil,
		HasSubmitted: false,
	}
}
