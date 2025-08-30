build:
	@go build -o bin/service cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/service

migration:
	@~/go/bin/migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-newest:
	@go run cmd/migrate/main.go newest

migrate-down:
	@go run cmd/migrate/main.go down
