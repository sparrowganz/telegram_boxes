FROM golang:latest as builder

WORKDIR /go/src/telegram_botes/services/admin

COPY . .

RUN go get && CGO_ENABLED=0 GO111MODULE=on go build -v -o admin_bot .


FROM alpine:latest

WORKDIR /home

COPY --from=builder /go/src/telegram_botes/services/admin/admin_bot .

CMD ["./admin_bot"]