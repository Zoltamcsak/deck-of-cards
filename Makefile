PKG_MAIN=./cmd/main.go
BINARY_NAME=deck-of-cards

.PHONY: run
run:
	@go run $(PKG_MAIN)

.PHONY: build
build:
	go build -o $(BINARY_NAME) $(PKG_MAIN)

.PHONY: test
test:
	go test ./...

.PHONY: run-db
run-db:
	docker-compose up -d

.PHONY: stop-db
stop-db:
	docker-compose down

.PHONY: tidy
tidy:
	go mod tidy & go mod vendor