.PHONY: dc build run test lint mock env

dc:
	@docker-compose up  --remove-orphans --build

build:
	@go build -v -o app cmd/s1-tts-restapi/*.go

run:
	@go build -o app cmd/s1-tts-restapi/main.go && HTTP_ADDR=:3000 ./app

test:
	@go test -race -v ./...

lint:
	@golangci-lint run

mock:
	@mockgen -package mock_service -destination service/mock/user.service.mock.go github.com/s1Sharp/s1-tts-restapi/service UserService

env: .env
	@echo ok


