package instax

import (
	"encoding/json"
	"fmt"
)

func (c *Client) Likes(mediaID string) ([]UserShort, error) {
	path := fmt.Sprintf("/media/%s/likes", mediaID)

	resp, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var likes []UserShort
	err = json.Unmarshal(resp.Data, &likes)
	if err != nil {
		return nil, err
	}

	return likes, nil
}

func (c *Client) Like(mediaID string) error {
	path := fmt.Sprintf("/media/%s/likes", mediaID)

	_, err := c.do("POST", path, nil)
	if err != nil {
		return err
	}

	return nil
}
