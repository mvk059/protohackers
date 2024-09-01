FROM golang:1.21-alpine AS builder

WORKDIR /

# Copy the source code
COPY main.go .

# List contents of the directory (for debugging)
RUN ls -la

# Build the application
RUN go build -o protohackers main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /protohackers .

# We don't EXPOSE any specific port as it's set via environment variable

CMD ["./protohackers"]