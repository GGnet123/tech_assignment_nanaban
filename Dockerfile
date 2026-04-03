FROM golang:1.26.1-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/app

FROM alpine:3.23

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/app ./app

CMD ["./app"]