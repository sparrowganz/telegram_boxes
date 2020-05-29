FROM golang:latest as builder

WORKDIR /go/src/fileServer/services/logs

COPY . .

RUN go get && CGO_ENABLED=0 GO111MODULE=on go build -v -o logs_service .


FROM alpine:latest

WORKDIR /home

COPY --from=builder /go/src/fileServer/services/logs/logs_service .

CMD ["./logs_service"]