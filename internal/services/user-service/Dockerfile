# Build stage
FROM golang:1.24.2-alpine AS builder

# Install essential packages and clean cache
RUN apk update && apk upgrade && \
    apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

# Create appuser with specific UID/GID
RUN addgroup -g 10001 appgroup && \
    adduser -D -u 10001 -G appgroup -s /bin/sh appuser

WORKDIR /build

# Sao chép toàn bộ mã nguồn trước khi tải dependencies
# Cần thiết vì có replace directives trong go.mod
COPY . .

# Di chuyển vào thư mục service và tải dependencies
WORKDIR /build/internal/services/user-service
RUN go mod download

# --- Build ---
# Di chuyển lại vào thư mục service để build
WORKDIR /build/internal/services/user-service

# Build the binary with security optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -trimpath \
    -o user-service ./cmd/main.go

# Final stage - Use Alpine instead of distroless for better compatibility
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory and change ownership
WORKDIR /app
RUN chown appuser:appgroup /app

# Copy the binary
COPY --from=builder /build/internal/services/user-service/user-service /app/user-service

# Copy templates if they exist
COPY --from=builder /build/internal/services/user-service/templates /app/templates

# Change ownership of copied files
RUN chown -R appuser:appgroup /app

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/user-service"]