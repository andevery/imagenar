package instax

import (
	"encoding/json"
	"fmt"
)

const (
	NONE         = "none"
	FOLLOWS      = "follows"
	FOLLOWED_BY  = "followed_by"
	REQUESTED    = "requested"
	REQUESTED_BY = "requested_by"
	BLOCKED      = "blocked_by_you"
)

type Relationship struct {
	IncomingStatus      string `json:"incoming_status"`
	OutgoingStatus      string `json:"outgoing_status"`
	TargetUserIsPrivate bool   `json:"target_user_is_private"`
}

func (c *Client) Relationship(userID string) (*Relationship, error) {
	path := fmt.Sprintf("/users/%s/relationship", userID)

	resp, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var r Relationship
	err = json.Unmarshal(resp.Data, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (c *Client) Follows(userID string) *UsersFeed {
	return &UsersFeed{
		client:   c,
		init:     true,
		query:    userID,
		endPoint: "follows",
	}
}

func (c *Client) FollowedBy(userID string) *UsersFeed {
	return &UsersFeed{
		client:   c,
		init:     true,
		query:    userID,
		endPoint: "followed-by",
	}
}

func (c *Client) RequestedBy(userID string) *UsersFeed {
	return &UsersFeed{
		client:   c,
		init:     true,
		query:    userID,
		endPoint: "requested-by",
	}
}
