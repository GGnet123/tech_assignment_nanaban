include .env
export

test:
	go test ./... -race

lint:
	docker compose --profile tools run --rm lint
lint-fix:
	docker compose --profile tools run --rm lint golangci-lint run --fix ./...
proto-gen:
	protoc \
		--proto_path=proto \
		--go_out=pkg/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=pkg/pb \
		--go-grpc_opt=paths=source_relative \
		proto/v1/*.proto

MIGRATE_CMD=docker compose run --rm migrate -path /migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable&search_path=public"

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

migrate-create:
	$(MIGRATE_CMD) create -ext sql -dir /migrations -seq $(name)

build:
	docker compose build
docker-build:
	docker compose up --build -d
run:
	docker compose up -d