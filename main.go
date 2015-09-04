package main

import (
	"github.com/andevery/instax"
	"log"
	"math/rand"
	"time"
)

var (
	count = 0
)

type Liker struct {
	Min            int
	Max            int
	MinBreak       time.Duration
	MaxBreak       time.Duration
	RateLimitPause time.Duration
	UserCond       struct {
		FollowedBy int
		Follows    int
		Media      int
	}

	tags        []string
	likesNumber int
	counter     int
	client      *instax.Client
}

func NewLiker(client *instax.Client, tags []string) *Liker {
	liker := new(Liker)
	liker.client = client
	liker.tags = tags

	liker.Min = 20
	liker.Max = 40
	liker.MinBreak = 40 * time.Minute
	liker.MaxBreak = 80 * time.Minute
	liker.RateLimitPause = 20 * time.Minute
	liker.UserCond.FollowedBy = 500
	liker.UserCond.Follows = 200
	liker.UserCond.Media = 50

	return liker
}

func (l *Liker) checkAndLike(media []instax.Media) {
	var user *instax.User
	var err error

	if l.client.Limit() < 100 {
		time.Sleep(l.RateLimitPause)
	}

	for i, _ := range media {
		if media[i].UserHasLiked {
			continue
		}

		user, err = l.client.User(media[i].User.ID)
		if err != nil {
			log.Println(err)
			return
		}

		if user.Counts.FollowedBy <= l.UserCond.FollowedBy &&
			user.Counts.Follows <= l.UserCond.Follows &&
			user.Counts.Media >= l.UserCond.Media {
			l.client.Like(media[i].ID)
			l.counter++
			if counter == l.Max {
				return
			}
		}
	}
}

func (l *Liker) Perform() {
	for {
		likesNumber := rand.Intn(l.Max-l.Min) + l.Min
		breakTime := time.Duration(rand.Int63n(int64(l.MaxBreak)-int64(l.MinBreak)) + int64(l.MinBreak))

		for _, tag := range l.tags {
			feed := l.client.MediaByTag(tag)
		}
	}
}

func CheckAndLike(c *instax.Client, media []instax.Media) {
	var user *instax.User
	for i, _ := range media {
		user, _ = c.User(media[i].User.ID)
		if user.Counts.FollowedBy < 500 && user.Counts.Follows < 200 && user.Counts.Media > 20 {
			c.Like(media[i].ID)
			log.Println("Liked")
			log.Println(c.Limit())
			count++
		}
		if count > 1 {
			return
		}
		time.Sleep(time.Duration(rand.Int31n(6)) * time.Second)
	}
}

func main() {
	client := instax.NewClient("2079178474.1fb234f.682a311e35334df3842ccb654516baf5")
	liker := NewLiker(client, []string{"animallovers"})
	liker.Perform()

	feed := client.MediaByTag("animallovers")
	m, _ := feed.Get()
	CheckAndLike(client, m)
	for {
		if count > 1 {
			return
		}
		m, _ = feed.Next()
		CheckAndLike(client, m)
	}
}
