package instaw

import (
	"encoding/json"
	"fmt"
	"github.com/andevery/instax"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	Scheme    = "https"
	Host      = "instagram.com"
	Namespace = "web"
)

type DelayerFunc func() time.Duration

type Client struct {
	Delayer           DelayerFunc
	UserAgent         string
	Cookie            string
	CSRFToken         string
	RateLimitDelayMin time.Duration
	RateLimitDelayMax time.Duration

	client *http.Client
}

func NewClient() *CLient {
	return &WebCLient{
		UserAgent:         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
		RateLimitDelayMin: 60 * time.Second,
		RateLimitDelayMax: 180 * time.Second,
		client:            new(http.Client),
	}
}

func (c *Client) do(method, path, referer string) error {
	if c.Delayer != nil {
		time.Sleep(c.Delayer())
	}

	u := &url.URL{
		Scheme: Scheme,
		Host:   Host,
		Path:   fmt.Sprintf("/%s/%s", Namespace, path),
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Cookie", c.Cookie)
	req.Header.Add("Referer", referer)
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("X-CSRFToken", c.CSRFToken)
	req.Header.Add("X-Instagram-AJAX", "1")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var status Response
	err = json.Unmarshal(body, &status)
	if err != nil {
		fmt.Println("")
		log.Println(resp.StatusCode)
		log.Println(string(body))
		time.Sleep(time.Duration(rand.Int63n(int64(c.RateLimitDelayMax)-int64(c.RateLimitDelayMin)) + int64(c.RateLimitDelayMin)))
	}

	return nil
}

func (c *Client) Like(media *instax.Media) error {
	path := fmt.Sprintf("likes/%s/like/", media.ID)

	return c.do("POST", path, media.Link)
}

func (c *Client) Follow(user *instax.User) error {
	path := fmt.Sprintf("friendships/%s/follow/", user.ID)
	referer := fmt.Sprintf("https://instagram.com/%s/", user.Username)

	return c.do("POST", path, referer)
}
