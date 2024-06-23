all: swag mod-tidy run

mod-tidy:
	go mod tidy

run:
	go run cmd/app/main.go

swag:
	swag init -g cmd/app/main.go