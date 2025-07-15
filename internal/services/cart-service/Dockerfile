FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# --- TỐI ƯU HÓA CACHE CỦA DOCKER ---
# Bước 1: Copy các file go.mod và go.sum của service vào trước.
# Điều này tận dụng Docker layer caching. Nếu các file này không đổi, Docker sẽ không cần
# tải lại dependencies ở lần build sau.
COPY internal/services/cart-service/go.mod internal/services/cart-service/go.sum ./
# Phải copy cả go.mod của các package `replace`
COPY internal/pkg/go.mod ./internal/pkg/go.mod
COPY proto/gen/go/go.mod ./proto/gen/go/go.mod

# Bước 2: Tải dependencies.
# Lệnh `go mod download` chỉ chạy khi các file mod/sum thay đổi.
RUN go mod download

# Bước 3: Copy toàn bộ source code của project vào.
# Chỉ khi source code thay đổi (chứ không phải dependencies), Docker mới build lại từ bước này.
COPY . .

# --- BUILD BINARY ---
# Di chuyển vào đúng thư mục của service để thực hiện build
WORKDIR /app/internal/services/cart-service

# Build ứng dụng Go.
# -o /app/cart-service: Lưu file binary ra thư mục gốc /app để stage 2 dễ dàng lấy nó.
RUN CGO_ENABLED=0 go build -o /app/cart-service -ldflags="-w -s" ./cmd/main.go

# ===================================================================
# STAGE 2: Create a minimal production image (Giai đoạn "production")
# ===================================================================
# Bắt đầu từ một image nền siêu nhỏ và an toàn.
FROM alpine:latest


# Cài đặt các gói cần thiết cho production
RUN apk --no-cache add ca-certificates

# Tạo user không phải root để chạy ứng dụng
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Đặt thư mục làm việc
WORKDIR /home/appuser

# Sao chép file binary đã được build từ giai đoạn `builder`
COPY --from=builder /app/cart-service .

# Thay đổi quyền sở hữu của file cho user mới
RUN chown appuser:appgroup /home/appuser/cart-service

# Chỉ định user để chạy container
USER appuser

# Expose port mà service sẽ lắng nghe
# Port này sẽ được map trong docker-compose.yml
EXPOSE 8083

# HEALTHCHECK để Docker biết tình trạng "sức khỏe" của container.
# Bạn cần tạo một endpoint /ping hoặc /health trong router của cart-service.
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD [ "wget", "-q", "-O", "/dev/null", "http://localhost:8083/ping" ] || exit 1

# Lệnh chạy ứng dụng khi container khởi động
CMD ["./cart-service"]