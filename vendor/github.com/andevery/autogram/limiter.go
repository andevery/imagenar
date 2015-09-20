package autogram

import (
	"github.com/andevery/instaw"
	"github.com/andevery/instax"
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

	counters struct {
		total uint32
		hour  uint32
		day   uint32
	}
	webTickers struct {
		day  *time.Ticker
		hour *time.Ticker
	}
	apiTicker *time.Ticker

	timer     chan time.Time
	done      chan bool
	apiClient *instax.Client
	webClient *instaw.Client
}

func DefaultLimiter(api *instax.Client, web *instaw.Client) *Limiter {
	l := &Limiter{
		apiClient: api,
		webClient: web,
	}
	l.MaxDelay = 25 * time.Second
	l.MinDelay = 15 * time.Second
	l.Rate.HourLimit = 180
	l.timer = make(chan time.Time)
	l.done = make(chan bool)
	l.webTickers.day = time.NewTicker(24 * time.Hour)
	l.webTickers.hour = time.NewTicker(1 * time.Hour)

	l.Start()

	return l
}

func (l *Limiter) ApiClient() *instax.Client {
	if l.apiTicker == nil {
		l.apiTicker = time.NewTicker(1 * time.Hour)
	}
	if l.apiClient.Limit() < 500 {
		<-l.apiTicker.C
	}
	return l.apiClient
}

func (l *Limiter) WebClient() *instaw.Client {
	<-l.timer
	return l.webClient
}

func (l *Limiter) Stop() {
	l.done <- true
}

func (l *Limiter) Start() {
	go func() {
		for {
			select {
			case <-l.webTickers.day.C:
				atomic.StoreUint32(&l.counters.day, 0)
			case <-l.webTickers.hour.C:
				atomic.StoreUint32(&l.counters.hour, 0)
			case <-l.done:
				close(l.timer)
				l.webTickers.day.Stop()
				l.webTickers.hour.Stop()
				return
			default:
				if l.notLimited() {
					l.timer <- time.Now()
					l.incCounters()
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
		flag = flag && l.counters.day < l.Rate.DayLimit
	}

	if l.Rate.HourLimit > 0 {
		flag = flag && l.counters.hour < l.Rate.HourLimit
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
