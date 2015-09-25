package instax

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type UserShort struct {
	ID             string `json:"id"`
	FullName       string `json:"full_name"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

type User struct {
	Bio    string `json:"bio"`
	Counts struct {
		FollowedBy int `json:"followed_by"`
		Follows    int `json:"follows"`
		Media      int `json:"media"`
	} `json:"counts"`
	FullName       string `json:"full_name"`
	ID             string `json:"id"`
	ProfilePicture string `json:"profile_picture"`
	Username       string `json:"username"`
	Website        string `json:"website"`
}

func (c *Client) User(userID string) (*User, error) {
	path := fmt.Sprintf("/users/%s", userID)

	resp, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(resp.Data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type UsersFeed struct {
	client     *Client
	init       bool
	query      string
	endPoint   string
	nextCursor string
}

func (f *UsersFeed) do(values *url.Values) ([]UserShort, error) {
	path := fmt.Sprintf("/users/%s/%s", f.query, f.endPoint)

	resp, err := f.client.do("GET", path, values)
	if err != nil {
		return nil, err
	}

	var u []UserShort
	err = json.Unmarshal(resp.Data, &u)
	if err != nil {
		return nil, err
	}

	f.nextCursor = resp.Pagination["next_cursor"]

	return u, nil
}

func (f *UsersFeed) Next() ([]UserShort, error) {
	if f.init {
		f.init = false
		return f.do(nil)
	} else if len(f.nextCursor) == 0 {
		return nil, EOF
	}

	values := &url.Values{}
	values.Add("cursor", f.nextCursor)

	return f.do(values)
}
