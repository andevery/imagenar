package main

import (
	"fmt"
	"github.com/andevery/instax"
	"log"
	"math/rand"
	"net/http"
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
	FeedDepth      int
	UserCond       struct {
		FollowedBy int
		Follows    int
		Media      int
	}
	MediaCond struct {
		TagsCount int
	}

	tags   []string
	client *instax.Client

	breakTime   time.Duration
	counter     int
	depth       int
	likesNumber int
}

func NewLiker(client *instax.Client, tags []string) *Liker {
	liker := new(Liker)
	liker.client = client
	liker.tags = tags

	liker.Min = 60
	liker.Max = 90
	liker.MinBreak = 50 * time.Minute
	liker.MaxBreak = 70 * time.Minute
	liker.RateLimitPause = 20 * time.Minute
	liker.FeedDepth = 10
	liker.UserCond.FollowedBy = 500
	liker.UserCond.Follows = 200
	liker.UserCond.Media = 50
	liker.MediaCond.TagsCount = 10

	return liker
}

func (l *Liker) flushCounts() {
	rand.Seed(time.Now().Unix())
	l.likesNumber = rand.Intn(l.Max-l.Min) + l.Min
	l.breakTime = time.Duration(rand.Int63n(int64(l.MaxBreak)-int64(l.MinBreak)) + int64(l.MinBreak))
	l.counter = 0
}

func (l *Liker) isUserMatch(user *instax.User) bool {
	return user.Counts.FollowedBy <= l.UserCond.FollowedBy &&
		user.Counts.Follows <= l.UserCond.Follows &&
		user.Counts.Media >= l.UserCond.Media
}

func (l *Liker) fakeLoadImages(media []instax.Media) {
	for i, _ := range media {
		http.Get(media[i].Images.Thumbnail.URL)
	}
}

func (l *Liker) checkAndLike(media []instax.Media) {
	var user *instax.User
	var err error

	if l.client.Limit() < 500 {
		fmt.Println("")
		log.Printf("Rate limit pause. Resume after %s.", l.RateLimitPause)
		time.Sleep(l.RateLimitPause)
	}

	for i, _ := range media {
		if media[i].UserHasLiked || len(media[i].Tags) > l.MediaCond.TagsCount {
			continue
		}

		user, err = l.client.User(media[i].User.ID)
		if err != nil {
			log.Fatal(err)
			return
		}

		if l.isUserMatch(user) {
			http.Get(media[i].Images.LowResolution.URL)
			err = l.client.Like(media[i].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(" ✔ ")
			l.counter++
			if l.counter >= l.likesNumber {
				fmt.Println("")
				log.Printf("Break time. Liked %v photos. Resume after %s", l.counter, l.breakTime)
				time.Sleep(l.breakTime)
				l.flushCounts()
			}
		}
	}
}

func (l *Liker) Start() {
	l.flushCounts()
	for _, tag := range l.tags {
		feed := l.client.MediaByTag(tag)
		media, err := feed.Get()
		if err != nil {
			log.Fatal(err)
		}
		l.depth = 1
		l.checkAndLike(media)
		for {
			l.depth++
			if !feed.CanNext() || l.depth > l.FeedDepth {
				fmt.Println("")
				log.Printf("End of feed: #%s.", tag)
				break
			}
			media, err = feed.Next()
			if err != nil {
				log.Fatal(err)
			}
			l.checkAndLike(media)
		}
	}
}

func main() {

	client := instax.NewClient("2079178474.7553a48.6c4d40d7782147cab2b5ff5ab44428e1", "5ac3e50811cc47c2a4cd1adda782eb4b")
	client.Delayer = func() time.Duration {
		return time.Duration(rand.Intn(7)+3) * time.Second
	}

	liker := NewLiker(client, []string{"city"})
	liker.Start()

	// feed := client.MediaByTag("animallovers")
	// m, _ := feed.Get()
	// CheckAndLike(client, m)
	// for {
	// 	if count > 1 {
	// 		return
	// 	}
	// 	m, _ = feed.Next()
	// 	CheckAndLike(client, m)
	// }
}
