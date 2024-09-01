FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 10000

CMD ["./main"]