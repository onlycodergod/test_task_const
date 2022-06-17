include .env
export

run:
	go run ./cmd/app
.PHONY: run

tidy:
	go mod tidy
.PHONY: tidy

test:
	go test -v -cover -race ./internal/...
.PHONY: test

migrate-create:
	migrate create -ext sql -dir migrations 'scheme'
.PHONY: migrate-create

migrate-up:
	migrate -path migrations -database '$(POSTGRES_DSN)' up
.PHONY: migrate-up

migrate-down:
	migrate -path migrations -database '$(POSTGRES_DSN)' down
.PHONY: migrate-up

compose-up:
	docker-compose up --build && docker-compose logs -f
.PHONY: compose-up

compose-down:
	docker-compose down --remove-orphans --volumes
.PHONY: compose-down

remove-volumes:
	docker system prune -a --volumes
.PHONY: remove-volumes

linter-golangci:
	CGO_ENABLED=0 golangci-lint run
.PHONY: linter-golangci