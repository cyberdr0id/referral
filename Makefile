run:
	go run cmd/main.go

lint:
	golangci-lint run --config .golangci.yml

dc-up:
	docker-compose up

.PHONY: run lint dc-up