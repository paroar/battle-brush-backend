package model

import (
	"encoding/json"
)

// Drawing struct
type Drawing struct {
	PlayerID string `redis:"playerid"`
	RoomID   string `redis:"roomid"`
	Img      string `redis:"img"`
}

// NewDrawing constructor
func NewDrawing(playerid, roomid, img string) *Drawing {
	return &Drawing{
		PlayerID: playerid,
		RoomID:   roomid,
		Img:      img,
	}
}

// MarshalBinary marshaler
func (d *Drawing) MarshalBinary() (data []byte, err error) {
	return json.Marshal(d)
}

// UnmarshalBinary unmarshaler
func (d *Drawing) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	return nil
}
