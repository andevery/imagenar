package autogram

import (
	// "time"
	"github.com/andevery/instax"
)

type Follower struct {
	Limiter        *Limiter
	WithLikes      bool
	Liker          *Liker
	UsersCondition struct {
		MaxFollowedBy int
		MinFollowedBy int
		MaxFollows    int
		MinFollows    int
		MaxMedia      int
	}

	Insta *instax.Client
}

func (f *Follower) FollowBatch(users []instax.User) {
	for i, _ := range users {

	}
}

func (f *Follower) isUserMathc(user *instax.User) bool {
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
