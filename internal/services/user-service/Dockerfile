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

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary with security optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -trimpath \
    -o user-service ./cmd/main.go

# Final stage - Use distroless instead of scratch
FROM gcr.io/distroless/static-debian12:nonroot

# Copy the binary
COPY --from=builder /build/user-service /app/user-service

# Copy templates if they exist
COPY --from=builder /build/templates /app/templates

# Set working directory
WORKDIR /app

# Use nonroot user (UID 65532)
USER 65532:65532

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/user-service"]