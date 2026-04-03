proto-gen:
	protoc \
		--proto_path=proto \
		--go_out=pkg/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=pkg/pb \
		--go-grpc_opt=paths=source_relative \
		proto/v1/*.proto

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