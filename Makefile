DEFAULT: mod run

mod:
	go mod tidy

run:
	go run cmd/main.go