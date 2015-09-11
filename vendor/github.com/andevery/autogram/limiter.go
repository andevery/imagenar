package autogram

import (
	"math/rand"
	"sync/atomic"
	"time"
)

type Limiter struct {
	Rate struct {
		HourLimit uint32
		DayLimit  uint32
	}
	MaxDelay time.Duration
	MinDelay time.Duration
	Timer    chan time.Time

	counters struct {
		total uint32
		hour  uint32
		day   uint32
	}
	tickers struct {
		day  *time.Ticker
		hour *time.Ticker
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

	limiter.tickers.day = time.NewTicker(24 * time.Hour)
	limiter.tickers.hour = time.NewTicker(1 * time.Hour)

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
			case <-l.tickers.day.C:
				l.counters.day = 0
			case <-l.tickers.hour.C:
				l.counters.hour = 0
			case <-l.done:
				close(l.Timer)
				l.tickers.day.Stop()
				l.tickers.hour.Stop()
				return
			default:
				if l.notLimited() {
					l.incCounters()
					l.Timer <- time.Now()
				}
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
	atomic.AddUint32(&l.counters.day, 1)
	atomic.AddUint32(&l.counters.hour, 1)
	atomic.AddUint32(&l.counters.total, 1)
}

func (l *Limiter) TotalAmount() uint32 {
	return atomic.LoadUint32(&l.counters.total)
}

func (l *Limiter) HourAmount() uint32 {
	return atomic.LoadUint32(&l.counters.hour)
}

func (l *Limiter) DayAmount() uint32 {
	return atomic.LoadUint32(&l.counters.day)
}
