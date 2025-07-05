# Auto Swagger Generation

Các service hiện tại đã được cập nhật để tự động generate swagger documentation mỗi khi khởi động service.

## Cách hoạt động

1. **User Service** (`internal/services/user-service/cmd/main.go`):
   - Tự động generate swagger docs khi service khởi động
   - Swagger docs được tạo trong folder `docs/`
   - Truy cập swagger UI tại: `http://localhost:8081/swagger/index.html`

2. **Shop Service** (`internal/services/shop-service/cmd/main.go`):
   - Tự động generate swagger docs khi service khởi động  
   - Swagger docs được tạo trong folder `docs/`
   - Truy cập swagger UI tại: `http://localhost:8082/swagger/index.html`

## Yêu cầu

- Cần cài đặt `swag` tool:
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

## Thay đổi

### ✅ Đã thực hiện:
1. **Cập nhật main.go** của cả user-service và shop-service:
   - Thêm function `generateSwaggerDocs()` để tự động generate docs
   - Gọi function này khi service khởi động

2. **Cập nhật go.mod**:
   - Thêm `github.com/swaggo/swag v1.16.4` vào dependencies

3. **Makefile**:
   - Không có lệnh swagger nào cần xóa (đã kiểm tra)

### 🔄 Workflow mới:
1. Khi chạy service, swagger docs sẽ tự động được generate
2. Không cần chạy lệnh `swag init` thủ công nữa
3. Service sẽ tiếp tục hoạt động ngay cả khi swagger generation thất bại

### ⚠️ Lưu ý:
- Service sẽ log warning nếu không thể generate swagger docs
- Cần đảm bảo `swag` tool đã được cài đặt trong môi trường production
- Swagger docs được generate trong working directory của service

## Test

Để test swagger generation, chạy một trong các service:

```bash
# User service
cd internal/services/user-service
go run cmd/main.go

# Shop service  
cd internal/services/shop-service
go run cmd/main.go
```

Kiểm tra log để xem swagger generation có thành công không.
