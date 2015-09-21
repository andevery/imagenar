package instax

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// "encoding/json"

// "net/http"

const (
	apiScheme  = "https"
	apiHost    = "api.instagram.com"
	apiVersion = "v1"
)

var (
	RateLimitException = errors.New("The maximum number of requests per hour has been exceeded.")
	NotFound           = errors.New("API not found error")
	NotAllowed         = errors.New("API not allowed error")
)

type DelayerFunc func() time.Duration

type Client struct {
	Delayer            DelayerFunc
	rateLimitRemaining uint32
	httpClient         *http.Client
	accessToken        string
	secret             string
}

func NewClient(accessToken string) *Client {
	return &Client{
		rateLimitRemaining: 5000,
		accessToken:        accessToken,
		httpClient:         new(http.Client),
	}
}

func (c *Client) Limit() uint32 {
	return atomic.LoadUint32(&c.rateLimitRemaining)
}

func generateSig(endpoint, secret string, values *url.Values) string {
	var keys []string
	for k := range *values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	params := []string{endpoint}
	for _, k := range keys {
		params = append(params, fmt.Sprintf("%s=%s", k, values.Get(k)))
	}
	sig := strings.Join(params, "|")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(sig))
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *Client) do(method, path string, values *url.Values) (*Response, error) {
	if c.Delayer != nil {
		d := c.Delayer()
		time.Sleep(d)
	}

	u := &url.URL{
		Scheme: apiScheme,
		Host:   apiHost,
		Path:   fmt.Sprintf("/%s%s", apiVersion, path),
	}

	if values == nil {
		values = &url.Values{}
	}

	values.Add("access_token", c.accessToken)
	u.RawQuery = values.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	limit, err := strconv.ParseUint(resp.Header.Get("X-Ratelimit-Remaining"), 10, 64)
	atomic.StoreUint32(&c.rateLimitRemaining, uint32(limit))

	switch r.Meta.Code {
	case 200:
		return &r, nil
	case 429:
		return nil, RateLimitException
	default:
		switch r.Meta.ErrorType {
		case "APINotAllowedError":
			return nil, NotAllowed
		case "APINotFoundError":
			return nil, NotFound
		default:
			return nil, errors.New(r.Meta.ErrorMessage)
		}
	}
}
