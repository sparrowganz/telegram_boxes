FROM golang:latest as builder

WORKDIR /go/src/telegram_botes/services/core

COPY . .

RUN go get && CGO_ENABLED=0 GO111MODULE=on go build -v -o core .


FROM alpine:latest

WORKDIR /home

COPY --from=builder /go/src/telegram_botes/services/core/core .

CMD ["./core"]