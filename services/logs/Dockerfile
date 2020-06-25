FROM golang:latest as builder

WORKDIR /go/src/telegram_botes/services/logs

COPY . .

RUN CGO_ENABLED=0 GO111MODULE=on go build -mod vendor -v -o logs_service .


FROM alpine:latest

WORKDIR /home

COPY --from=builder /go/src/telegram_botes/services/logs/logs_service .

CMD ["./logs_service"]