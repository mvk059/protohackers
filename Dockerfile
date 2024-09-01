FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o tcp-server main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/tcp-server ./

EXPOSE 10000

CMD ["./tcp-server"]