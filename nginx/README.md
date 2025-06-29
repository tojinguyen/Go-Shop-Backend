# Nginx API Gateway Setup Guide

## Tổng quan

Nginx được cấu hình để hoạt động như một API Gateway và Load Balancer cho hệ thống microservices Go-Shop. Nó cung cấp:

- **API Gateway**: Định tuyến request đến các microservice tương ứng
- **Load Balancing**: Phân tải request giữa nhiều instance của cùng một service
- **Rate Limiting**: Giới hạn số request để tránh DDoS
- **CORS Support**: Hỗ trợ Cross-Origin Resource Sharing
- **Health Checks**: Kiểm tra sức khỏe của các service
- **Logging**: Ghi log chi tiết cho monitoring

## Cấu trúc thư mục

```
nginx/
├── Dockerfile              # Docker image cho Nginx
├── nginx.conf              # Cấu hình chính của Nginx
├── conf.d/
│   ├── upstreams.conf      # Định nghĩa upstream cho load balancing
│   └── api-gateway.conf    # Cấu hình routing cho API Gateway
└── logs/                   # Thư mục chứa log files
```

## API Endpoints

Tất cả các API sẽ được truy cập thông qua Nginx tại `http://localhost:80`:

### User Service
- `GET/POST /api/v1/users` - Quản lý user
- `POST /api/v1/auth/login` - Đăng nhập
- `POST /api/v1/auth/register` - Đăng ký
- `GET/PUT /api/v1/profile` - Quản lý profile

### Order Service
- `GET/POST /api/v1/orders` - Quản lý đơn hàng

### Payment Service
- `POST /api/v1/payments` - Xử lý thanh toán

### Shop Service
- `GET/POST /api/v1/shops` - Quản lý cửa hàng
- `GET /api/v1/products` - Danh sách sản phẩm

### Review Service
- `GET/POST /api/v1/reviews` - Quản lý đánh giá

### Shipping Service
- `GET/POST /api/v1/shipping` - Quản lý vận chuyển

### Notification Service
- `GET /api/v1/notifications` - Quản lý thông báo

## Rate Limiting

- API thông thường: 10 requests/second
- Login endpoint: 5 requests/minute
- Có thể điều chỉnh trong file `nginx.conf`

## Load Balancing

- Sử dụng thuật toán `least_conn` (kết nối ít nhất)
- Hỗ trợ health check và failover
- Có thể thêm nhiều instance của cùng một service

## Monitoring

### Health Check
```bash
curl http://localhost:80/health
```

### Kiểm tra logs
- Access logs: `nginx/logs/access.log`
- Error logs: `nginx/logs/error.log`

## Cấu hình nâng cao

### Thêm service mới
1. Thêm upstream definition trong `conf.d/upstreams.conf`
2. Thêm location block trong `conf.d/api-gateway.conf`
3. Reload cấu hình: `nginx.ps1 -Command reload`

### Điều chỉnh rate limiting
Sửa file `nginx.conf`:
```nginx
limit_req_zone $binary_remote_addr zone=api:10m rate=20r/s;
```

### Thêm HTTPS
1. Thêm SSL certificate vào thư mục `nginx/ssl/`
2. Cập nhật cấu hình trong `api-gateway.conf`
3. Thay đổi port 443 trong docker-compose.yml

## Troubleshooting

### Nginx không khởi động được
1. Kiểm tra cú pháp cấu hình: `nginx.ps1 -Command test`
2. Xem error logs: `nginx.ps1 -Command logs`
3. Kiểm tra port conflicts

### Service không response
1. Kiểm tra upstream service đang chạy
2. Xem Nginx access logs
3. Test direct connection đến service

### Performance tuning
1. Điều chỉnh `worker_processes` trong `nginx.conf`
2. Tăng `worker_connections`
3. Điều chỉnh buffer sizes cho proxy

## Security Best Practices

1. **Rate Limiting**: Đã được cấu hình sẵn
2. **Security Headers**: X-Frame-Options, X-XSS-Protection, etc.
3. **CORS**: Được cấu hình cho development, cần điều chỉnh cho production
4. **Input Validation**: Client request size limited to 100MB
5. **Error Handling**: Custom error pages không leak thông tin

## Next Steps

1. Implement SSL/TLS certificates
2. Add authentication middleware
3. Implement request/response logging
4. Add metrics collection (Prometheus)
5. Configure log rotation
