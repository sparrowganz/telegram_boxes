FROM golang:latest as builder

WORKDIR /go/src/telegram_botes/services/box

COPY . .

RUN  CGO_ENABLED=0 GO111MODULE=on go build -mod=vendor -v -o box_bot .


FROM alpine:latest

WORKDIR /home

COPY --from=builder /go/src/telegram_botes/services/box/box_bot .

CMD ["./box_bot"]