package lobby

import "github.com/paroar/battle-brush-backend/drawing"

// Image struct
type Image struct {
	UserID   string `json:"userid"`
	UserName string `json:"username"`
	Img      string `json:"img"`
}

//Do retrieves a Drawing
func (i *Image) Do(c *Client) {
	drawing := &drawing.Drawing{
		ClientID: i.UserID,
		Img:      i.Img,
	}
	game := c.room.GetGame().(IGame)
	game.Drawing(drawing)
}
