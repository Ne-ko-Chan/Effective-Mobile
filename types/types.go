package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type SubscriptionStore interface {
	GetSumPeriod(from, to, uuid, serviceName string) (int, error)
	GetSubscriptionByID(id int) (*Subscription, error)
	GetSubscriptions() ([]Subscription, error)
	CreateSubscription(s Subscription) error
	UpdateSubscription(s Subscription) error
	DeleteSubscription(id int) error
}

type Subscription struct {
	ID          int          `json:"id"`
	ServiceName string       `json:"service_name"`
	Price       int          `json:"price"`
	UserID      string       `json:"user_id"`
	StartDate   time.Time    `json:"start_date"`
	EndDate     sql.NullTime `json:"end_date"`
}

// TODO: check out more validate options
type CreateSubscriptionPayload struct {
	ServiceName string     `json:"service_name" validate:"required"`
	Price       int        `json:"price"        validate:"required"`
	UserID      string     `json:"user_id"      validate:"required"`
	StartDate   CustomTime `json:"start_date"   validate:"required"`
	EndDate     CustomTime `json:"end_date"`
}

type UpdateSubscriptionPayload struct {
	ID          int        `json:"id" validate:"required"`
	ServiceName string     `json:"service_name"    validate:"required"`
	Price       int        `json:"price"           validate:"required"`
	UserID      string     `json:"user_id"         validate:"required"`
	StartDate   CustomTime `json:"start_date"      validate:"required"`
	EndDate     CustomTime `json:"end_date"`
}

type DeleteSubscriptionPayload struct {
	ID int `json:"id" validate:"required"`
}

type CustomTime time.Time

func (t CustomTime) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%s", time.Time(t).Format("01-2006")), nil
}

func (t *CustomTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := time.Parse("01-2006", s)
	if err != nil {
		return fmt.Errorf("invalid format: %s, expected MM-YYYY", s)
	}

	*t = CustomTime(parsed)
	return nil
}

func (t CustomTime) Time() time.Time {
	return time.Time(t)
}

func (t CustomTime) String() string {
	return t.Time().Format("01-2006")
}
