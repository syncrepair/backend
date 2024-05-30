all: mod-tidy run

mod-tidy:
	go mod tidy

run:
	go run cmd/main.go