# Go-Shop - Microservice E-commerce Delivery Platform

Nền tảng thương mại điện tử và giao hàng được xây dựng theo kiến trúc microservice sử dụng Go (Golang), tương tự như Shopee.

## 📋 Mục lục

- [Yêu cầu chức năng](#yêu-cầu-chức-năng)
- [Yêu cầu phi chức năng](#yêu-cầu-phi-chức-năng)
- [Thiết kế API](#thiết-kế-api)
- [Kiến trúc hệ thống](#kiến-trúc-hệ-thống)
- [Cấu trúc thư mục](#cấu-trúc-thư-mục)
- [Tech Stack](#tech-stack)
- [Database Schema](#database-schema)
- [Microservices](#microservices)
- [Development](#development)
- [Deployment](#deployment)

## 🎯 Yêu cầu chức năng

### User Management Service
- Đăng ký, đăng nhập (JWT)
- Quên mật khẩu (OTP, email)
- Đổi mật khẩu 
- Đăng xuất 
- Phân quyền (User, Shipper)
- CRUD profile
- Địa chỉ giao hàng (Nhiều địa chỉ)


### Search & Recommendation Service
- Gợi ý sản phẩm khi người dùng đăng nhập (tương tác hành vi người dùng)
- Advanced search với filters
- Auto-complete và search suggestions
- Personalized recommendations
- Recently viewed products
- Trending products 
- Price comparison và similar products


### Product Catalog Service
- CRUD shop
- Quản lý thông tin profile shop 
- Báo cáo doanh thu và analytics theo shop 
- Xử lý đơn hàng và order fulfillment
- Quản lý khuyến mãi, tạo discount campaigns 

- CRUD product
- Quản lý catalog sản phẩm (title, description, media, brand, model)
- Quản lý giá cả sản phẩm 
- Quản lý stock và inventory
- Phân loại sản phẩm theo categories/subcategories
- Tìm kiếm và lọc sản phẩm (price, rating, location, category)
- Bulk import/export sản phẩm

### Shopping Cart Service
- Quản lý giỏ hàng của user
- Thêm, sửa, xóa sản phẩm
- Tính tổng tiền 
- Lưu lại giỏ hàng

### Order Service
- Tạo đơn hàng mới từ giỏ hàng 
- Quản lý trạng thái đơn hàng (pending, confirmed, shipped, delivered, cancelled)
- Tính toán tổng tiền (product price, shipping fee, taxes, discount)
- Hủy đơn hàng và return/refund processing
- Lịch sử mua hàng

### Payment Service
- Xử lý thanh toán (Credit Card, E-wallet, Bank Transfer, COD)
- Tích hợp payment gateway (Stripe, PayPal, VNPay, Momo)
- Quản lý refund và chargeback
- Escrow service cho buyer protection
- Payment history và transaction logs

### Shipping Service
- Quản lý shipper
- Tracking đơn hàng real-time
- Tính toán shipping cost theo distance và weight
- Address validation và geocoding

### Review Service
- Đánh giá sản phẩm và vendor/shop
- Đánh giá delivery service
- Upload hình ảnh và video review
- Quản lý comments và rating
- Verified purchase reviews


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

## 🏗️ Kiến trúc hệ thống
- Coming Soon


## 📁 Cấu trúc thư mục

```
go-food/
├── api/                          # API Gateway & Shared API specs
│   ├── gateway/                  # Kong/Nginx configurations
│   │   ├── plugins/
│   │   ├── routes/
│   │   └── middleware/
│   ├── proto/                    # Protocol buffer definitions
│   │   ├── user/
│   │   ├── restaurant/
│   │   ├── menu/
│   │   ├── order/
│   │   ├── payment/
│   │   ├── delivery/
│   │   └── common/
│   └── openapi/                  # OpenAPI specifications
│       ├── user-service.yaml
│       ├── restaurant-service.yaml
│       ├── menu-service.yaml
│       ├── order-service.yaml
│       ├── payment-service.yaml
│       ├── delivery-service.yaml
│       └── gateway.yaml
├── internal/
│   └── services/                 # Microservices
│       ├── user-service/
│       │   ├── cmd/
│       │   ├── internal/
│       │   │   ├── handler/
│       │   │   ├── service/
│       │   │   ├── repository/
│       │   │   └── domain/
│       │   ├── migrations/
│       │   └── configs/
│       ├── restaurant-service/   # Vendor/Seller service
│       │   ├── cmd/
│       │   ├── internal/
│       │   ├── migrations/
│       │   └── configs/
│       ├── menu-service/         # Product catalog service
│       │   ├── cmd/
│       │   ├── internal/
│       │   ├── migrations/
│       │   └── configs/
│       ├── cart-service/
│       │   ├── cmd/
│       │   ├── internal/
│       │   └── configs/
│       ├── order-service/
│       │   ├── cmd/
│       │   ├── internal/
│       │   ├── migrations/
│       │   └── configs/
│       ├── payment-service/
│       │   ├── cmd/
│       │   ├── internal/
│       │   │   ├── gateway/      # Payment gateway integrations
│       │   │   ├── escrow/
│       │   │   └── webhook/
│       │   ├── migrations/
│       │   └── configs/
│       ├── delivery-service/     # Shipping service
│       │   ├── cmd/
│       │   ├── internal/
│       │   │   ├── tracking/
│       │   │   ├── routing/
│       │   │   └── shipper/
│       │   ├── migrations/
│       │   └── configs/
│       ├── notification-service/
│       │   ├── cmd/
│       │   ├── internal/
│       │   │   ├── email/
│       │   │   ├── sms/
│       │   │   ├── push/
│       │   │   └── websocket/
│       │   └── configs/
│       ├── review-service/
│       │   ├── cmd/
│       │   ├── internal/
│       │   ├── migrations/
│       │   └── configs/
│       ├── search-service/       # Search & recommendation
│       │   ├── cmd/
│       │   ├── internal/
│       │   │   ├── elasticsearch/
│       │   │   ├── recommendation/
│       │   │   └── indexing/
│       │   └── configs/
│       ├── media-service/        # File upload & processing
│       │   ├── cmd/
│       │   ├── internal/
│       │   │   ├── upload/
│       │   │   ├── processing/
│       │   │   └── cdn/
│       │   └── configs/
│       └── analytics-service/    # Business intelligence
│           ├── cmd/
│           ├── internal/
│           │   ├── aggregation/
│           │   ├── reporting/
│           │   └── dashboard/
│           └── configs/
├── pkg/                          # Shared libraries
│   ├── auth/                     # JWT, OAuth2 utilities
│   ├── config/                   # Configuration management
│   ├── database/                 # Database utilities
│   │   ├── postgres/
│   │   ├── mongodb/
│   │   ├── redis/
│   │   └── migrations/
│   ├── middleware/               # HTTP middleware
│   │   ├── cors/
│   │   ├── rate-limit/
│   │   ├── validation/
│   │   └── logging/
│   ├── models/                   # Shared domain models
│   │   ├── user/
│   │   ├── restaurant/
│   │   ├── menu/
│   │   ├── order/
│   │   └── common/
│   ├── messaging/                # Message broker utilities
│   │   ├── rabbitmq/
│   │   ├── kafka/
│   │   └── events/
│   ├── external/                 # External service clients
│   │   ├── payment/
│   │   ├── maps/
│   │   ├── email/
│   │   └── sms/
│   └── utils/                    # Common utilities
│       ├── crypto/
│       ├── validator/
│       ├── logger/
│       └── http/
├── deployments/                  # Deployment configurations
│   ├── docker/
│   │   ├── Dockerfile.user-service
│   │   ├── Dockerfile.restaurant-service
│   │   ├── Dockerfile.menu-service
│   │   ├── Dockerfile.order-service
│   │   ├── Dockerfile.payment-service
│   │   ├── Dockerfile.delivery-service
│   │   ├── docker-compose.yml
│   │   └── docker-compose.prod.yml
│   ├── kubernetes/
│   │   ├── namespace.yaml
│   │   ├── services/
│   │   ├── deployments/
│   │   ├── configmaps/
│   │   ├── secrets/
│   │   └── ingress/
│   └── terraform/
│       ├── aws/
│       ├── gcp/
│       └── azure/
├── scripts/                      # Build and deployment scripts
│   ├── build.sh
│   ├── deploy.sh
│   ├── test.sh
│   └── migrate.sh
├── docs/                         # Documentation
│   ├── api/                      # API documentation
│   ├── architecture/             # Architecture diagrams
│   ├── deployment/               # Deployment guides
│   └── development/              # Development guides
├── tests/                        # Integration tests
│   ├── e2e/                      # End-to-end tests
│   ├── integration/              # Service integration tests
│   └── load/                     # Performance tests
├── tools/                        # Development tools
│   ├── proto-gen/                # Protocol buffer generation
│   ├── mock-gen/                 # Mock generation
│   └── migrate/                  # Database migration tools
├── monitoring/                   # Monitoring configurations
│   ├── prometheus/
│   ├── grafana/
│   ├── jaeger/
│   └── elk/
├── Makefile                      # Build automation
├── go.mod
├── go.sum
└── README.md
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

### Core Tables (PostgreSQL)

#### User Management
- **users**: User profiles (customers, restaurant owners, delivery drivers)
- **user_addresses**: Multiple delivery addresses per user
- **user_sessions**: Active user sessions and JWT tokens
- **user_preferences**: Food preferences, dietary restrictions

#### Restaurant Management
- **restaurants**: Restaurant/vendor information and business details
- **restaurant_owners**: Ownership relationships
- **restaurant_hours**: Operating hours and availability
- **restaurant_areas**: Delivery coverage areas

#### Menu & Products
- **categories**: Food categories (appetizers, mains, desserts, etc.)
- **menu_items**: Food items with pricing and descriptions
- **menu_item_variants**: Size, spice level, customizations
- **menu_item_options**: Add-ons and modifiers
- **menu_availability**: Time-based availability of items

#### Orders & Shopping
- **shopping_carts**: User shopping cart state
- **cart_items**: Items in cart with customizations
- **orders**: Order header information
- **order_items**: Detailed order line items
- **order_status_history**: Order state tracking

#### Payment & Financial
- **payments**: Payment transaction records
- **payment_methods**: Stored payment methods
- **refunds**: Refund processing records
- **restaurant_payouts**: Earnings distribution to restaurants
- **transaction_fees**: Platform commission tracking

#### Delivery & Logistics
- **delivery_drivers**: Driver profiles and vehicle info
- **delivery_assignments**: Order-driver assignments
- **delivery_tracking**: Real-time location updates
- **delivery_routes**: Optimized delivery routes
- **delivery_fees**: Dynamic pricing for delivery

#### Reviews & Ratings
- **reviews**: Restaurant and food reviews
- **review_media**: Review photos and videos
- **review_votes**: Helpful/unhelpful votes
- **driver_reviews**: Delivery service ratings

#### Notifications & Communication
- **notifications**: User notification history
- **notification_preferences**: User notification settings
- **push_tokens**: Device tokens for push notifications

### Document Storage (MongoDB)

#### Analytics & Logs
- **user_behavior_logs**: Click streams, search patterns
- **order_analytics**: Aggregated order metrics
- **restaurant_analytics**: Business intelligence data
- **system_logs**: Application and error logs
- **audit_trails**: Security and compliance logs

#### Media & Content
- **menu_images**: Food photos and restaurant images
- **review_media**: User-uploaded review content
- **promotional_content**: Marketing materials and banners

### Cache Layer (Redis)

#### Session Management
- **user_sessions**: Active user sessions
- **cart_cache**: Real-time shopping cart state
- **driver_locations**: Live driver position tracking

#### Performance Optimization
- **menu_cache**: Frequently accessed menu data
- **restaurant_cache**: Popular restaurant information
- **search_cache**: Search results and suggestions
- **recommendation_cache**: Personalized recommendations

#### Rate Limiting & Security
- **api_rate_limits**: Request throttling per user/IP
- **login_attempts**: Failed login tracking
- **otp_codes**: Temporary verification codes

### Search Index (Elasticsearch)

#### Search Optimization
- **restaurant_index**: Restaurant search with location-based queries
- **menu_item_index**: Food search with filters (cuisine, price, dietary)
- **review_index**: Review content for sentiment analysis
- **user_preference_index**: Personalization data

### Database Relationships

#### Key Foreign Keys
- users ↔ user_addresses (1:N)
- restaurants ↔ menu_items (1:N)
- menu_items ↔ menu_item_variants (1:N)
- users ↔ orders (1:N)
- orders ↔ order_items (1:N)
- orders ↔ payments (1:1)
- orders ↔ delivery_assignments (1:1)
- restaurants ↔ reviews (1:N)
- delivery_drivers ↔ delivery_assignments (1:N)

#### Indexing Strategy
- **Geospatial**: Restaurant locations, delivery areas
- **Composite**: User + restaurant for order history
- **Text**: Menu item names and descriptions
- **Time-based**: Order timestamps, delivery windows

## 🎯 Microservices

### Core Services

#### 1. User Management Service
- **Chức năng**: Authentication, authorization, profile management
- **Database**: PostgreSQL (user profiles, addresses)
- **Cache**: Redis (sessions, tokens)
- **Communication**: gRPC + REST API

#### 2. Restaurant Service (Vendor/Seller)
- **Chức năng**: Restaurant/shop management, analytics, order fulfillment
- **Database**: PostgreSQL (restaurant info, business data)
- **Communication**: gRPC for internal, REST for dashboard

#### 3. Menu Service (Product Catalog)
- **Chức năng**: Food items, categories, pricing, inventory
- **Database**: PostgreSQL (menu items, categories)
- **Search**: Elasticsearch indexing
- **Media**: MongoDB (images, descriptions)

#### 4. Shopping Cart Service
- **Chức năng**: Cart management, temporary storage
- **Cache**: Redis (cart state, session-based)
- **Database**: PostgreSQL (persistent carts)

#### 5. Order Service
- **Chức năng**: Order lifecycle, status tracking, history
- **Database**: PostgreSQL (orders, order_items)
- **Events**: Order state changes via message broker

#### 6. Payment Service
- **Chức năng**: Payment processing, refunds, escrow
- **Database**: PostgreSQL (payment records, transactions)
- **External**: Payment gateways (Stripe, VNPay, Momo)
- **Security**: PCI compliance, encryption

#### 7. Delivery Service (Shipping)
- **Chức năng**: Delivery assignment, tracking, route optimization
- **Database**: PostgreSQL (delivery info, shipper data)
- **Real-time**: WebSocket for live tracking
- **External**: Maps API for geocoding

#### 8. Notification Service
- **Chức năng**: Real-time notifications, push notifications
- **Message Queue**: RabbitMQ/Kafka for async messaging
- **Channels**: Email, SMS, push, in-app notifications

#### 9. Review Service
- **Chức năng**: Ratings, reviews, feedback management
- **Database**: PostgreSQL (reviews, ratings)
- **Media**: MongoDB (review images/videos)

#### 10. Search & Recommendation Service
- **Chức năng**: Search, filters, personalized recommendations
- **Search Engine**: Elasticsearch (full-text search, filters)
- **ML**: Recommendation algorithms, user behavior tracking
- **Cache**: Redis (search results, suggestions)

#### 11. Media Service
- **Chức năng**: File upload, image processing, CDN
- **Storage**: AWS S3/MinIO for file storage
- **Processing**: Image resizing, compression
- **CDN**: CloudFront for global delivery

#### 12. Analytics Service
- **Chức năng**: Business intelligence, reporting, metrics
- **Database**: MongoDB (analytics data, logs)
- **Processing**: Real-time data aggregation
- **Visualization**: Dashboard APIs for reporting

### Service Communication Patterns

#### Synchronous Communication
- **gRPC**: Internal service-to-service calls
- **REST API**: External client communications
- **GraphQL**: Unified API layer (optional)

#### Asynchronous Communication
- **Event Sourcing**: Domain events for state changes
- **Message Queues**: Background job processing
- **Pub/Sub**: Real-time notifications and updates

#### Data Management
- **Database per Service**: Each service owns its data
- **Event-driven Architecture**: Services communicate via events
- **CQRS**: Separate read/write models for complex queries
- **Saga Pattern**: Distributed transaction management

#### Service Discovery & Load Balancing
- **Service Registry**: Consul/Etcd for service registration
- **Load Balancer**: Nginx/HAProxy for traffic distribution
- **Health Checks**: Automatic service health monitoring
- **Circuit Breaker**: Fault tolerance and resilience

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