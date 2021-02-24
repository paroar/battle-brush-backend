package lobby

import (
	"log"

	"github.com/paroar/battle-brush-backend/drawing"
)

// Image struct
type Image struct {
	UserID   string `json:"userid"`
	UserName string `json:"username"`
	Img      string `json:"img"`
}

//Do retrieves a Drawing
func (i *Image) Do(c *Client) {
	client, err := c.room.getClient(i.UserID)
	if err != nil {
		log.Println(err)
	}
	drawing := &drawing.Drawing{
		ClientID: client.id,
		Img:      i.Img,
	}
	c.room.game.drawChan <- drawing
}
