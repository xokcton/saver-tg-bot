FROM golang:1.21.1-alpine3.17 AS builder

COPY . /github.com/xokcton/saver-tg-bot/
WORKDIR /github.com/xokcton/saver-tg-bot/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/bot ./cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/xokcton/saver-tg-bot/.bin/bot .
COPY --from=0 /github.com/xokcton/saver-tg-bot/configs configs/

EXPOSE 80

CMD ["./bot"]