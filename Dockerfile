FROM golang:alpine as builder
WORKDIR /go/src/telegram-bot
COPY . /go/src/telegram-bot
RUN go build -mod=vendor -o ./dist/ems

FROM alpine:3.11.3
RUN apk add --update ca-certificates
RUN apk add --no-cache tzdata && \
  cp -f /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
  apk del tzdata

COPY ./config/config.yaml .
COPY --from=builder /go/src/telegram-bot/dist/ems .
EXPOSE 9090
ENTRYPOINT ["./ems"]
