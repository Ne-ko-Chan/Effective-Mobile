package subscription

import (
	"database/sql"
	"fmt"
	"rest-service/types"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetSumPeriod(from, to time.Time, uuid, serviceName string) (int, error) {
	query := "SELECT SUM(price) FROM subscriptions WHERE start_date >= $1 AND start_date <= $2"
	params := []any{from, to}
	paramCount := 2

	if uuid != "" {
		paramCount++
		params = append(params, uuid)
		query += fmt.Sprintf(" AND user_id = $%d", paramCount)
	}

	if serviceName != "" {
		paramCount++
		params = append(params, serviceName)
		query += fmt.Sprintf(" AND service_name = $%d", paramCount)
	}

	var sum int
	err := s.db.QueryRow(query, params...).Scan(&sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (s *Store) GetSubscriptionByID(id int) (*types.Subscription, error) {
	row := s.db.QueryRow("SELECT id, user_id, service_name, price, start_date, end_date FROM subscriptions WHERE id = $1", id)

	var sub types.Subscription
	err := row.Scan(
		&sub.ID,
		&sub.UserID,
		&sub.ServiceName,
		&sub.Price,
		&sub.StartDate,
		&sub.EndDate,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("there is no subscription with requested id")
		}
		return nil, err
	}

	return &sub, nil
}

func (s *Store) GetSubscriptions() ([]types.Subscription, error) {
	rows, err := s.db.Query("SELECT id, user_id, service_name, price, start_date, end_date FROM subscriptions")
	if err != nil {
		return nil, err
	}

	var subscriptions []types.Subscription
	var sub types.Subscription
	for rows.Next() {
		err = rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.ServiceName,
			&sub.Price,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			return  nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (s *Store) CreateSubscription(sub types.Subscription) error {
	if !sub.EndDate.Valid {
		_, err := s.db.Exec("INSERT INTO subscriptions(user_id, service_name, price, start_date) VALUES ($1,$2,$3,$4)", sub.UserID, sub.ServiceName, sub.Price, sub.StartDate)
		if err != nil {
			return err
		}
	} else {
		_, err := s.db.Exec("INSERT INTO subscriptions(user_id, service_name, price, start_date, end_date) VALUES ($1,$2,$3,$4,$5)", sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) UpdateSubscription(sub types.Subscription) error {
	_, err := s.GetSubscriptionByID(sub.ID)
	if err != nil {
		return err
	}
	if !sub.EndDate.Valid {
		_, err := s.db.Exec("UPDATE subscriptions SET user_id=$1, service_name=$2, price=$3, start_date=$4 WHERE id=$5", sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.ID)
		if err != nil {
			return err
		}
	} else {
		_, err := s.db.Exec("UPDATE subscriptions SET user_id=$1, service_name=$2, price=$3, start_date=$4, end_date=$5 WHERE id=$6", sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, sub.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) DeleteSubscription(id int) error {
	_, err := s.GetSubscriptionByID(id)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM subscriptions WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
