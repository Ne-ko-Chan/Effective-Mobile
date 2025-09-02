CREATE TABLE IF NOT EXISTS subscriptions (
  service_name VARCHAR(255) NOT NULL,
  price        INTEGER      NOT NULL,
  user_id      UUID         NOT NULL,
  start_date   DATE         NOT NULL,
  end_date     DATE,
  PRIMARY KEY(service_name, user_id)
)
