
.PHONY: run_tests

run_tests:
	@go test ./... -race
run_serv:
	@go run ./cmd/server  