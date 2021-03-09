package message

// Envelope struct for websockets messages
type Envelope struct {
	Type    int         `json:"type"`
	Content interface{} `json:"content"`
}
