.PHONY: up down build test lint k6 clean

up:
	docker-compose up --build

down:
	docker-compose down -v

build:
	docker-compose build

test:
	go test -count=1 ./...

lint:
	golangci-lint run ./...

k6:
	k6 run tests/k6/reserve_test.js

clean:
	docker-compose down -v
	rm -f ./wishlist-api
	go clean -testcache