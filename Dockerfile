FROM golang:1.15-alpine3.12 AS builder

COPY . /pocket-bot/
WORKDIR /pocket-bot/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /pocket-bot/bin/bot .
COPY --from=0 /pocket-bot/configs configs/

EXPOSE 8080

CMD ["./bot"]

