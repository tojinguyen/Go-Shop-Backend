# CreateShop Feature - Vertical Slice Architecture

## Overview
Feature CreateShop được implement theo kiến trúc VSA (Vertical Slice Architecture), trong đó tất cả logic liên quan đến việc tạo shop được tổ chức trong một slice độc lập.

## Structure
```
features/shop/commands/create_shop/
├── command.go      # Command definition
├── handler.go      # Business logic handler
├── request.go      # HTTP request DTO
└── api_handler.go  # HTTP API handler
```

## Flow
1. **HTTP Request** → `api_handler.go`
2. **Request Validation** → `request.go` 
3. **Business Logic** → `handler.go`
4. **Database Operations** → `repository/postgres_shop_repository.go`
5. **Response** → `api_handler.go`

## Key Components

### Command (command.go)
- Định nghĩa CreateShopCommand struct
- Chứa tất cả data cần thiết để tạo shop

### Handler (handler.go)
- Chứa toàn bộ business logic:
  - Validation owner không có shop trước đó
  - Validation email không bị trùng
  - Tạo domain object
  - Lưu vào database
- Return CreateShopResponse

### Request (request.go)
- HTTP request validation
- Convert HTTP request to Command
- Validation rules bằng gin binding tags

### API Handler (api_handler.go)
- Handle HTTP request/response
- Error handling và status codes
- Integration với Gin framework

## Benefits của VSA
1. **High Cohesion**: Tất cả logic liên quan đến CreateShop ở cùng một nơi
2. **Low Coupling**: Feature độc lập, không phụ thuộc vào service layer
3. **Easy Testing**: Có thể test từng feature riêng biệt
4. **Maintainability**: Dễ maintain và extend feature
5. **Clear Boundaries**: Feature boundaries rõ ràng

## Usage

### HTTP Request
```bash
POST /api/v1/shops
Content-Type: application/json

{
  "owner_id": "123e4567-e89b-12d3-a456-426614174000",
  "shop_name": "My Shop",
  "avatar_url": "https://example.com/avatar.jpg",
  "banner_url": "https://example.com/banner.jpg",
  "shop_description": "My awesome shop",
  "address_id": "123e4567-e89b-12d3-a456-426614174001",
  "phone": "0123456789",
  "email": "shop@example.com"
}
```

### Response
```json
{
  "success": true,
  "message": "Created",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174002",
    "owner_id": "123e4567-e89b-12d3-a456-426614174000",
    "shop_name": "My Shop",
    "avatar_url": "https://example.com/avatar.jpg",
    "banner_url": "https://example.com/banner.jpg",
    "shop_description": "My awesome shop",
    "address_id": "123e4567-e89b-12d3-a456-426614174001",
    "phone": "0123456789",
    "email": "shop@example.com",
    "rating": 0.0,
    "status": "inactive",
    "created_at": "2025-07-05T10:30:00Z"
  }
}
```

## Error Handling
- **400 Bad Request**: Invalid request data
- **409 Conflict**: Owner already has shop hoặc email already in use  
- **500 Internal Server Error**: Database errors

## Dependencies
- Repository: `repository.ShopRepository`
- Domain: `domain.Shop`
- Infrastructure: PostgreSQL database
