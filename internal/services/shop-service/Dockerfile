# Build stage 
FROM golang:1.23-alpine AS builder

# Install essential packages and clean cache
RUN apk update && apk upgrade && \
    apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Production stage
FROM alpine:latest

# Install ca-certificates and curl for HTTPS requests and health checks
RUN apk --no-cache add ca-certificates tzdata curl

# Create a non-root user
RUN adduser -D -s /bin/sh shopuser

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy any configuration files if needed
COPY --from=builder /app/internal/config/ ./internal/config/

# Change ownership to non-root user
RUN chown -R shopuser:shopuser /root
USER shopuser

# Expose port (default 8080 or from config)
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Command to run the application
CMD ["./main"]

