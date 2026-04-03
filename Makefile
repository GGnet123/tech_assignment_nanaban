migrate-up:
	docker compose run migrate up

migrate-down:
	docker compose run migrate down

migrate-create:
	docker compose run migrate create -ext sql -dir /migrations -seq $(name)

build:
	docker compose build
docker-build:
	docker compose up --build -d
run:
	docker compose up -d