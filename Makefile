include .env
export


.PHONY: build run test test-race lint docker-up docker-down docker-build clean

APP_NAME=kv-server
CMD_PATH=./cmd/server
BIN_PATH=./bin/$(APP_NAME)

build:
	@go build -o $(BIN_PATH) $(CMD_PATH)

run:
	@go run $(CMD_PATH)

test:
	@go test ./...

test-race:
	@go test ./... -race

clean:
	rm -rf ./bin


docker-build:
	docker compose build

docker-up:
	docker compose up

docker-down:
	docker compose down
