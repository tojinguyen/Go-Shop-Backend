# Go-Shop - Microservice E-commerce Delivery Platform

Nền tảng thương mại điện tử và giao hàng được xây dựng theo kiến trúc microservice sử dụng Go (Golang), tương tự như Shopee.

## 📋 Mục lục

- [Yêu cầu chức năng](#yêu-cầu-chức-năng)
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


## 🔗 Thiết kế API

### Authentication APIs
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/refresh
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
POST   /api/v1/auth/change-password
POST   /api/v1/auth/verify-otp
```

### User Management APIs
```
GET    /api/v1/users/profile
PUT    /api/v1/users/profile
GET    /api/v1/users/{id}
DELETE /api/v1/users/{id}

# Address Management
GET    /api/v1/users/addresses
POST   /api/v1/users/addresses
GET    /api/v1/users/addresses/{id}
PUT    /api/v1/users/addresses/{id}
DELETE /api/v1/users/addresses/{id}
PUT    /api/v1/users/addresses/{id}/default

# Role Management (User, Shipper)
POST   /api/v1/users/shipper/register
```

### Shop Management APIs
```
# Shop CRUD
GET    /api/v1/shops
POST   /api/v1/shops
GET    /api/v1/shops/{id}
PUT    /api/v1/shops/{id}
DELETE /api/v1/shops/{id}
GET    /api/v1/shops/search?location={lat,lng}&category={category}&radius={radius}

# Shop Profile Management
PUT    /api/v1/shops/{id}/profile

# Shop Product Management
GET    /api/v1/shops/{id}/products
POST   /api/v1/shops/{id}/products
GET    /api/v1/shops/{id}/products/{product_id}
PUT    /api/v1/shops/{id}/products/{product_id}
DELETE /api/v1/shops/{id}/products/{product_id}
GET    /api/v1/shops/{id}/products/categories
GET    /api/v1/shops/{id}/products/inventory
PUT    /api/v1/shops/{id}/products/{product_id}/inventory
GET    /api/v1/shops/{id}/products/low-stock

# Shop Orders & Fulfillment
GET    /api/v1/shops/{id}/orders
PUT    /api/v1/shops/{id}/orders/{order_id}/status
POST   /api/v1/shops/{id}/orders/{order_id}/fulfill
GET    /api/v1/shops/{id}/orders/pending
GET    /api/v1/shops/{id}/orders/history

# Shop Analytics & Reports
GET    /api/v1/shops/{id}/analytics/revenue
GET    /api/v1/shops/{id}/analytics/orders
GET    /api/v1/shops/{id}/analytics/products
GET    /api/v1/shops/{id}/reports/sales
GET    /api/v1/shops/{id}/reports/performance

# Promotions & Campaigns
GET    /api/v1/shops/{id}/promotions
POST   /api/v1/shops/{id}/promotions
GET    /api/v1/shops/{id}/promotions/{promo_id}
PUT    /api/v1/shops/{id}/promotions/{promo_id}
DELETE /api/v1/shops/{id}/promotions/{promo_id}
POST   /api/v1/shops/{id}/campaigns
GET    /api/v1/shops/{id}/campaigns
PUT    /api/v1/shops/{id}/campaigns/{campaign_id}/status
```

### Product Catalog APIs
```
# Product CRUD
GET    /api/v1/products
POST   /api/v1/products
GET    /api/v1/products/{id}
PUT    /api/v1/products/{id}
DELETE /api/v1/products/{id}

# Product Catalog Management
PUT    /api/v1/products/{id}/catalog
POST   /api/v1/products/{id}/media
DELETE /api/v1/products/{id}/media/{media_id}
PUT    /api/v1/products/{id}/brand
PUT    /api/v1/products/{id}/model

# Price Management
PUT    /api/v1/products/{id}/price
GET    /api/v1/products/{id}/price-history
POST   /api/v1/products/{id}/discount
DELETE /api/v1/products/{id}/discount

# Stock & Inventory
GET    /api/v1/products/{id}/inventory
PUT    /api/v1/products/{id}/inventory
POST   /api/v1/products/{id}/inventory/adjustment
GET    /api/v1/products/low-stock

# Categories & Classification
GET    /api/v1/products/categories
POST   /api/v1/products/categories
GET    /api/v1/products/categories/{id}
PUT    /api/v1/products/categories/{id}
DELETE /api/v1/products/categories/{id}
GET    /api/v1/products/categories/{id}/subcategories
POST   /api/v1/products/categories/{id}/subcategories

# Search & Filtering
GET    /api/v1/products/search?q={query}&category={category}&price_min={min}&price_max={max}&rating={rating}&location={location}
GET    /api/v1/products/filter?brand={brand}&model={model}&attributes={attributes}

# Product Relations
GET    /api/v1/products/{id}/related
GET    /api/v1/products/{id}/reviews
GET    /api/v1/products/{id}/variants
POST   /api/v1/products/{id}/variants
```

### Shopping Cart APIs
```
GET    /api/v1/cart
POST   /api/v1/cart/items
PUT    /api/v1/cart/items/{id}
DELETE /api/v1/cart/items/{id}
DELETE /api/v1/cart/clear
GET    /api/v1/cart/summary
POST   /api/v1/cart/save
GET    /api/v1/cart/saved
POST   /api/v1/cart/restore/{saved_cart_id}
GET    /api/v1/cart/calculate-total
POST   /api/v1/cart/apply-coupon
DELETE /api/v1/cart/remove-coupon
```

### Order Management APIs
```
# Order Creation & Management
POST   /api/v1/orders
GET    /api/v1/orders
GET    /api/v1/orders/{id}
PUT    /api/v1/orders/{id}/status
DELETE /api/v1/orders/{id}

# Order Status Management
GET    /api/v1/orders/{id}/status-history
PUT    /api/v1/orders/{id}/confirm
PUT    /api/v1/orders/{id}/ship
PUT    /api/v1/orders/{id}/deliver
PUT    /api/v1/orders/{id}/cancel

# Order Calculations
GET    /api/v1/orders/{id}/calculation
POST   /api/v1/orders/calculate-preview
GET    /api/v1/orders/{id}/fees/breakdown

# Returns & Refunds
POST   /api/v1/orders/{id}/return
GET    /api/v1/orders/{id}/return-status
POST   /api/v1/orders/{id}/refund/request
GET    /api/v1/orders/returns
GET    /api/v1/orders/refunds

# Order Tracking
GET    /api/v1/orders/{id}/tracking
GET    /api/v1/orders/{id}/timeline

# Purchase History
GET    /api/v1/orders/history
GET    /api/v1/orders/history/summary
GET    /api/v1/orders/repeat/{id}
```

### Payment APIs
```
# Payment Processing
POST   /api/v1/payments
GET    /api/v1/payments/{id}
PUT    /api/v1/payments/{id}/status
POST   /api/v1/payments/{id}/capture

# Payment Methods
GET    /api/v1/payments/methods
POST   /api/v1/payments/methods
DELETE /api/v1/payments/methods/{id}
PUT    /api/v1/payments/methods/{id}/default

# Payment Gateway Integration
POST   /api/v1/payments/stripe/webhook
POST   /api/v1/payments/paypal/webhook
POST   /api/v1/payments/vnpay/webhook
POST   /api/v1/payments/momo/webhook

# Refunds & Chargebacks
POST   /api/v1/payments/{id}/refund
GET    /api/v1/payments/{id}/refund-status
GET    /api/v1/payments/chargebacks
POST   /api/v1/payments/{id}/dispute

# Escrow Service
POST   /api/v1/payments/escrow/hold
POST   /api/v1/payments/escrow/release
GET    /api/v1/payments/escrow/{id}/status

# Transaction History
GET    /api/v1/payments/history
GET    /api/v1/payments/transactions
GET    /api/v1/payments/receipts/{id}
```

### Shipping & Delivery APIs
```
# Shipping Methods & Calculation
GET    /api/v1/shipping/methods
POST   /api/v1/shipping/calculate
GET    /api/v1/shipping/providers
POST   /api/v1/shipping/labels

# Shipper Management
GET    /api/v1/shipping/shippers
POST   /api/v1/shipping/shippers/register
GET    /api/v1/shipping/shippers/{id}
PUT    /api/v1/shipping/shippers/{id}/status
GET    /api/v1/shipping/shippers/available

# Order Assignment & Tracking
POST   /api/v1/shipping/assign/{order_id}
GET    /api/v1/shipping/{order_id}/tracking
PUT    /api/v1/shipping/{order_id}/status
POST   /api/v1/shipping/{order_id}/location

# Real-time Tracking
GET    /api/v1/shipping/{order_id}/live-tracking
POST   /api/v1/shipping/{order_id}/update-location
GET    /api/v1/shipping/shipper/{shipper_id}/location

# Address & Geocoding
POST   /api/v1/shipping/validate-address
POST   /api/v1/shipping/geocode
GET    /api/v1/shipping/distance-matrix

# Shipping Costs & Fees
GET    /api/v1/shipping/cost-calculator
POST   /api/v1/shipping/calculate-fees
GET    /api/v1/shipping/weight-pricing
```

### Search & Recommendation APIs
```
# Search
GET    /api/v1/search?q={query}&filters={filters}&location={location}&sort={sort}
GET    /api/v1/search/suggestions?q={partial_query}
GET    /api/v1/search/autocomplete?q={query}
POST   /api/v1/search/advanced
GET    /api/v1/search/filters/available

# Personalized Recommendations
GET    /api/v1/recommendations/products
GET    /api/v1/recommendations/shops
GET    /api/v1/recommendations/based-on-behavior
GET    /api/v1/recommendations/similar-users

# Trending & Popular
GET    /api/v1/trending/products
GET    /api/v1/trending/shops
GET    /api/v1/trending/categories
GET    /api/v1/popular/searches

# User Behavior Tracking
GET    /api/v1/recent/products
GET    /api/v1/recent/searches
POST   /api/v1/behavior/view-product
POST   /api/v1/behavior/search
POST   /api/v1/behavior/click

# Price Comparison
GET    /api/v1/products/{id}/price-comparison
GET    /api/v1/products/similar-price?product_id={id}
GET    /api/v1/price-alerts
POST   /api/v1/price-alerts
DELETE /api/v1/price-alerts/{id}
```

### Review & Rating APIs
```
# Product Reviews
GET    /api/v1/products/{id}/reviews
POST   /api/v1/products/{id}/reviews
PUT    /api/v1/reviews/{id}
DELETE /api/v1/reviews/{id}
GET    /api/v1/reviews/{id}

# Shop Reviews
GET    /api/v1/shops/{id}/reviews
POST   /api/v1/shops/{id}/reviews
GET    /api/v1/shops/{id}/rating-summary

# Delivery Reviews
POST   /api/v1/delivery/{order_id}/review
GET    /api/v1/delivery/reviews
GET    /api/v1/shippers/{id}/reviews

# Media Upload for Reviews
POST   /api/v1/reviews/{id}/media
DELETE /api/v1/reviews/{id}/media/{media_id}
GET    /api/v1/reviews/{id}/media

# Review Management
GET    /api/v1/reviews/moderation/pending
PUT    /api/v1/reviews/{id}/approve
PUT    /api/v1/reviews/{id}/reject
POST   /api/v1/reviews/{id}/report

# Review Analytics
GET    /api/v1/reviews/verified-purchases
GET    /api/v1/reviews/sentiment-analysis
GET    /api/v1/reviews/rating-distribution
```

### Notification APIs
```
# Real-time Notifications
GET    /api/v1/notifications
PUT    /api/v1/notifications/{id}/read
DELETE /api/v1/notifications/{id}
POST   /api/v1/notifications/mark-all-read

# Push Notifications
POST   /api/v1/notifications/push/register
DELETE /api/v1/notifications/push/unregister
GET    /api/v1/notifications/preferences
PUT    /api/v1/notifications/preferences

# Email & SMS
POST   /api/v1/notifications/email/send
POST   /api/v1/notifications/sms/send
GET    /api/v1/notifications/templates
POST   /api/v1/notifications/templates

# WebSocket Events
WS     /api/v1/notifications/live
WS     /api/v1/orders/{id}/live-updates
WS     /api/v1/shipping/{order_id}/live-tracking
```

## 🏗️ Kiến trúc hệ thống
- Coming Soon

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