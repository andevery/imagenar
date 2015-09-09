package autogram

import (
	"time"
)

type Follower struct {
}

type Limiter struct {
	Max struct {
		PerHour int
		PerDay  int
	}
	MaxDelay time.Duration
	MinDelay time.Duration

	counter int
}
