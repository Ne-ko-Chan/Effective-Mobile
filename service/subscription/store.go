package subscription

import (
	"database/sql"
	"rest-service/types"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetSumPeriod(from, to time.Time, uuid, serviceName string) {
}

func (s *Store) GetSubscriptionByID(id int) {

}

func (s *Store) CreateSubscription(sub types.Subscription) error {
	return nil
}
