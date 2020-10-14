FROM golang:alpine as builder
WORKDIR /go/src/transport/ems
COPY . /go/src/transport/ems
RUN go build -mod=vendor -o ./dist/ems

FROM alpine:3.11.3
RUN apk add --update ca-certificates
RUN apk add --no-cache tzdata && \
  cp -f /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
  apk del tzdata

COPY ./config/config.yaml .
COPY --from=builder /go/src/transport/ems/dist/ems .
EXPOSE 9090
ENTRYPOINT ["./ems"]
