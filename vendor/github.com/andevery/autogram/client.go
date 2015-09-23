package autogram

import (
	"github.com/andevery/instaw"
	"github.com/andevery/instax"
)

type Client struct {
	api *instax.Client
	web *instaw.Client
}

func NewClient(user, pass, token string) *Client {
	return &Client{}
}

func (c *Client) Api() *instax.Client {
	return c.api
}

func (c *Client) Web() *instaw.Client {
	return c.web
}
