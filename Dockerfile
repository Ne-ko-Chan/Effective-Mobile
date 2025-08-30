FROM golang:1.25-alpine

COPY . /app

WORKDIR /app

RUN go mod tidy

RUN go build -o build/service cmd/main.go

CMD [ "./build/service" ]
