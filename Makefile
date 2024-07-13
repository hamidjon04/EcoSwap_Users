CURRENT_DIR := $(shell pwd)
DATABASE_URL="postgres://postgres:hamidjon4424@localhost:5432/ecoswap?sslmode=disable"

run-router:
	@go run cmd/router/main.go

run-service:
	@go run cmd/service/main.go

gen-proto:
	@./scripts/gen-proto.sh $(CURRENT_DIR)

tidy:
	@go mod tidy

mig-create:
	@if [ -z "$(name)" ]; then \
	read -p "Enter migration name: " name; \
	fi; \
	migrate create -ext sql -dir db/migrations -seq $$name

mig-up:
	@migrate -database "$(DATABASE_URL)" -path db/migrations up

mig-down:
	@migrate -database "$(DATABASE_URL)" -path db/migrations down

mig-force:
	@if [ -z "$(version)" ]; then \
	read -p "Enter migration version: " version; \
	fi; \
	migrate -database "$(DATABASE_URL)" -path db/migrations force $$version

permission:
	@chmod +x scripts/gen-proto.sh

test:
	@go test ./storage/postgres

swag-gen:
	~/go/bin/swag init -g ./api/router.go -o api/docs

