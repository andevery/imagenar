package autogram

import (
	// "time"
	"github.com/andevery/instaw"
)

type Follower struct {
	Limiter   *Limiter
	WithLikes bool
	Liker     *Liker
}

func (f *Follower) FollowBatch([]instaw.User) {

}
