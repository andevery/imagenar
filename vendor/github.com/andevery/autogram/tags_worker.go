package autogram

import (
	"github.com/andevery/instax"
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

type TagsWorker struct {
	Follow       bool
	Like         bool
	LikesPerUser struct {
		Max int
		Min int
	}
	MediaCondition struct {
		MaxTags int
	}
	UserCondition struct {
		MaxFollowedBy int
		MinFollowedBy int
		MaxFollows    int
		MinFollows    int
		MinMedia      int
	}
	Delay time.Duration

	tags   []string
	counts struct {
		likes   uint32
		follows uint32
	}

	client *Client
}

func NewTagsWorker(client *Client, tags []string) *TagsWorker {
	return &TagsWorker{client: client, tags: tags}
}

func DefaultTagsWorker(client *Client, tags []string) *TagsWorker {
	w := NewTagsWorker(client, tags)
	w.Follow = true
	w.Like = true
	w.LikesPerUser.Min = 2
	w.LikesPerUser.Max = 4
	// w.MediaCondition.MaxTags = 15
	w.UserCondition.MaxFollowedBy = 500
	w.UserCondition.MaxFollows = 300
	w.UserCondition.MinFollows = 100
	w.UserCondition.MinMedia = 20
	w.Delay = 1 * time.Minute

	return w
}

func (w *TagsWorker) Start() {
	for _, t := range w.tags {
		feed := w.client.Api().MediaByTag(t)
		log.Printf("Feed get\n")
		go w.perform(feed)
	}
}

func (w *TagsWorker) perform(feed *instax.MediaFeed) {
	for {
		media, err := feed.Prev()
		log.Printf("Got %v media\n", len(media))
		if err != nil {
			log.Fatal(err)
		}
		for i, _ := range media {
			if !w.mediaMatch(&media[i]) {
				continue
			}
			user, ok := w.userMatch(media[i].User.ID)
			if !ok {
				continue
			}
			log.Printf("User %v got %v \n", user.ID, ok)

			if w.Like {
				recent, err := w.client.Api().RecentMediaByUser(media[i].User.ID)
				if err == nil {
					w.client.Web().Like(&recent[0])
					count := rand.Intn(w.LikesPerUser.Max-w.LikesPerUser.Min) + w.LikesPerUser.Min - 1
					w.client.LikeAFew(recent[1:], count)
					atomic.AddUint32(&w.counts.likes, uint32(count+1))
				}
			}
			if w.Follow {
				err := w.client.Web().Follow(user)
				if err == nil {
					atomic.AddUint32(&w.counts.follows, 1)
				}
			}
		}
		time.Sleep(w.Delay)
	}
}

func (w *TagsWorker) mediaMatch(media *instax.Media) bool {
	match := true
	if w.MediaCondition.MaxTags > 0 {
		match = match && len(media.Tags) <= w.MediaCondition.MaxTags
	}
	return match
}

func (w *TagsWorker) userMatch(userID string) (*instax.User, bool) {
	r, err := w.client.Api().Relationship(userID)
	if err != nil || r.TargetUserIsPrivate || r.OutgoingStatus != instax.NONE {
		return nil, false
	}

	user, err := w.client.Api().User(userID)
	if err != nil {
		return nil, false
	}

	flag := true

	if w.UserCondition.MaxFollowedBy > 0 {
		flag = flag && user.Counts.FollowedBy <= w.UserCondition.MaxFollowedBy
	}
	if w.UserCondition.MinFollowedBy > 0 {
		flag = flag && user.Counts.FollowedBy >= w.UserCondition.MinFollowedBy
	}
	if w.UserCondition.MaxFollows > 0 {
		flag = flag && user.Counts.Follows <= w.UserCondition.MaxFollows
	}
	if w.UserCondition.MinFollows > 0 {
		flag = flag && user.Counts.Follows >= w.UserCondition.MinFollows
	}
	if w.UserCondition.MinMedia > 0 {
		flag = flag && user.Counts.Media >= w.UserCondition.MinMedia
	}

	return user, flag
}

func (w *TagsWorker) LikesCount() int {
	return int(atomic.LoadUint32(&w.counts.likes))
}

func (w *TagsWorker) FollowsCount() int {
	return int(atomic.LoadUint32(&w.counts.follows))
}
