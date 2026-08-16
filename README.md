# Aftercare Service

Aftercare Service is a small HTTP service for accepting after-sales issue descriptions. It validates the request, normalizes whitespace, applies the supported description length, stores the ticket in memory, and returns Unicode character statistics together with loading, error, and result states.

## Requirements

- Go 1.24.13
- `GOTOOLCHAIN=local`

## Run

From the module root:

```sh
export GOTOOLCHAIN=local
go run ./cmd/aftercare
```

The service listens on `:8080` by default. Set `AFTERCARE_ADDR` to use another address.

Create a ticket:

```sh
curl -i \
  -H 'Content-Type: application/json' \
  -d '{"description":"  商品   外包装破损，需要售后处理  "}' \
  http://localhost:8080/api/tickets
```

Read it back using the returned ticket ID:

```sh
curl -i http://localhost:8080/api/tickets/T-000001
```

Check service health:

```sh
curl -i http://localhost:8080/healthz
```

## Test

```sh
export GOTOOLCHAIN=local
go test -count=1 ./...
```

The in-memory repository assigns deterministic IDs and has no external service dependencies.
