-include .env

tidy:
	go mod tidy

.PHONY: tidy

run:
	@go run cmd/api/main.go

.PHONY: run

migrateUp:
	@go run util/migrations/main.go up ${DB_DSN}

.PHONY: migrateUp

migrateDown:
	go run util/migrations/main.go down ${DB_DSN}

.PHONY: migrateDown

lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2
	@golangci-lint run

.PHONY: lint