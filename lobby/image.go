package lobby

import "log"

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
	drawing := &Drawing{
		Client: client,
		Img:    i.Img,
	}
	c.room.game.drawChan <- drawing
}
