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
		if count > 30 {
			return
		}
		time.Sleep(time.Duration(rand.Int31n(6)) * time.Second)
	}
}

func main() {
	client := instax.NewClient("2079178474.1fb234f.682a311e35334df3842ccb654516baf5")
	feed := client.MediaByTag("ночь")
	m, _ := feed.Get()
	CheckAndLike(client, m)
	for {
		if count > 30 {
			return
		}
		m, _ = feed.Next()
		CheckAndLike(client, m)
	}
}
