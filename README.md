# Go-Shop - Microservice E-commerce Delivery Platform

Nền tảng thương mại điện tử và giao hàng được xây dựng theo kiến trúc microservice sử dụng Go (Golang), tương tự như Shopee.

## 📋 Mục lục

- [Yêu cầu chức năng](#yêu-cầu-chức-năng)
- [Thiết kế API](#thiết-kế-api)
- [Tech Stack](#tech-stack)
- [Microservices](#microservices)
- [Development](#development)
- [Deployment](#deployment)

## 🎯 Yêu cầu chức năng

### User Management Service
- Đăng ký, đăng nhập (JWT)
- Đổi mật khẩu 
- Đăng xuất 
- Phân quyền (Admin, Seller, Customer, Shipper)
- CRUD profile
- Địa chỉ giao hàng (Nhiều địa chỉ)


### Shop Service
- CRUD shop
- Xử lý đơn hàng và order fulfillment
- Báo cáo doanh thu và analytics theo shop 
- Quản lý khuyến mãi, tạo discount campaigns 

### Product Service
- CRUD product

### Cart Service
- Update (Thêm, sửa, xóa sản phẩm) giỏ hàng
- Apply Promotion
- Tính tổng tiền 

### Order Service
- Tạo đơn hàng mới từ giỏ hàng 
- Lấy thông tin giỏ hàng
- Quản lý trạng thái đơn hàng (pending, confirmed, shipped, delivered, cancelled)

### Payment Service
- Xử lý thanh toán (E-wallet)
- Tích hợp payment gateway (Momo)
- Payment history

### Review Service
- Đánh giá sản phẩm và vendor/shop
- Đánh giá delivery service
- Upload hình ảnh và video review
- Quản lý comments và rating
- Verified purchase reviews

### Search & Recommendation Service
- Gợi ý sản phẩm khi người dùng đăng nhập (tương tác hành vi người dùng)
- Advanced search với filters
- Auto-complete và search suggestions
- Personalized recommendations
- Recently viewed products
- Trending products 
- Price comparison và similar products

## 🔗 Thiết kế API

### Authentication APIs
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/refresh
POST   /api/v1/auth/change-password
```

### User Management APIs
```
POST   /api/v1/users/profile
GET    /api/v1/users/profile
GET    /api/v1/users/{id}
PUT    /api/v1/users/profile
DELETE /api/v1/users

# Address Management
GET    /api/v1/users/addresses
POST   /api/v1/users/addresses
GET    /api/v1/users/addresses/{id}
PUT    /api/v1/users/addresses/{id}
DELETE /api/v1/users/addresses/{id}
PUT    /api/v1/users/addresses/{id}/default

# Role Management (User, Shipper)
POST   /api/v1/users/shippers/register
GET    /api/v1/users/shippers/profile
GET    /api/v1/users/shippers/{id}/profile
PUT    /api/v1/users/shippers/profile
```

### Shop Management APIs
```
# Shop CRUD
GET    /api/v1/shops
POST   /api/v1/shops
GET    /api/v1/shops/{id}
PUT    /api/v1/shops/{id}
DELETE /api/v1/shops/{id} - SHOP_OWNER_ONLY

# Shop Orders & Fulfillment
GET    /api/v1/shops/{id}/orders
PUT    /api/v1/shops/{id}/orders/{order_id}/status
POST   /api/v1/shops/{id}/orders/{order_id}/fulfill
GET    /api/v1/shops/{id}/orders/pending

# Shop Analytics & Reports
GET    /api/v1/shops/{id}/analytics/revenue
GET    /api/v1/shops/{id}/analytics/orders
GET    /api/v1/shops/{id}/analytics/products

# Promotions & Campaigns
GET    /api/v1/shops/{id}/promotions
POST   /api/v1/shops/{id}/promotions
GET    /api/v1/shops/{id}/promotions/{promo_id}
PUT    /api/v1/shops/{id}/promotions/{promo_id}
DELETE /api/v1/shops/{id}/promotions/{promo_id}
```

### Product Catalog APIs
```
# Product CRUD
GET    /api/v1/products
POST   /api/v1/products
GET    /api/v1/products/{id}
PUT    /api/v1/products/{id}
DELETE /api/v1/products/{id}
```

### Cart APIs
```
GET    /api/v1/cart
DELETE /api/v1/cart
POST   /api/v1/cart/items
POST   /api/v1/cart/apply-promotion
DELETE /api/v1/cart/remove-promotion
```

### Order Management APIs
```
# Order Creation & Management
POST   /api/v1/orders
GET    /api/v1/orders
GET    /api/v1/orders/{id}

# Order Status Management
GET    /api/v1/orders/{id}/status-history
PUT    /api/v1/orders/{id}/confirm
PUT    /api/v1/orders/{id}/ship
PUT    /api/v1/orders/{id}/deliver
PUT    /api/v1/orders/{id}/cancel

# Order Calculations
POST   /api/v1/orders/calculate-preview
```

### Payment APIs
```
# Payment Processing
POST   /api/v1/payments/initiate

# Payment Gateway Integration
POST   /api/v1/payments/ipn/:provider

# Transaction History
GET    /api/v1/payments/history
GET    /api/v1/payments/receipts/{id}
```

### Search & Recommendation APIs
```
# Search
GET    /api/v1/search?q={query}                   // Tìm kiếm cơ bản
GET    /api/v1/search/suggestions?q={query}       // Gợi ý khi người dùng gõ

# Personalized Recommendations
GET    /api/v1/recommendations/products           // Gợi ý sản phẩm (theo hành vi đơn giản)

# Trending & Popular
GET    /api/v1/trending/products                  // Sản phẩm đang hot
```

### Review & Rating APIs
```
# Product Reviews
GET    /api/v1/products/{id}/reviews      // Lấy danh sách review theo sản phẩm
POST   /api/v1/products/{id}/reviews      // Gửi đánh giá mới

# Shop Reviews
GET    /api/v1/shops/{id}/reviews         // Lấy danh sách review theo shop
POST   /api/v1/shops/{id}/reviews         // Gửi đánh giá mới cho shop
GET    /api/v1/shops/{id}/rating-summary  // Tóm tắt đánh giá (số sao trung bình)
```

## 🛠️ Tech Stack

### Backend
- **Language**: Go (Golang) 1.24+
- **Framework**: Gin/Echo for HTTP, gRPC for inter-service communication
- **Database**: PostgreSQL (primary), MongoDB (logs), Redis (cache)
- **Message Broker**: RabbitMQ/Apache Kafka
- **Search**: Elasticsearch
- **Authentication**: JWT, OAuth2

### Infrastructure
- **Containerization**: Docker, Docker Compose
- **Orchestration**: Kubernetes
- **API Gateway**: Kong/Nginx
- **Monitoring**: Prometheus, Grafana, Jaeger
- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **CI/CD**: GitHub Actions, Jenkins

### Development Tools
- **Package Manager**: Go Modules
- **Testing**: Testify, GoMock
- **Code Quality**: golangci-lint, SonarQube
- **Documentation**: Swagger/OpenAPI

### Data Management
- **Database per Service**: Mỗi service sở hữu data riêng
- **Event-driven Architecture**: Services giao tiếp qua events
- **CQRS**: Tách read/write models cho complex queries
- **Saga Pattern**: Distributed transaction management
- **Data Consistency**: Eventually consistent với compensation patterns

### Service Discovery & Load Balancing
- **Service Registry**: Consul/Etcd cho service registration
- **Load Balancer**: Nginx/HAProxy cho traffic distribution
- **Health Checks**: Automatic service health monitoring
- **Circuit Breaker**: Fault tolerance và resilience patterns
- **Rate Limiting**: API throttling và abuse prevention

### Security & Cross-cutting Concerns
- **API Gateway**: Kong/Nginx cho unified entry point
- **Authentication**: JWT token validation across services
- **Authorization**: Role-based access control (RBAC)
- **Audit Logging**: Distributed tracing với Jaeger
- **Monitoring**: Prometheus metrics với Grafana dashboards

## 🚀 Development

### Yêu cầu hệ thống
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+
- RabbitMQ 3.12+

### Setup môi trường development
```bash
# Clone repository
git clone https://github.com/toji-dev/go-shop.git
cd go-shop

# Start infrastructure services
docker-compose up -d postgres redis rabbitmq elasticsearch

# Run individual services
make run-user-service
make run-vendor-service
make run-product-service
make run-cart-service
# ...
```

## 🚀 Deployment

### Docker Deployment
```bash
# Build all services
make build-all

# Deploy with docker-compose
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes Deployment
```bash
# Apply Kubernetes manifests
kubectl apply -f deployments/kubernetes/
```

### Environment Variables
- Tạo file `.env` cho mỗi service
- Sử dụng Kubernetes secrets cho production
- Configure external services (databases, message queues)

---

## 📈 Roadmap

- [ ] Phase 1: Core services (User, Vendor, Product, Cart, Order)
- [ ] Phase 2: Payment integration và escrow service
- [ ] Phase 3: Search & recommendation engine
- [ ] Phase 4: Advanced features (live chat, flash sales, affiliate program)
- [ ] Phase 5: International expansion features

## 📄 License

This project is licensed under the MIT License.