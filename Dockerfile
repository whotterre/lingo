# Build stage
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache git && cd cmd/api && go build -o main .

# Runtime stage
FROM alpine:latest
WORKDIR /root/
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main .
EXPOSE 7000
CMD ["./main"]
