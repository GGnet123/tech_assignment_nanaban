# Nanaban

gRPC service that fetches order book data from the Grinex exchange and calculates bid/ask rates using configurable methods.

## Requirements

- Docker & Docker Compose
- [grpcurl](https://github.com/fullstorydev/grpcurl) for manual testing

```sh
brew install grpcurl
```

## Running

```sh
cp .env.example .env
docker compose up --build -d
```

## Testing

The server has reflection enabled, so no `-proto` flag is needed.

### Health Check

Checks whether the gRPC server is up and ready to serve.

```sh
grpcurl -plaintext localhost:8099 grpc.health.v1.Health/Check
```

### GetRates — TOP_N

Returns the bid and ask prices at position `n` in the order book.

```sh
grpcurl -plaintext \
  -d '{"method": 1, "n": 0}' \
  localhost:{YOUR_PORT} \
  v1.RateService/GetRates
```

`n: 0` returns the top of the book. Increase `n` to go deeper.

### GetRates — AVG_NM

Returns the average bid and ask prices across entries from index `n` to `m` (inclusive).

```sh
grpcurl -plaintext \
  -d '{"method": 2, "n": 0, "m": 4}' \
  localhost:{YOUR_PORT} \
  v1.RateService/GetRates
```

## Configuration

| Variable      | Description                  | Default   |
|---------------|------------------------------|-----------|
| `APP_HOST`    | Server host                  | `0.0.0.0` |
| `APP_PORT`    | gRPC server port             | `8080`    |
| `DB_HOST`     | Postgres host                | `postgres`|
| `DB_PORT`     | Postgres port                | `5432`    |
| `DB_USER`     | Postgres user                | `main`    |
| `DB_PASSWORD` | Postgres password            | `secret`  |
| `DB_NAME`     | Postgres database name       | `postgres`|