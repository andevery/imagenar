package instax

import (
	"encoding/json"
)

// "net/http"

const (
	apiUrl = "https://api.instagram.com/v1"
)

type Client struct {
	authToken string
}

func NewClient(authToken string) *Client {
	return &Client{authToken: authToken}
}

func (c *Client) send() (*json.RawMessage, error) {
	return nil, nil
}
