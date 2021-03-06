package model

// Vote struct
type Vote struct {
	PlayerID string
	Vote     float64
}

// NewVote constructor
func NewVote(playerid string, vote float64) *Vote {
	return &Vote{
		PlayerID: playerid,
		Vote:     vote,
	}
}
