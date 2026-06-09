.PHONY: build run test docker-build up down logs

build:
	go build -o bin/vigia ./cmd/vigia

run:
	go run ./cmd/vigia

test:
	go test ./...

docker-build:
	docker build -t vigia .

up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f vigia
