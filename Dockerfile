FROM golang:1.21-alpine AS builder

WORKDIR /

# Copy the source code
COPY main.go .

# Build the application
RUN go build -o protohackers

FROM alpine:latest

WORKDIR /

COPY --from=builder /protohackers .

# We don't EXPOSE any specific port as it's set via environment variable

CMD ["./protohackers"]