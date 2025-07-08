# ===================================================================
# STAGE 1: Build a Go binary (Giai đoạn "builder")
# ===================================================================
FROM golang:1.24.2-alpine AS builder

# Cài đặt các gói cần thiết
RUN apk add --no-cache git ca-certificates

# Đặt thư mục làm việc là thư mục gốc của service cụ thể này
# <--- THAY ĐỔI QUAN TRỌNG: Đây là thay đổi cốt lõi
WORKDIR /app/internal/services/product-service

# ---- Tối ưu hóa caching của Docker ----
# Sao chép các file dependency của product-service VÀ package dùng chung
# Docker sẽ chỉ build lại lớp này nếu các file go.mod/sum thay đổi
COPY internal/pkg/go.mod internal/pkg/go.sum /app/internal/pkg/
COPY internal/services/product-service/go.mod internal/services/product-service/go.sum ./

# Tải dependencies. Lệnh này giờ được chạy trong đúng thư mục của module.
RUN go mod download

# Sao chép TOÀN BỘ mã nguồn của dự án vào.
# Lý do cần copy toàn bộ là vì product-service có `replace` và `import` các package khác
# trong cùng một monorepo (như `internal/pkg`).
COPY . /app

# Biên dịch ứng dụng Go.
# Lệnh `go build` giờ đây rất đơn giản vì chúng ta đã ở đúng thư mục.
# -o /app/product-service: Lưu file binary ra thư mục gốc /app để stage 2 dễ dàng lấy nó.
# <--- THAY ĐỔI QUAN TRỌNG: Đường dẫn build đơn giản hơn nhiều
RUN CGO_ENABLED=0 go build -o /app/product-service -ldflags="-w -s" ./cmd/main.go

# ===================================================================
# STAGE 2: Create a minimal production image (Giai đoạn "production")
# ===================================================================
FROM alpine:latest

# Cài đặt các gói cần thiết cho production
RUN apk --no-cache add ca-certificates curl

# Tạo user không phải root để chạy ứng dụng (bảo mật hơn)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

# Sao chép file binary đã được build từ giai đoạn `builder`
# <--- THAY ĐỔI QUAN TRỌNG: Đường dẫn nguồn trong `builder` đã thay đổi
COPY --from=builder /app/product-service .

# Thay đổi quyền sở hữu
RUN chown -R appuser:appgroup /home/appuser

# Chỉ định user để chạy container
USER appuser

# Expose port
EXPOSE 8082

# HEALTHCHECK để Docker biết tình trạng "sức khỏe" của container.
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:8082/health || exit 1

# Lệnh chạy ứng dụng khi container khởi động
CMD ["./product-service"]