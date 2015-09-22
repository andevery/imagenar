package autogram

import (
	"log"
	"math/rand"
	// "time"
	"github.com/andevery/instax"
)

type Follower struct {
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

	Provider *Provider
}

func DefaultFollower(p *Provider, l *Liker) *Follower {
	f := &Follower{
		MinLikes: 7,
		MaxLikes: 11,
		Liker:    l,
		Provider: p,
	}

	f.UsersCondition.MaxFollowedBy = 500
	f.UsersCondition.MaxFollows = 300
	f.UsersCondition.MinFollows = 100
	f.UsersCondition.MinMedia = 50

	return f
}

func (f *Follower) Start() {
	for {
		users, err := f.Provider.NextUsers()
		if err == EOF {
			return
		} else if err == instax.NotFound {
			continue
		} else if err != nil {
			log.Fatal(err)
		}

		f.FollowABatch(users)
	}
}

func (f *Follower) FollowAFew(users []instax.UserShort, count int) {
	for _, i := range randomIndexes(len(users), count) {
		if _, ok := f.isUserMatch(users[i].ID); ok {
			// f.Provider.WebClient().Follow(u)
			if f.Liker != nil {
				media, err := f.Provider.ApiClient().RecentMediaByUser(users[i].ID)
				if err == nil {
					f.Liker.LikeAFew(media, rand.Intn(f.MaxLikes-f.MinLikes)+f.MinLikes)
				}
			}
		}
	}
}

func (f *Follower) FollowABatch(users []instax.UserShort) {
	f.FollowAFew(users, len(users))
}

func (f *Follower) isUserMatch(userID string) (*instax.User, bool) {
	r, err := f.Provider.ApiClient().Relationship(userID)
	if err != nil || r.TargetUserIsPrivate || r.OutgoingStatus != instax.NONE {
		return nil, false
	}

	user, err := f.Provider.ApiClient().User(userID)
	if err != nil {
		return nil, false
	}

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

	return user, flag
}
