CREATE TABLE IF NOT EXISTS subscriptions (
  service_name VARCHAR(255) NOT NULL,
  price        INTEGER      NOT NULL,
  user_id      UUID         PRIMARY KEY,
  start_date   DATE         NOT NULL,
  end_date     DATE
)
