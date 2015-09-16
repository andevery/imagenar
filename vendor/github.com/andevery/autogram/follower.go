package autogram

import (
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
		MaxMedia      int
	}

	WebClient *instaw.Client
	ApiClient *instax.Client
}

func (f *Follower) FollowABatch(users []instax.User) {
	for i, _ := range users {
		if f.isUserMatch(&users[i]) {
			if f.WithLikes {
				media, err := f.ApiClient.RecentMediaByUser(users[i].ID)
				if err == nil {
					f.Liker.LikeAFew(media, rand.Intn(f.MaxLikes-f.MinLikes)+f.MinLikes)
				}
			}
			<-f.Limiter.Timer
			f.WebClient.Follow(&users[i])
		}
	}
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
	if f.UsersCondition.MaxMedia > 0 {
		flag = flag && user.Counts.Media <= f.UsersCondition.MaxMedia
	}

	return flag
}
