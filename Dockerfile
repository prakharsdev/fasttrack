# Build stage (Debian-based, more stable than Alpine)
FROM golang:1.24-bookworm AS builder

# Set working directory
WORKDIR /app

# Install git (required for go mod if private repos used)
RUN apt-get update && apt-get install -y git

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code
COPY . .

# Build the Go binary
RUN go build -o fasttrack cmd/app/main.go

# Runtime stage (smaller)
FROM debian:bookworm-slim

# Set working directory in container
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/fasttrack .

# Copy .env file to container (optional — if you want .env inside container)
COPY .env .

# Expose health check port
EXPOSE 8080

# Run the Go binary
CMD ["./fasttrack"]
