# Nanaban

gRPC service that fetches order book data from the Grinex exchange and calculates bid/ask rates using configurable methods.

## Requirements

- [Docker & Docker Compose](https://docs.docker.com/compose/install/)
- [grpcurl](https://github.com/fullstorydev/grpcurl) — for manual testing

```sh
brew install grpcurl
```

## Running

```sh
cp .env.example .env
docker compose up --build -d
```

## Configuration

Copy `.env.example` to `.env` and adjust as needed.

| Variable          | Description                        | Default       |
|-------------------|------------------------------------|---------------|
| `APP_NAME`        | Service name                       | `nanaban`     |
| `APP_ENV`         | Environment (`dev`/`production`)   | `dev`         |
| `APP_HOST`        | Server host                        | `http://localhost` |
| `APP_PORT`        | gRPC server port                   | `8080`        |
| `PROMETHEUS_PORT` | Prometheus metrics port            | `9090`        |
| `GOMEMLIMIT`      | Soft memory limit for GC           | `400MiB`      |
| `DB_HOST`         | Postgres host                      | `postgres`    |
| `DB_PORT`         | Postgres port                      | `5432`        |
| `DB_USER`         | Postgres user                      | `main`        |
| `DB_PASSWORD`     | Postgres password                  | `secret`      |
| `DB_NAME`         | Postgres database name             | `postgres`    |
| `LOG_LEVEL`       | Log level (`debug`/`info`/`error`) | `debug`       |

## Testing

Reflection is enabled — no `-proto` flag needed.

### Health Check

```sh
grpcurl -plaintext localhost:{APP_PORT} grpc.health.v1.Health/Check
```

### GetRates — TOP_N

Returns bid and ask at position `n` in the order book. `n: 0` is the top.

```sh
grpcurl -plaintext \
  -d '{"method": 1, "n": 0}' \
  localhost:{APP_PORT} \
  v1.RateService/GetRates
```

### GetRates — AVG_NM

Returns the average bid and ask across entries from index `n` to `m` (inclusive).

```sh
grpcurl -plaintext \
  -d '{"method": 2, "n": 0, "m": 4}' \
  localhost:{APP_PORT} \
  v1.RateService/GetRates
```

## Metrics

Prometheus metrics are exposed at:

```
http://localhost:{PROMETHEUS_PORT}/metrics
```

Default port: `9090`.