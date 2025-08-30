FROM golang:1.25-alpine

COPY ./cmd/migrate/ /app/cmd/migrate/

COPY ./config /app/config/

COPY ./go.mod /app/

COPY ./db/ /app/db/

WORKDIR /app

RUN go mod tidy
