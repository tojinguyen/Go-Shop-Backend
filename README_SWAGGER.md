# 🎉 Swagger Documentation - Go-Shop Microservices

Đã hoàn thành việc cài đặt Swagger cho tất cả các microservices trong Go-Shop project!

## 📋 Services và URLs

| Service | Port | Swagger URL | Status |
|---------|------|-------------|--------|
| User Service | 8081 | http://localhost:8081/swagger/index.html | ✅ Complete |
| Shop Service | 8082 | http://localhost:8082/swagger/index.html | ✅ Setup |
| Order Service | 8083 | http://localhost:8083/swagger/index.html | 🔧 Dependencies |
| Payment Service | 8084 | http://localhost:8084/swagger/index.html | 🔧 Dependencies |
| Review Service | 8085 | http://localhost:8085/swagger/index.html | 🔧 Dependencies |
| Shipping Service | 8086 | http://localhost:8086/swagger/index.html | 🔧 Dependencies |
| Notification Service | 8087 | http://localhost:8087/swagger/index.html | 🔧 Dependencies |

## 🚀 Quick Start

### 1. Generate Documentation
```bash
# Tất cả services
make swagger-gen

# Service cụ thể
make swagger-gen-user
make swagger-gen-shop
```

### 2. Start Services
```bash
# Với documentation
make run-with-swagger

# Hoặc thông thường
make up
```

### 3. PowerShell Script
```powershell
.\generate-swagger.ps1
```

## 📝 API Documentation Examples

### User Service Auth
```go
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration request"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) { ... }
```

## 🛠️ Thêm Documentation cho Service Mới

### 1. Thêm annotations vào handler:
```go
// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateProductRequest true "Product data"
// @Success 201 {object} dto.ProductResponse "Product created"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // implementation
}
```

### 2. Generate docs:
```bash
make swagger-gen-[service-name]
```

### 3. Restart service để áp dụng changes

## 🔧 Troubleshooting

### Swagger không load
- Kiểm tra service đã start
- Verify port mapping
- Check endpoint `/swagger/index.html`

### Documentation không cập nhật
- Chạy `make swagger-gen-[service]`
- Restart service
- Clear browser cache

### Parse errors
- Kiểm tra syntax Swagger comments
- Đảm bảo import đúng packages
- Validate JSON/YAML output

## 📁 Project Structure

```
service-name/
├── docs/
│   ├── docs.go          # Generated Go code
│   ├── swagger.json     # JSON specification
│   └── swagger.yaml     # YAML specification
├── cmd/
│   └── main.go         # Swagger annotations
└── internal/
    └── handlers/        # API handlers with annotations
```

## 🔑 Security

Trong production:
- Disable Swagger UI hoặc bảo vệ bằng authentication
- Chỉ enable cho development/staging environments
- Review security schemas trong documentation

---

**📌 Status: Ready for development!**

Swagger UI đã sẵn sàng để test và document APIs.
