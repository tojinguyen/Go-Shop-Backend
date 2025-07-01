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

# --- Tối ưu hóa Caching ---
# 1. Sao chép các file go.mod/go.sum để tải dependencies
# Chỉ những thay đổi trong các file này mới làm cho bước `go mod download` chạy lại.

# Sao chép go.mod/sum của user-service
COPY internal/services/user-service/go.mod internal/services/user-service/go.sum ./internal/services/user-service/

# Sao chép go.mod/sum của các package dùng chung (nếu có)
COPY internal/pkg/go.mod internal/pkg/go.sum ./internal/pkg/

# Di chuyển vào thư mục service và tải dependencies
WORKDIR /build/internal/services/user-service
RUN go mod download

# 2. Sao chép toàn bộ mã nguồn
# Bước này sẽ chạy lại khi có bất kỳ thay đổi nào trong mã nguồn,
# nhưng bước tải dependencies ở trên sẽ được lấy từ cache.
WORKDIR /build
COPY . .

# --- Build ---
# Di chuyển lại vào thư mục service để build
WORKDIR /build/internal/services/user-service

# Build the binary with security optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -trimpath \
    -o user-service ./cmd/main.go

# Final stage - Use distroless instead of scratch
FROM gcr.io/distroless/static-debian12:nonroot

# Copy the binary
COPY --from=builder /build/internal/services/user-service/user-service /app/user-service

# Copy templates if they exist
COPY --from=builder /build/internal/services/user-service/templates /app/templates

# Set working directory
WORKDIR /app

# Use nonroot user (UID 65532)
USER 65532:65532

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/user-service"]