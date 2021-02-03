package game

// Message struct for messages of websockets
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
