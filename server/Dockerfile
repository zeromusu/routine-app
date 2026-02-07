# Build Stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the application
RUN go build -o routine-app-server ./cmd/routine-app-server/main.go

# Run Stage
FROM alpine:latest
WORKDIR /app
# Copy the built binary from the builder stage
COPY --from=builder /app/routine-app-server .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./routine-app-server"]