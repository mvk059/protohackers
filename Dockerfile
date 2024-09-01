FROM golang:1.20-alpine AS builder

WORKDIR /
COPY . .
RUN go build -o main .

FROM alpine:latest

WORKDIR /
COPY --from=builder /main .

EXPOSE 10000

CMD ["./main"]