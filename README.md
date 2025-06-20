# Go-Food - Microservice Food Delivery Application

Ứng dụng giao đồ ăn được xây dựng theo kiến trúc microservice sử dụng Go (Golang).

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
- Đăng ký, đăng nhập, đăng xuất người dùng
- Quản lý thông tin profile (Customer, Restaurant Owner, Delivery Driver)
- Xác thực và phân quyền (JWT, OAuth2)
- Reset password, verify email

### Restaurant Service
- Quản lý thông tin nhà hàng
- Đăng ký nhà hàng mới
- Cập nhật menu, giá cả, thời gian hoạt động
- Quản lý đánh giá và rating
- Upload hình ảnh nhà hàng và món ăn

### Menu Service
- Quản lý danh sách món ăn
- Phân loại món ăn (categories)
- Quản lý giá cả và khuyến mãi
- Tìm kiếm và lọc món ăn
- Quản lý tính khả dụng của món ăn

### Order Service
- Tạo đơn hàng mới
- Quản lý trạng thái đơn hàng
- Tính toán tổng tiền (bao gồm tax, delivery fee)
- Hủy đơn hàng
- Lịch sử đặt hàng

### Payment Service
- Xử lý thanh toán (Credit Card, E-wallet, COD)
- Tích hợp payment gateway
- Quản lý refund
- Lưu trữ payment history

### Delivery Service
- Quản lý delivery drivers
- Tracking đơn hàng real-time
- Tính toán route tối ưu
- Cập nhật trạng thái giao hàng
- Ước tính thời gian giao hàng

### Notification Service
- Push notification cho mobile app
- Email notification
- SMS notification
- In-app notification

### Review Service
- Đánh giá nhà hàng và món ăn
- Đánh giá delivery service
- Quản lý comments và rating
- Báo cáo review spam/inappropriate

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
┌─────────────┬─────────────┬─────────────┬─────────────────────┐
│   User      │ Restaurant  │   Menu      │      Order          │
│  Service    │  Service    │  Service    │     Service         │
└─────────────┴─────────────┴─────────────┴─────────────────────┘
┌─────────────┬─────────────┬─────────────┬─────────────────────┐
│  Payment    │  Delivery   │Notification │     Review          │
│  Service    │  Service    │  Service    │     Service         │
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

### Restaurant APIs
```
GET    /api/v1/restaurants
POST   /api/v1/restaurants
GET    /api/v1/restaurants/{id}
PUT    /api/v1/restaurants/{id}
DELETE /api/v1/restaurants/{id}
GET    /api/v1/restaurants/search?location={lat,lng}&radius={km}
```

### Menu APIs
```
GET    /api/v1/restaurants/{restaurant_id}/menu
POST   /api/v1/restaurants/{restaurant_id}/menu/items
GET    /api/v1/menu/items/{id}
PUT    /api/v1/menu/items/{id}
DELETE /api/v1/menu/items/{id}
GET    /api/v1/menu/search?q={query}&category={category}
```

### Order APIs
```
POST   /api/v1/orders
GET    /api/v1/orders
GET    /api/v1/orders/{id}
PUT    /api/v1/orders/{id}/status
DELETE /api/v1/orders/{id}
GET    /api/v1/orders/{id}/tracking
```

### Payment APIs
```
POST   /api/v1/payments
GET    /api/v1/payments/{id}
POST   /api/v1/payments/{id}/refund
GET    /api/v1/payments/history
```

### Delivery APIs
```
GET    /api/v1/delivery/drivers/available
POST   /api/v1/delivery/assign
GET    /api/v1/delivery/{order_id}/tracking
PUT    /api/v1/delivery/{order_id}/status
```

## 📁 Cấu trúc thư mục

```
go-food/
├── api/                          # API Gateway & Shared API specs
│   ├── gateway/
│   ├── proto/                    # Protocol buffer definitions
│   └── openapi/                  # OpenAPI specifications
├── services/                     # Microservices
│   ├── user-service/
│   ├── restaurant-service/
│   ├── menu-service/
│   ├── order-service/
│   ├── payment-service/
│   ├── delivery-service/
│   ├── notification-service/
│   └── review-service/
├── shared/                       # Shared libraries
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
- **users**: Thông tin người dùng
- **restaurants**: Thông tin nhà hàng
- **menu_items**: Món ăn và thông tin
- **orders**: Đơn hàng
- **order_items**: Chi tiết món ăn trong đơn hàng
- **payments**: Thông tin thanh toán
- **deliveries**: Thông tin giao hàng
- **reviews**: Đánh giá và nhận xét
- **notifications**: Thông báo

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
git clone https://github.com/your-username/go-food.git
cd go-food

# Start infrastructure services
docker-compose up -d postgres redis rabbitmq

# Run individual services
make run-user-service
make run-restaurant-service
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

- [ ] Phase 1: Core services (User, Restaurant, Menu, Order)
- [ ] Phase 2: Payment integration
- [ ] Phase 3: Real-time tracking và delivery
- [ ] Phase 4: Advanced features (recommendations, analytics)
- [ ] Phase 5: Mobile app development

## 🤝 Contributing

1. Fork repository
2. Tạo feature branch
3. Commit changes
4. Push to branch  
5. Tạo Pull Request

## 📄 License

This project is licensed under the MIT License.