.PHONY: run build docker-up docker-down docker-build docker-logs download-deps

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-unit:
	go test -v -race ./test/unit/...

test-integration:
	go test -v -race ./test/integration/...

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