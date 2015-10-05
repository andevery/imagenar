package autogram

import (
	"github.com/andevery/instax"
	"log"
)

type UnfollowWorker struct {
	exclude map[string]bool
	count   int
	client  *Client
}

func NewUnfollowWorker(client *Client, ids []string) *UnfollowWorker {
	w := &UnfollowWorker{client: client}
	w.exclude = make(map[string]bool)
	for _, id := range ids {
		w.exclude[id] = true
	}

	return w
}

func (w *UnfollowWorker) Start() {
	go func() {
		var users []instax.UserShort
		feed := w.client.Api().Follows("self")
		for u, err := feed.Next(); err != instax.EOF; u, err = feed.Next() {
			users = append(users, u...)
		}
		for i := range users {
			userID := users[len(users)-1-i].ID
			if _, ok := w.exclude[userID]; ok {
				continue
			}
			user, err := w.client.Api().User(userID)
			if err != nil {
				log.Println(err)
				continue
			}
			err = w.client.Unfollow(user)
			if err == nil {
				log.Printf("Unfollowed user %s", user.ID)
				w.count++
			}
		}
	}()
}

func (w *UnfollowWorker) Count() int {
	return w.count
}
