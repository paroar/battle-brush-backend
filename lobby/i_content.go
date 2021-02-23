package lobby

// IContent interface of Content in Message
type IContent interface {
	Do(c *Client)
}
