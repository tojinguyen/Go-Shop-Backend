# Go-Shop - Microservice E-commerce Delivery Platform

Nền tảng thương mại điện tử và giao hàng được xây dựng theo kiến trúc microservice sử dụng Go (Golang), tương tự như Shopee.

## 📋 Mục lục

- [Yêu cầu chức năng](#yêu-cầu-chức-năng)
- [Yêu cầu phi chức năng](#yêu-cầu-phi-chức-năng)
- [Kiến trúc hệ thống](#kiến-trúc-hệ-thống)
- [Thiết kế API](#thiết-kế-api)
- [Cấu trúc thư mục](#cấu-trúc-thư-mục)
- [Tech Stack](#tech-stack)
- [Database Schema](#database-schema)
- [Microservices](#microservices)
- [Development](#development)
- [Deployment](#deployment)

## 🎯 Yêu cầu chức năng

### User Management Service
- Đăng ký, đăng nhập, đăng xuất người dùng (Buyer, Seller, Admin)
- Quản lý thông tin profile (Customer, Vendor, Delivery Partner)
- Xác thực và phân quyền (JWT, OAuth2)
- Reset password, verify email
- Quản lý địa chỉ giao hàng multiple addresses

### Vendor/Seller Service
- Đăng ký shop/store mới
- Quản lý thông tin shop (tên, mô tả, logo, banner)
- Quản lý sản phẩm và inventory
- Xử lý đơn hàng và order fulfillment
- Báo cáo doanh thu và analytics
- Upload hình ảnh sản phẩm và shop

### Product Service
- Quản lý catalog sản phẩm
- Phân loại sản phẩm theo categories/subcategories
- Quản lý giá cả, khuyến mãi và discount campaigns
- Tìm kiếm và lọc sản phẩm (price, rating, location, category)
- Quản lý stock và inventory
- Product variations (size, color, model)
- Bulk import/export sản phẩm

### Shopping Cart Service
- Quản lý giỏ hàng của user
- Add/remove/update items
- Calculate total với taxes và shipping
- Save for later functionality
- Cross-selling suggestions

### Order Service
- Tạo đơn hàng mới từ multiple vendors
- Quản lý trạng thái đơn hàng (pending, confirmed, shipped, delivered, cancelled)
- Tính toán tổng tiền (product price, shipping fee, taxes, discount)
- Hủy đơn hàng và return/refund processing
- Order splitting theo vendor
- Lịch sử mua hàng

### Payment Service
- Xử lý thanh toán (Credit Card, E-wallet, Bank Transfer, COD)
- Tích hợp payment gateway (Stripe, PayPal, VNPay, Momo)
- Quản lý refund và chargeback
- Split payment cho multiple vendors
- Escrow service cho buyer protection
- Payment history và transaction logs

### Shipping & Delivery Service
- Quản lý shipping partners và delivery methods
- Tích hợp với 3rd party logistics (Giao Hàng Nhanh, Giao Hàng Tiết Kiệm)
- Tracking đơn hàng real-time
- Tính toán shipping cost theo distance và weight
- Delivery time estimation
- Address validation và geocoding
- Proof of delivery (POD)

### Notification Service
- Push notification cho mobile app
- Email notification
- SMS notification
- In-app notification

### Review Service
- Đánh giá sản phẩm và vendor/shop
- Đánh giá delivery service
- Upload hình ảnh và video review
- Q&A section cho sản phẩm
- Quản lý comments và rating
- Báo cáo review spam/inappropriate
- Verified purchase reviews

### Search & Recommendation Service
- Advanced search với filters
- Auto-complete và search suggestions
- Personalized recommendations
- Recently viewed products
- Trending products và bestsellers
- Price comparison và similar products

## ⚡ Yêu cầu phi chức năng

### Performance
- Hỗ trợ 10,000+ concurrent users
- Response time < 200ms cho các API chính
- Database query optimization
- Caching strategy (Redis)
- Load balancing

### Scalability
- Horizontal scaling cho các microservices
- Auto-scaling based on traffic
- Database sharding/partitioning
- Message queue for async processing

### Reliability
- 99.9% uptime
- Circuit breaker pattern
- Retry mechanism
- Graceful degradation
- Health check endpoints

### Security
- Authentication & Authorization (JWT)
- Data encryption (at rest và in transit)
- API rate limiting
- Input validation và sanitization
- HTTPS only
- SQL injection prevention

### Monitoring & Logging
- Centralized logging (ELK stack)
- Metrics collection (Prometheus)
- Distributed tracing (Jaeger)
- Error tracking
- Performance monitoring

## 🏗️ Kiến trúc hệ thống

```
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway (Kong/Nginx)                 │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway (Kong/Nginx)                 │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────┬─────────────┬─────────────┬─────────────────────┐
│   User      │   Vendor/   │  Product    │   Shopping Cart     │
│  Service    │   Seller    │  Service    │     Service         │
│             │  Service    │             │                     │
└─────────────┴─────────────┴─────────────┴─────────────────────┘
┌─────────────┬─────────────┬─────────────┬─────────────────────┐
│    Order    │  Payment    │ Shipping &  │  Notification       │
│   Service   │  Service    │ Delivery    │    Service          │
│             │             │  Service    │                     │
└─────────────┴─────────────┴─────────────┴─────────────────────┘
┌─────────────┬─────────────┬─────────────┬─────────────────────┐
│   Review    │ Search &    │   Admin     │    Analytics        │
│  Service    │ Recommend   │  Service    │    Service          │
│             │  Service    │             │                     │
└─────────────┴─────────────┴─────────────┴─────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Message Broker (RabbitMQ/Kafka)               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────┬─────────────┬─────────────┬─────────────────────┐
│ PostgreSQL  │   MongoDB   │    Redis    │    Elasticsearch    │
│ (Primary)   │ (Logs/Docs) │  (Cache)    │     (Search)        │
└─────────────┴─────────────┴─────────────┴─────────────────────┘
```

## 🔗 Thiết kế API

### Authentication APIs
```
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/logout
POST /api/v1/auth/refresh
POST /api/v1/auth/forgot-password
POST /api/v1/auth/reset-password
```

### User APIs
```
GET    /api/v1/users/profile
PUT    /api/v1/users/profile
GET    /api/v1/users/{id}
DELETE /api/v1/users/{id}
```

### Vendor/Seller APIs
```
GET    /api/v1/vendors
POST   /api/v1/vendors
GET    /api/v1/vendors/{id}
PUT    /api/v1/vendors/{id}
DELETE /api/v1/vendors/{id}
GET    /api/v1/vendors/search?location={lat,lng}&category={category}
GET    /api/v1/vendors/{id}/products
GET    /api/v1/vendors/{id}/orders
GET    /api/v1/vendors/{id}/analytics
```

### Product APIs
```
GET    /api/v1/products
POST   /api/v1/products
GET    /api/v1/products/{id}
PUT    /api/v1/products/{id}
DELETE /api/v1/products/{id}
GET    /api/v1/products/search?q={query}&category={category}&price_min={min}&price_max={max}
GET    /api/v1/products/categories
GET    /api/v1/products/{id}/reviews
GET    /api/v1/products/{id}/related
```

### Shopping Cart APIs
```
GET    /api/v1/cart
POST   /api/v1/cart/items
PUT    /api/v1/cart/items/{id}
DELETE /api/v1/cart/items/{id}
DELETE /api/v1/cart/clear
GET    /api/v1/cart/summary
```

### Order APIs
```
POST   /api/v1/orders
GET    /api/v1/orders
GET    /api/v1/orders/{id}
PUT    /api/v1/orders/{id}/status
DELETE /api/v1/orders/{id}
GET    /api/v1/orders/{id}/tracking
POST   /api/v1/orders/{id}/return
POST   /api/v1/orders/{id}/cancel
```

### Payment APIs
```
POST   /api/v1/payments
GET    /api/v1/payments/{id}
POST   /api/v1/payments/{id}/refund
GET    /api/v1/payments/history
GET    /api/v1/payments/methods
POST   /api/v1/payments/escrow/release
```

### Shipping & Delivery APIs
```
GET    /api/v1/shipping/methods
POST   /api/v1/shipping/calculate
GET    /api/v1/shipping/{order_id}/tracking
PUT    /api/v1/shipping/{order_id}/status
GET    /api/v1/shipping/providers
POST   /api/v1/shipping/labels
```

### Search & Recommendation APIs
```
GET    /api/v1/search?q={query}&filters={filters}
GET    /api/v1/search/suggestions?q={partial_query}
GET    /api/v1/recommendations/products
GET    /api/v1/recommendations/vendors
GET    /api/v1/trending/products
GET    /api/v1/recent/products
```

## 📁 Cấu trúc thư mục

```
go-shop/
├── api/                          # API Gateway & Shared API specs
│   ├── gateway/
│   ├── proto/                    # Protocol buffer definitions
│   └── openapi/                  # OpenAPI specifications
├── internal/
│   └── services/                 # Microservices
│       ├── user-service/
│       ├── vendor-service/
│       ├── product-service/
│       ├── cart-service/
│       ├── order-service/
│       ├── payment-service/
│       ├── shipping-service/
│       ├── notification-service/
│       ├── review-service/
│       ├── search-service/
│       └── admin-service/
├── pkg/                          # Shared libraries
│   ├── auth/
│   ├── config/
│   ├── database/
│   ├── middleware/
│   ├── models/
│   └── utils/
├── deployments/                  # Deployment configurations
│   ├── docker/
│   ├── kubernetes/
│   └── terraform/
├── scripts/                      # Build and deployment scripts
├── docs/                         # Documentation
├── tests/                        # Integration tests
└── tools/                        # Development tools
```

## 🛠️ Tech Stack

### Backend
- **Language**: Go (Golang) 1.21+
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

## 🗄️ Database Schema

### Các bảng chính:
- **users**: Thông tin người dùng (buyers, sellers, admins)
- **vendors**: Thông tin shop/seller
- **products**: Catalog sản phẩm và thông tin chi tiết
- **product_variants**: Biến thể sản phẩm (size, color, etc.)
- **categories**: Danh mục sản phẩm
- **shopping_carts**: Giỏ hàng của user
- **cart_items**: Chi tiết sản phẩm trong giỏ hàng
- **orders**: Đơn hàng
- **order_items**: Chi tiết sản phẩm trong đơn hàng
- **payments**: Thông tin thanh toán
- **shipping**: Thông tin vận chuyển
- **reviews**: Đánh giá sản phẩm và vendor
- **notifications**: Thông báo
- **addresses**: Địa chỉ giao hàng của user

## 🎯 Microservices

### Service Discovery
- Consul/Etcd cho service registration và discovery
- Health check endpoints cho tất cả services

### Inter-service Communication
- gRPC cho internal communication
- REST API cho external clients
- Message queue cho async operations

### Data Management
- Database per service pattern
- Event-driven architecture
- CQRS pattern cho complex queries

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
git clone https://github.com/your-username/go-shop.git
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
- [ ] Phase 3: Shipping integration và tracking
- [ ] Phase 4: Search & recommendation engine
- [ ] Phase 5: Advanced features (live chat, flash sales, affiliate program)
- [ ] Phase 6: Mobile app development
- [ ] Phase 7: Seller analytics dashboard
- [ ] Phase 8: International expansion features

## 🤝 Contributing

1. Fork repository
2. Tạo feature branch
3. Commit changes
4. Push to branch  
5. Tạo Pull Request

## 📄 License

This project is licensed under the MIT License.