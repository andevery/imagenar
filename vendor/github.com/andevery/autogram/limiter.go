package autogram

import (
	"math/rand"
	"time"
)

type Limiter struct {
	Max struct {
		PerHour int
		PerDay  int
	}
	MaxDelay time.Duration
	MinDelay time.Duration
	Timer    chan time.Time

	counters struct {
		total int
		hour  int
		day   int
	}
	done chan bool
}

func NewLimiter() *Limiter {
	limiter := &Limiter{
		MaxDelay: 15 * time.Second,
		MinDelay: 5 * time.Second,
	}
	limiter.Max.PerHour = 60
	limiter.Timer = make(chan time.Time)

	limiter.startTimer()

	return limiter
}

func (l *Limiter) startTimer() {
	go func() {
		for {
			select {
			case <-done:
				close(l.Timer)
				return
			default:
				l.Timer <- time.Now()
				delay := time.Duration(rand.Int63n(int64(l.MaxDelay) - int64(l.MinDelay) + int64(l.MinDelay)))
				time.Sleep(delay)
			}
		}
	}()
}
