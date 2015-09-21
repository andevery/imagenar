package instax

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Media struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Link         string    `json:"link"`
	User         UserShort `json:"user"`
	Location     Location  `json:"location"`
	Tags         []string  `json:"tags"`
	Filter       string    `json:"filter"`
	UserHasLiked bool      `json:"user_has_liked"`
	CreatedTime  string    `json:"created_time"`

	Caption struct {
		ID          string    `json:"id"`
		Text        string    `json:"text"`
		CreatedTime string    `json:"created_time"`
		From        UserShort `json:"from"`
	} `json:"caption"`

	Images struct {
		LowResolution      MediaItem `json:"low_resolution"`
		StandardResolution MediaItem `json:"standard_resolution"`
		Thumbnail          MediaItem `json:"thumbnail"`
	} `json:"images"`

	Videos struct {
		LowBandwidth       MediaItem `json:"low_bandwidth"`
		LowResolution      MediaItem `json:"low_resolution"`
		StandardResolution MediaItem `json:"standard_resolution"`
	} `json:"videos"`

	Comments struct {
		Count int       `json:"count"`
		Data  []Comment `json:"data"`
	} `json:"comments"`

	Likes struct {
		Count int         `json:"count"`
		Data  []UserShort `json:"data"`
	} `json:"likes"`

	UsersInPhoto []struct {
		Position struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"position"`
		User UserShort `json:"user"`
	} `json:"users_in_photo"`
}

type Comment struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	CreatedTime string    `json:"created_time"`
	From        UserShort `json:"from"`
}

type MediaItem struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type Location struct {
	ID        int     `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
}

func (c *Client) Media(id string) (*Media, error) {
	path := fmt.Sprintf("/media/%s", id)

	resp, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var m Media
	err = json.Unmarshal(resp.Data, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (c *Client) RecentMediaByUser(userID string) ([]Media, error) {
	path := fmt.Sprintf("/users/%s/media/recent", userID)

	resp, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var m []Media
	err = json.Unmarshal(resp.Data, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

type MediaFeed struct {
	client          *Client
	endPoint        string
	query           string
	paramMaxID      string
	paginationMaxID string
	nextMaxID       string
}

func (f *MediaFeed) CanNext() bool {
	if len(f.nextMaxID) > 0 {
		return true
	}

	return false
}

func (f *MediaFeed) do(values *url.Values) ([]Media, error) {
	path := fmt.Sprintf("/%s/%s/media/recent", f.endPoint, f.query)

	resp, err := f.client.do("GET", path, values)
	if err != nil {
		return nil, err
	}

	var m []Media
	err = json.Unmarshal(resp.Data, &m)
	if err != nil {
		return nil, err
	}

	f.nextMaxID = resp.Pagination[f.paginationMaxID]

	return m, nil
}

func (f *MediaFeed) Get() ([]Media, error) {
	return f.do(nil)
}

func (f *MediaFeed) Next() ([]Media, error) {
	values := &url.Values{}
	values.Add(f.paramMaxID, f.nextMaxID)

	return f.do(values)
}

func (c *Client) MediaByTag(tag string) *MediaFeed {
	return &MediaFeed{
		endPoint:        "tags",
		paramMaxID:      "max_tag_id",
		paginationMaxID: "next_max_tag_id",
		query:           tag,
		client:          c,
	}
}

func (c *Client) MediaByUser(userID string) *MediaFeed {
	return &MediaFeed{
		endPoint:        "users",
		paramMaxID:      "max_id",
		paginationMaxID: "next_max_id",
		query:           userID,
		client:          c,
	}
}
