package instaw

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andevery/instax"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	Scheme    = "https"
	Host      = "www.instagram.com"
	Namespace = "web"
)

var (
	TooManyRequests = errors.New("Sorry, too many requests. Please try again later.")
	BadRequest      = errors.New("Bad request")
)

type DelayerFunc func() time.Duration

type Client struct {
	Delayer           DelayerFunc
	UserAgent         string
	Cookie            string
	CSRFToken         string
	RateLimitDelayMin time.Duration
	RateLimitDelayMax time.Duration

	username string
	password string

	client *http.Client
}

func NewClient(login, password string) (*Client, error) {
	c := &Client{
		UserAgent:         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
		RateLimitDelayMin: 40 * time.Second,
		RateLimitDelayMax: 80 * time.Second,
		username:          login,
		password:          password,
		client:            new(http.Client),
	}
	err := c.signIn()
	if err != nil {
		return nil, err
	}
	return c, nil
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

	switch resp.StatusCode {
	case 403:
		return TooManyRequests
	case 400:
		return BadRequest
	}

	var status Response
	err = json.Unmarshal(body, &status)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) signIn() error {
	var cookies []string
	resp, err := http.Get("https://instagram.com/accounts/login/")
	if err != nil {
		return err
	}

	rc := resp.Cookies()
	for _, v := range rc {
		cookies = append(cookies, fmt.Sprintf("%s=%s;", v.Name, v.Value))
		if v.Name == "csrftoken" {
			c.CSRFToken = v.Value
		}
	}
	c.Cookie = strings.Join(cookies, " ")

	u := &url.URL{
		Scheme: Scheme,
		Host:   Host,
		Path:   "/accounts/login/ajax/",
	}

	values := &url.Values{}
	values.Add("username", c.username)
	values.Add("password", c.password)

	req, err := http.NewRequest("POST", u.String(), bytes.NewBufferString(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Cookie", c.Cookie)
	req.Header.Add("Referer", fmt.Sprintf("%s://%s/accounts/login", Scheme, Host))
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("X-CSRFToken", c.CSRFToken)
	req.Header.Add("X-Instagram-AJAX", "1")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err = c.client.Do(req)
	if err != nil {
		return err
	}

	rc = resp.Cookies()
	cookies = []string{}
	for _, v := range rc {
		cookies = append(cookies, fmt.Sprintf("%s=%s;", v.Name, v.Value))
		if v.Name == "csrftoken" {
			c.CSRFToken = v.Value
		}
	}
	c.Cookie = strings.Join(cookies, " ")

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

func (c *Client) Unfollow(user *instax.User) error {
	path := fmt.Sprintf("friendships/%s/unfollow/", user.ID)
	referer := fmt.Sprintf("https://instagram.com/%s/", user.Username)

	return c.do("POST", path, referer)
}
