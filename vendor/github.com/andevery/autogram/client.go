package autogram

import (
	"github.com/andevery/instaw"
	"github.com/andevery/instax"
	"log"
	"math/rand"
	"time"
)

type Client struct {
	WebRate struct {
		HourLimit int
		DayLimit  int
		MaxDelay  time.Duration
		MinDelay  time.Duration
	}
	ApiRate struct {
		MinRemaining uint32
	}

	webCounters struct {
		day  int
		hour int
	}
	tickers struct {
		webDay  *time.Ticker
		webHour *time.Ticker
		apiHour *time.Ticker
	}

	webAllow   chan bool
	clientDone chan bool

	api *instax.Client
	web *instaw.Client
}

func DefaultClient(user, pass, token string) (*Client, error) {
	c := new(Client)
	var err error
	c.web, err = instaw.NewClient(user, pass)
	if err != nil {
		return nil, err
	}
	c.api = instax.NewClient(token)

	c.WebRate.HourLimit = 180
	c.WebRate.DayLimit = 0
	c.WebRate.MaxDelay = 35 * time.Second
	c.WebRate.MinDelay = 25 * time.Second
	c.ApiRate.MinRemaining = 500
	c.tickers.webDay = time.NewTicker(24 * time.Hour)
	c.tickers.webHour = time.NewTicker(time.Hour)
	c.webAllow = make(chan bool)
	c.clientDone = make(chan bool)

	c.start()

	return c, nil
}

func (c *Client) Api() *instax.Client {
	if c.tickers.apiHour == nil {
		c.tickers.apiHour = time.NewTicker(time.Hour)
	}
	if c.api.Limit() < c.ApiRate.MinRemaining {
		<-c.tickers.apiHour.C
	}
	return c.api
}

func (c *Client) Web() *instaw.Client {
	<-c.webAllow
	return c.web
}

func (c *Client) start() {
	go func() {
		for {
			select {
			case <-c.tickers.webDay.C:
				c.webCounters.day = 0
			case <-c.tickers.webHour.C:
				c.webCounters.hour = 0
			case <-c.clientDone:
				c.tickers.apiHour.Stop()
				c.tickers.webHour.Stop()
				c.tickers.webDay.Stop()
				close(c.webAllow)
				close(c.clientDone)
				return
			default:
				if c.allowed() {
					c.webAllow <- true
					c.webCounters.day++
					c.webCounters.hour++
				}
				delay := time.Duration(rand.Int63n(int64(c.WebRate.MaxDelay)-int64(c.WebRate.MinDelay)) + int64(c.WebRate.MinDelay))
				log.Printf("Delay: %v\n", delay)
				time.Sleep(delay)
			}
		}
	}()
}

func (c *Client) Stop() {
	c.clientDone <- true
}

func (c *Client) allowed() bool {
	allow := true
	if c.WebRate.DayLimit > 0 {
		allow = allow && c.webCounters.day < c.WebRate.DayLimit
	}

	if c.WebRate.HourLimit > 0 {
		allow = allow && c.webCounters.hour < c.WebRate.HourLimit
	}
	return allow
}

func (c *Client) LikeAFew(media []instax.Media, count int) {
	for _, i := range randomIndexes(len(media), count) {
		c.Web().Like(&media[i])
	}
}
