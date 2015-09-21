package instax

import (
	"encoding/json"
	"fmt"
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
