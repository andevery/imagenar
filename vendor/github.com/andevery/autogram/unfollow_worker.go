package autogram

import (
	"github.com/andevery/instax"
	"log"
	"sync"
)

type UnfollowWorker struct {
	ID        int64
	exclude   map[string]bool
	count     int
	client    *Client
	send      Reporter
	done      chan bool
	waitGroup *sync.WaitGroup
}

func NewUnfollowWorker(id int64, client *Client, whitelist []string, reporter Reporter) *UnfollowWorker {
	w := &UnfollowWorker{
		ID:        id,
		client:    client,
		done:      make(chan bool),
		send:      reporter,
		waitGroup: &sync.WaitGroup{},
	}
	w.exclude = make(map[string]bool)
	for _, username := range whitelist {
		w.exclude[username] = true
	}

	return w
}

func (w *UnfollowWorker) Start() {
	go func() {
		defer w.waitGroup.Done()
		var users []instax.UserShort
		feed := w.client.Api().Follows("self")
		for u, err := feed.Next(); err != instax.EOF; u, err = feed.Next() {
			if w.stopped() {
				return
			}
			users = append(users, u...)
		}
		for i := range users {
			if w.stopped() {
				return
			}
			username := users[len(users)-1-i].Username
			if _, ok := w.exclude[username]; ok {
				continue
			}
			userID := users[len(users)-1-i].ID
			user, err := w.client.Api().User(userID)
			if err != nil {
				w.error(err)
				continue
			}
			err = w.client.Web().Unfollow(user)
			if err != nil {
				w.fatal(err)
				return
			}
			log.Printf("Unfollowed user %s", user.ID)
			w.count++
			w.report()
		}
	}()
	w.waitGroup.Add(1)
	log.Println("Unfollowing task started")
}

func (w *UnfollowWorker) Stop() {
	select {
	case <-w.done:
		return
	default:
		close(w.done)
		w.waitGroup.Wait()
		log.Println("Unfollowing task stopped")
	}
}

func (w *UnfollowWorker) stopped() bool {
	select {
	case <-w.done:
		return true
	default:
	}
	return false
}

func (w *UnfollowWorker) Count() int {
	return w.count
}

func (w *UnfollowWorker) report() {
	if w.send != nil {
		report := map[string]int{
			"unfollowed": w.Count(),
		}
		w.send.Report(w.ID, report)
	}
}

func (w *UnfollowWorker) error(err error) {
	if w.send != nil {
		w.send.Error(w.ID, err)
	}
}

func (w *UnfollowWorker) fatal(err error) {
	if w.send != nil {
		w.send.Fatal(w.ID, err)
	}
}
