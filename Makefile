.PHONY:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t xokcton/saver-tg-bot:0.1 .

start-container:
	docker run --env-file .env -p 80:80 xokcton/saver-tg-bot:0.1