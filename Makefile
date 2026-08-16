export GOTOOLCHAIN := local

.PHONY: run test build

run:
	go run ./cmd/aftercare

test:
	go test -count=1 ./...

build:
	CGO_ENABLED=0 go build ./...
