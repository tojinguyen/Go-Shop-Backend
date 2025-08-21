# Build stage
FROM golang:1.24-alpine AS builder

# Install essential packages and clean cache
RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.work go.work.sum ./

COPY internal/pkg/go.mod internal/pkg/go.sum ./internal/pkg/
COPY internal/services/user-service/go.mod internal/services/user-service/go.sum ./internal/services/user-service/
COPY internal/services/shop-service/go.mod internal/services/shop-service/go.sum ./internal/services/shop-service/
COPY internal/services/product-service/go.mod internal/services/product-service/go.sum ./internal/services/product-service/
COPY internal/services/cart-service/go.mod internal/services/cart-service/go.sum ./internal/services/cart-service/
COPY internal/services/order-service/go.mod internal/services/order-service/go.sum ./internal/services/order-service/
COPY internal/services/payment-service/go.mod internal/services/payment-service/go.sum ./internal/services/payment-service/
COPY proto/gen/go/go.mod proto/gen/go/go.sum ./proto/gen/go/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -ldflags='-w -s' \
    -trimpath \
    -o /app/payment-service-binary ./internal/services/payment-service/cmd/main.go

# Final stage - Use Alpine for better compatibility
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory and change ownership
WORKDIR /home/appuser

COPY --from=builder /app/payment-service-binary .

RUN chown -R appuser:appgroup /home/appuser

# Chuyển sang user mới
USER appuser

# Expose port
EXPOSE 8085

# Lệnh chạy ứng dụng
ENTRYPOINT ["./payment-service-binary"]
