.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t telegram-bot-behran:v0.1 .

start-container:
	docker run --name telegram-bot -p 8080:8080 --env-file .env telegram-bot-behran:v0.1

