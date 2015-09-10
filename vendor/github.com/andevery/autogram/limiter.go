package autogram

import (
	"math/rand"
	"sync/atomic"
	"time"
)

type Limiter struct {
	Rate struct {
		HourLimit int
		DayLimit  int
	}
	MaxDelay time.Duration
	MinDelay time.Duration
	Timer    chan time.Time

	counters struct {
		total int32
		hour  int32
		day   int32
	}
	done chan bool
}

func NewLimiter() *Limiter {
	limiter := &Limiter{
		MaxDelay: 2 * time.Second,
		MinDelay: 1 * time.Second,
	}
	limiter.Rate.HourLimit = 60
	limiter.Timer = make(chan time.Time)
	limiter.done = make(chan bool)

	limiter.startTimer()

	return limiter
}

func (l *Limiter) Stop() {
	l.done <- true
}

func (l *Limiter) startTimer() {
	go func() {
		for {
			select {
			case <-l.done:
				close(l.Timer)
				return
			default:
				l.Timer <- time.Now()
				l.incCounters()
				delay := time.Duration(rand.Int63n(int64(l.MaxDelay) - int64(l.MinDelay) + int64(l.MinDelay)))
				time.Sleep(delay)
			}
		}
	}()
}

func (l *Limiter) notLimited() bool {
	flag := true
	if l.Rate.DayLimit > 0 {
		flag = flag && l.counters.day <= l.Rate.DayLimit
	}

	if l.Rate.HourLimit > 0 {
		flag = flag && l.counters.hour <= l.Rate.HourLimit
	}
	return flag
}

func (l *Limiter) incCounters() {
	atomic.AddInt32(&l.counters.day, 1)
	atomic.AddInt32(&l.counters.hour, 1)
	atomic.AddInt32(&l.counters.total, 1)
}
