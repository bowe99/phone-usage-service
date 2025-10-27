.PHONY: run build docker-up docker-down docker-build docker-logs download-deps

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

docker-logs:
	docker-compose logs -f

download-deps:
	go mod download
	go mod tidy