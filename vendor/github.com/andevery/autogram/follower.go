package autogram

import (
	"log"
	"math/rand"
	// "time"
	"github.com/andevery/instaw"
	"github.com/andevery/instax"
)

type Follower struct {
	Limiter        *Limiter
	WithLikes      bool
	MinLikes       int
	MaxLikes       int
	Liker          *Liker
	UsersCondition struct {
		MaxFollowedBy int
		MinFollowedBy int
		MaxFollows    int
		MinFollows    int
		MinMedia      int
	}

	WebClient *instaw.Client
	ApiClient *instax.Client
}

func DefaultFollower(limiter *Limiter) {

}

func (f *Follower) FollowAFew(users []instax.UserShort, count int) {
	for _, i := range randomIndexes(len(users), count) {
		u, err := f.ApiClient.User(users[i].ID)
		if err != nil {
			log.Println(err, users[i].ID)
			continue
		}
		if f.isUserMatch(u) {
			if f.WithLikes {
				media, err := f.ApiClient.RecentMediaByUser(users[i].ID)
				if err == nil {
					f.Liker.LikeAFew(media, rand.Intn(f.MaxLikes-f.MinLikes)+f.MinLikes)
				}
			}
			<-f.Limiter.Timer
			f.WebClient.Follow(u)
		}
	}
}

func (f *Follower) FollowABatch(users []instax.UserShort) {
	f.FollowAFew(users, len(users))
}

func (f *Follower) isUserMatch(user *instax.User) bool {
	flag := true

	if f.UsersCondition.MaxFollowedBy > 0 {
		flag = flag && user.Counts.FollowedBy <= f.UsersCondition.MaxFollowedBy
	}
	if f.UsersCondition.MinFollowedBy > 0 {
		flag = flag && user.Counts.FollowedBy >= f.UsersCondition.MinFollowedBy
	}
	if f.UsersCondition.MaxFollows > 0 {
		flag = flag && user.Counts.Follows <= f.UsersCondition.MaxFollows
	}
	if f.UsersCondition.MinFollows > 0 {
		flag = flag && user.Counts.Follows >= f.UsersCondition.MinFollows
	}
	if f.UsersCondition.MinMedia > 0 {
		flag = flag && user.Counts.Media >= f.UsersCondition.MinMedia
	}

	return flag
}
