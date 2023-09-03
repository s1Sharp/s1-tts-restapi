.PHONY: dc run test lint

dc:
	@docker-compose up  --remove-orphans --build

run:
	@go build -o app cmd/s1-tts-restapi/main.go && HTTP_ADDR=:3000 ./app

test:
	@go test -race ./...

lint:
	@golangci-lint run

env: app.env
	@echo ok