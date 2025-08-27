package types

import "time"

type SubscriptionStore interface {
	GetSumPeriod(from, to time.Time, uuid, serviceName string)
	GetSubscriptionByID(id int)
	CreateSubscription(s Subscription) error
}

type Subscription struct {
	ID          int
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
	EndDate     time.Time
}
