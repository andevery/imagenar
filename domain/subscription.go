package domain

import (
	"time"
)

type SubscriptionRepository interface {
	Store(s *Subscription)
	Find(id int64) *Subscription
}

type Subscription struct {
	ID        int64
	Customer  *Customer
	Plan      *Plan
	BeginDate time.Time
	EndDate   time.Time
}
