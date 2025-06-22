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

### Tổng quan kiến trúc

```
                        ┌─────────────────────────────────────┐
                        │           Frontend Layer            │
                        │                                     │
                ┌───────┼──────┐ ┌────────┼────────┐  ┌───────┼─────┐
                │  Web Client  │ │  Mobile Apps    │  │ Admin Panel │
                │   (React)    │ │ (React Native)  │  │  (Vue.js)   │
                └──────────────┘ └─────────────────┘  └─────────────┘
                        │                 │                   │
                        └─────────────────┼───────────────────┘
                                          │
                                ┌─────────▼─────────┐
                                │   Load Balancer   │
                                │  (Nginx/HAProxy)  │
                                │  - SSL/TLS        │
                                │  - Static Content │
                                └─────────┬─────────┘
                                          │
                                ┌─────────▼─────────┐
                                │   API Gateway     │
                                │   (Kong/Nginx)    │
                                │   - Rate Limiting │
                                │   - Authentication│
                                │   - Request Routing│
                                │   - Circuit Breaker│
                                └─────────┬─────────┘
                                          │
                        ┌─────────────────┼─────────────────┐
                        │                 │                 │
              ┌─────────▼─────────┐ ┌─────▼─────┐ ┌─────────▼─────────┐
              │  Authentication   │ │ Business  │ │   Data Layer      │
              │    Services       │ │ Services  │ │    Services       │
              │                   │ │           │ │                   │
              │ - User Management │ │ - Products│ │ - PostgreSQL      │
              │ - JWT Validation  │ │ - Orders  │ │ - Redis Cache     │
              │ - Role Management │ │ - Payment │ │ - Elasticsearch   │
              └───────────────────┘ │ - Shipping│ │ - MongoDB         │
                                    │ - Reviews │ │ - RabbitMQ        │
                                    └───────────┘ └───────────────────┘
```

### Microservices Architecture Detail

#### Detailed Service Interaction Flow

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            Complete Order Flow                                  │
└─────────────────────────────────────────────────────────────────────────────────┘

1. User Authentication Flow:
   ┌─────────┐    HTTP POST     ┌─────────────┐    SQL Query    ┌─────────────┐
   │ Client  │ ───────────────► │User Service │ ──────────────► │ PostgreSQL  │
   │         │    /auth/login   │Port: 8001   │                 │ Users DB    │
   └─────────┘                  └─────┬───────┘                 └─────────────┘
        ▲                             │
        │         JWT Token           ▼
        └─────────────────────── Redis Cache
                                  Session Store

2. Product Search Flow:
   ┌─────────┐   HTTP GET      ┌─────────────┐   Elasticsearch  ┌─────────────┐
   │ Client  │ ──────────────► │Search Svc   │ ──────────────►  │Elasticsearch│
   │         │ /products/search│Port: 8007   │                  │Index        │
   └─────────┘                 └─────┬───────┘                  └─────────────┘
        ▲                            │
        │         Results            ▼ gRPC Call
        └────────────────────┌─────────────┐
                             │Product Svc  │
                             │Port: 8003   │
                             └─────────────┘

3. Add to Cart Flow:
   ┌─────────┐   HTTP POST     ┌─────────────┐    gRPC         ┌─────────────┐
   │ Client  │ ──────────────► │Cart Service │ ──────────────► │Product Svc  │
   │         │ /cart/items     │Port: 8004   │ Check Inventory │Port: 8003   │
   └─────────┘                 └─────┬───────┘                 └─────────────┘
                                     │
                                     ▼ Redis SET
                               ┌─────────────┐
                               │Redis Cache  │
                               │Cart Data    │
                               └─────────────┘

4. Order Creation Flow:
   ┌─────────┐   HTTP POST     ┌─────────────┐    gRPC         ┌─────────────┐
   │ Client  │ ──────────────► │Order Service│ ──────────────► │Cart Service │
   │         │   /orders       │Port: 8005   │  Get Cart       │Port: 8004   │
   └─────────┘                 └─────┬───────┘                 └─────────────┘
                                     │
                                     ▼ SQL INSERT
                               ┌─────────────┐
                               │PostgreSQL   │
                               │Orders DB    │
                               └─────┬───────┘
                                     │
                                     ▼ Publish Event
                               ┌─────────────┐
                               │RabbitMQ     │
                               │order.created│
                               └─────┬───────┘
                                     │
                    ┌────────────────┼────────────────┐
                    ▼                ▼                ▼
            ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
            │Payment Svc  │  │Notification │  │Analytics    │
            │Subscribe    │  │Service      │  │Service      │
            │Port: 8006   │  │Port: 8008   │  │Port: 8010   │
            └─────────────┘  └─────────────┘  └─────────────┘

5. Payment Processing Flow:
   ┌─────────────┐   Process     ┌─────────────┐   HTTP POST   ┌─────────────┐
   │Payment Svc  │ ────────────► │External     │ ────────────► │Stripe API   │
   │Port: 8006   │   Payment     │Gateway      │   Payment     │             │
   └─────┬───────┘               └─────────────┘               └─────┬───────┘
         │                                                           │
         ▼ SQL INSERT                                                 ▼ Webhook
   ┌─────────────┐                                             ┌─────────────┐
   │PostgreSQL   │                                             │Payment      │
   │Payments DB  │                                             │Confirmation │
   └─────┬───────┘                                             └─────┬───────┘
         │                                                           │
         ▼ Publish Event                                             ▼
   ┌─────────────┐                                             ┌─────────────┐
   │RabbitMQ     │ ◄───────────────────────────────────────────│Update Order │
   │payment.done │                                             │Status       │
   └─────────────┘                                             └─────────────┘

6. Shipping Assignment Flow:
   ┌─────────────┐   Subscribe   ┌─────────────┐   gRPC Call   ┌─────────────┐
   │Shipping Svc │ ────────────► │RabbitMQ     │               │Order Service│
   │Port: 8009   │  payment.done │             │ ◄─────────────│Get Order    │
   └─────┬───────┘               └─────────────┘               │Details      │
         │                                                     └─────────────┘
         ▼ Find Available Shipper
   ┌─────────────┐
   │PostgreSQL   │
   │Shippers DB  │
   └─────┬───────┘
         │
         ▼ Assign & Notify
   ┌─────────────┐   WebSocket   ┌─────────────┐
   │Notification │ ────────────► │Shipper App  │
   │Service      │   Real-time   │             │
   └─────────────┘               └─────────────┘
```

#### Service Dependencies Matrix

```
┌─────────────────┬──────────────────────────────────────────────────────────────┐
│    Service      │                    Dependencies                              │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│ User Service    │ • PostgreSQL (users_db)                                      │
│                 │ • Redis (sessions, OTP)                                      │
│                 │ • SMTP (email service)                                       │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│ Shop Service    │ • PostgreSQL (shops_db)                                      │
│                 │ • User Service (gRPC - owner validation)                     │
│                 │ • Analytics Service (gRPC - metrics)                         │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│ Product Service │ • PostgreSQL (products_db)                                   │
│                 │ • Elasticsearch (product indexing)                           │
│                 │ • Shop Service (gRPC - shop validation)                      │
│                 │ • Media Service (gRPC - image processing)                    │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│ Cart Service    │ • Redis (cart data)                                          │
│                 │ • Product Service (gRPC - inventory check)                   │
│                 │ • User Service (gRPC - user validation)                      │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│ Order Service   │ • PostgreSQL (orders_db)                                     │
│                 │ • Cart Service (gRPC - get cart)                             │
│                 │ • Product Service (gRPC - reserve inventory)                 │
│                 │ • RabbitMQ (publish order events)                            │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│ Payment Service │ • PostgreSQL (payments_db)                                   │
│                 │ • External APIs (Stripe, VNPay, Momo)                        │
│                 │ • Order Service (gRPC - order details)                       │
│                 │ • RabbitMQ (publish payment events)                          │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│ Shipping Service│ • PostgreSQL (shipping_db)                                   │
│                 │ • Maps API (geocoding, routing)                              │
│                 │ • Order Service (gRPC - delivery address)                    │
│                 │ • RabbitMQ (subscribe order events)                          │
│                 │ • WebSocket (real-time tracking)                             │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│Search & Rec Svc │ • Elasticsearch (search index)                               │
│                 │ • Redis (cache results)                                      │
│                 │ • Product Service (gRPC - product data)                      │
│                 │ • Analytics Service (gRPC - user behavior)                   │
│                 │ • ML Models (recommendation engine)                          │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│ Review Service  │ • PostgreSQL (reviews_db)                                    │
│                 │ • Order Service (gRPC - verify purchase)                     │
│                 │ • Media Service (gRPC - upload images)                       │
│                 │ • Notification Service (gRPC - notify shop)                  │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│Notification Svc │ • RabbitMQ (consume all events)                              │
│                 │ • Redis (notification preferences)                           │
│                 │ • SMTP (email), SMS Gateway, FCM (push)                      │
│                 │ • WebSocket (real-time notifications)                        │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│ Media Service   │ • S3/MinIO (file storage)                                    │
│                 │ • Redis (upload sessions)                                    │
│                 │ • CDN (content delivery)                                     │
│                 │ • Image Processing (resize, compress)                        │
├─────────────────┼──────────────────────────────────────────────────────────────┤
│Analytics Service│ • MongoDB (analytics data)                                   │
│                 │ • InfluxDB (time series metrics)                             │
│                 │ • Kafka (event streaming)                                    │
│                 │ • All Services (gRPC - collect metrics)                      │
└─────────────────┴──────────────────────────────────────────────────────────────┘
```

#### Communication Protocols & Ports

```
┌─────────────────┬──────────┬──────────┬─────────────────────────────────────────┐
│    Service      │HTTP Port │gRPC Port │            Protocols Used               │
├─────────────────┼──────────┼──────────┼─────────────────────────────────────────┤
│ API Gateway     │   8000   │    -     │ HTTP/HTTPS, WebSocket                   │
│ User Service    │   8001   │   9001   │ HTTP, gRPC, SMTP                        │
│ Shop Service    │   8002   │   9002   │ HTTP, gRPC                              │
│ Product Service │   8003   │   9003   │ HTTP, gRPC                              │
│ Cart Service    │   8004   │   9004   │ HTTP, gRPC, Redis Protocol              │
│ Order Service   │   8005   │   9005   │ HTTP, gRPC, AMQP                        │
│ Payment Service │   8006   │   9006   │ HTTP, gRPC, AMQP, HTTPS (Gateways)      │
│ Search Service  │   8007   │   9007   │ HTTP, gRPC, Elasticsearch API           │
│ Notification    │   8008   │   9008   │ gRPC, AMQP, WebSocket, SMTP, SMS        │
│ Shipping Service│   8009   │   9009   │ HTTP, gRPC, AMQP, WebSocket, Maps API   │
│ Media Service   │   8010   │   9010   │ HTTP, gRPC, S3 API                      │
│ Analytics       │   8011   │   9011   │ gRPC, Kafka, MongoDB, InfluxDB          │
│ Review Service  │   8012   │   9012   │ HTTP, gRPC                              │
└─────────────────┴──────────┴──────────┴─────────────────────────────────────────┘

Database Connections:
├── PostgreSQL: 5432 (User, Shop, Product, Order, Payment, Shipping, Review)
├── Redis: 6379 (Cache, Sessions, Cart, Pub/Sub)
├── MongoDB: 27017 (Analytics, Media metadata, Logs)
├── Elasticsearch: 9200 (Search indices)
├── RabbitMQ: 5672 (Message queues)
└── InfluxDB: 8086 (Time series metrics)
```

### Service Communication Flow

#### 1. User Registration & Authentication Flow
```
Mobile/Web Client
       │
       ▼
   API Gateway ──────────┐
       │                │
       ▼                ▼
User Management    Rate Limiting
   Service             │
       │               ▼
       ▼         Request Validation
  PostgreSQL           │
   (Users DB)          ▼
       │         JWT Generation
       ▼               │
   Redis Cache         ▼
  (Sessions)      Response to Client
```

#### 2. Product Search & Discovery Flow
```
Client Request
       │
       ▼
   API Gateway
       │
       ▼
Search & Recommendation Service
       │
       ├─────────┬─────────────┐
       ▼         ▼             ▼
 Elasticsearch Product      User Behavior
  (Indexing)   Catalog      Analytics
       │       Service           │
       ▼         │               ▼
   Results       ▼           MongoDB
     │     PostgreSQL        (Analytics)
     │    (Products DB)          │
     ▼         │                 ▼
   Redis       ▼             Recommendation
  (Cache)   Product Info      Engine (ML)
     │         │                 │
     └─────────┼─────────────────┘
               ▼
        Aggregated Response
```

#### 3. Order Processing Flow
```
Shopping Cart ────────┐
       │              │
       ▼              ▼
   Cart Service   Order Service
       │              │
       ▼              ▼
   Redis Cache    PostgreSQL
  (Cart Data)    (Orders DB)
       │              │
       ▼              ▼
Product Catalog   Payment Service
   Service            │
       │              ▼
       ▼         Payment Gateway
Inventory Check   (Stripe/VNPay)
       │              │
       ▼              ▼
Stock Update     Payment Result
       │              │
       ▼              ▼
   Message Queue ──▶ Order Status
  (RabbitMQ)         Update
       │              │
       ▼              ▼
Notification     Shipping Service
  Service            │
       │              ▼
       ▼         Shipper Assignment
Email/SMS/Push       │
Notifications        ▼
                Real-time Tracking
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

## 🎯 Microservices

### Core Services

#### 1. User Management Service
- **API Endpoints**: Authentication (`/api/v1/auth/*`), User Profile (`/api/v1/users/*`)
- **Chức năng**:
  - Authentication & Authorization (JWT, OAuth2)
  - User registration, login, logout, password management
  - Profile management và address management
  - Role management (User, Shipper registration)
  - OTP verification và forgot password flow
- **Database**: PostgreSQL (user profiles, addresses, roles)
- **Cache**: Redis (JWT tokens, sessions, OTP)
- **Security**: JWT tokens, bcrypt hashing, rate limiting
- **Communication**: gRPC + REST API

#### 2. Shop Management Service (Vendor/Seller)
- **API Endpoints**: Shop Management (`/api/v1/shops/*`)
- **Chức năng**:
  - CRUD shop management và profile
  - Shop product management và inventory
  - Order fulfillment và status management  
  - Revenue analytics và performance reports
  - Promotion campaigns và discount management
  - Location-based shop search
- **Database**: PostgreSQL (shop info, business data, promotions)
- **Analytics**: Revenue tracking, order analytics, product performance
- **Communication**: gRPC for internal, REST for dashboard

#### 3. Product Catalog Service
- **API Endpoints**: Product Management (`/api/v1/products/*`)
- **Chức năng**:
  - CRUD product management với media support
  - Category và subcategory management
  - Price management và price history
  - Stock & inventory management với low-stock alerts
  - Product variants và related products
  - Advanced search với filters (price, rating, brand, location)
  - Brand và model management
- **Database**: PostgreSQL (products, categories, pricing, inventory)
- **Search**: Elasticsearch indexing cho full-text search
- **Media**: MongoDB (product images, descriptions)
- **Cache**: Redis (popular products, search results)

#### 4. Shopping Cart Service
- **API Endpoints**: Cart Management (`/api/v1/cart/*`)
- **Chức năng**:
  - Real-time cart management (add, update, remove items)
  - Cart persistence và saved carts
  - Total calculation với shipping fees và taxes
  - Coupon application và discount calculation
  - Cart restoration và multiple saved carts
- **Cache**: Redis (active cart state, session-based)
- **Database**: PostgreSQL (persistent carts, saved carts)
- **Real-time**: WebSocket cho cart updates

#### 5. Order Management Service
- **API Endpoints**: Order Processing (`/api/v1/orders/*`)
- **Chức năng**:
  - Order creation từ shopping cart
  - Order lifecycle management (pending → confirmed → shipped → delivered)
  - Order status tracking và timeline
  - Return và refund request processing
  - Purchase history và repeat orders
  - Order calculation với fees breakdown
- **Database**: PostgreSQL (orders, order_items, status_history)
- **Events**: Order state changes via message broker
- **Integration**: Payment service, shipping service

#### 6. Payment Service
- **API Endpoints**: Payment Processing (`/api/v1/payments/*`)
- **Chức năng**:
  - Multi-gateway payment processing (Stripe, PayPal, VNPay, Momo)
  - Payment method management
  - Refund và chargeback handling
  - Escrow service cho buyer protection
  - Transaction history và receipt generation
  - Webhook handling cho payment gateways
- **Database**: PostgreSQL (payment records, transactions, refunds)
- **External**: Payment gateways integration
- **Security**: PCI compliance, payment encryption

#### 7. Shipping & Delivery Service
- **API Endpoints**: Shipping Management (`/api/v1/shipping/*`)
- **Chức năng**:
  - Shipper registration và management
  - Shipping cost calculation based on distance/weight
  - Order assignment to shippers
  - Real-time tracking và location updates
  - Address validation và geocoding
  - Live tracking với WebSocket
- **Database**: PostgreSQL (delivery info, shipper data, tracking)
- **Real-time**: WebSocket cho live tracking
- **External**: Maps API cho geocoding và route optimization
- **Integration**: Order service cho delivery updates

#### 8. Search & Recommendation Service
- **API Endpoints**: Search (`/api/v1/search/*`), Recommendations (`/api/v1/recommendations/*`)
- **Chức năng**:
  - Advanced search với filters và autocomplete
  - Personalized recommendations based on behavior
  - Trending products và popular searches
  - User behavior tracking (view, click, search)
  - Price comparison và similar products
  - Price alerts và notifications
- **Search Engine**: Elasticsearch (full-text search, filters)
- **ML**: Recommendation algorithms, collaborative filtering
- **Cache**: Redis (search results, suggestions, trending data)
- **Analytics**: User behavior tracking và recommendation metrics

#### 9. Review & Rating Service
- **API Endpoints**: Reviews (`/api/v1/products/{id}/reviews`, `/api/v1/shops/{id}/reviews`)
- **Chức năng**:
  - Product và shop reviews với rating
  - Delivery service reviews
  - Media upload cho reviews (images, videos)
  - Review moderation và spam detection
  - Verified purchase reviews
  - Sentiment analysis và rating distribution
- **Database**: PostgreSQL (reviews, ratings, moderation)
- **Media**: MongoDB (review images/videos)
- **ML**: Sentiment analysis, spam detection

#### 10. Notification Service
- **API Endpoints**: Notifications (`/api/v1/notifications/*`)
- **Chức năng**:
  - Real-time notifications (order updates, delivery status)
  - Multi-channel notifications (email, SMS, push, in-app)
  - Notification preferences management
  - Template management cho automated notifications
  - WebSocket cho live notifications
- **Message Queue**: RabbitMQ/Kafka cho async messaging
- **Channels**: Email, SMS, push notifications, in-app
- **Real-time**: WebSocket connections cho live updates

#### Supporting Services

#### 11. Media Service
- **Chức năng**: 
  - File upload và image processing
  - Image resizing, compression, watermarking
  - CDN integration cho fast delivery
  - Video processing cho review media
- **Storage**: AWS S3/MinIO cho file storage
- **Processing**: Image/video processing pipeline
- **CDN**: CloudFront cho global content delivery
- **Integration**: Product service, review service

#### 12. Analytics Service
- **Chức năng**:
  - Business intelligence và reporting
  - Real-time metrics aggregation
  - Shop performance analytics
  - User behavior analytics
  - Revenue tracking và forecasting
- **Database**: MongoDB (analytics data, logs)
- **Processing**: Real-time data aggregation với Apache Kafka
- **Visualization**: Dashboard APIs cho business reporting
- **Integration**: All services for data collection

### Service Communication Patterns

#### Synchronous Communication
- **gRPC**: Internal service-to-service calls
  - User authentication validation
  - Product inventory checks
  - Payment processing
- **REST API**: External client communications
  - Mobile app integration
  - Web dashboard
  - Third-party integrations
- **GraphQL**: Unified API layer (optional)
  - Frontend data aggregation
  - Flexible query capabilities

#### Asynchronous Communication
- **Event Sourcing**: Domain events cho state changes
  - Order status updates
  - Payment confirmations
  - Inventory changes
- **Message Queues**: Background job processing
  - Email notifications
  - Image processing
  - Analytics data processing
- **Pub/Sub**: Real-time notifications và updates
  - Live order tracking
  - Real-time notifications
  - Price updates

#### Inter-Service Communication Flow
1. **Order Flow**:
   ```
   Cart Service → Order Service → Payment Service → Shipping Service
                                ↓
   Notification Service ← Analytics Service
   ```

2. **Search Flow**:
   ```
   Client → Search Service → Product Service → Shop Service
          ↓
   Recommendation Service → User Behavior Tracking
   ```

3. **Review Flow**:
   ```
   Client → Review Service → Media Service (for uploads)
          ↓
   Order Service (verify purchase) → Notification Service
   ```

#### Data Management
- **Database per Service**: Mỗi service sở hữu data riêng
- **Event-driven Architecture**: Services giao tiếp qua events
- **CQRS**: Tách read/write models cho complex queries
- **Saga Pattern**: Distributed transaction management
- **Data Consistency**: Eventually consistent với compensation patterns

#### Service Discovery & Load Balancing
- **Service Registry**: Consul/Etcd cho service registration
- **Load Balancer**: Nginx/HAProxy cho traffic distribution
- **Health Checks**: Automatic service health monitoring
- **Circuit Breaker**: Fault tolerance và resilience patterns
- **Rate Limiting**: API throttling và abuse prevention

#### Security & Cross-cutting Concerns
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