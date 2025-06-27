# SMTP Email Service - Go-Shop

Email service cơ bản dùng chung cho tất cả các microservices trong hệ thống Go-Shop.

## Tính năng

- ✅ Hỗ trợ SMTP với TLS/SSL
- ✅ Gửi email text và HTML
- ✅ Template engine cho email đẹp
- ✅ Cấu hình sẵn cho Gmail, Outlook, Yahoo
- ✅ Có thể tùy chỉnh SMTP server
- ✅ Interface đơn giản, linh hoạt để từng service tự implement logic riêng

## Interface

```go
type EmailService interface {
    SendEmail(to []string, subject, body string) error
    SendHTMLEmail(to []string, subject, htmlBody string) error
    SendTemplateEmail(to []string, subject, templateName string, data interface{}) error
}
```

## Cài đặt

```go
import "github.com/toji-dev/go-shop/internal/pkg/email"
```

## Cách sử dụng

### 1. Cấu hình cơ bản

```go
// Sử dụng Gmail
config := email.GmailSMTPConfig(
    "your-email@gmail.com",
    "your-app-password", // App password, không phải password thường
    "your-email@gmail.com",
    "Go-Shop System",
)

// Hoặc sử dụng SMTP tùy chỉnh
config := email.CustomSMTPConfig(
    "smtp.example.com", // host
    587,                // port
    "username",         // username
    "password",         // password
    "noreply@example.com", // from
    "Go-Shop",          // from name
    true,               // use TLS
    false,              // use SSL
)

// Thêm đường dẫn template (tùy chọn)
config.TemplatePath = "./templates/email"
```

### 2. Tạo Email Service

```go
emailService, err := email.NewSMTPEmailService(config)
if err != nil {
    log.Fatalf("Failed to create email service: %v", err)
}
```

### 3. Gửi email

```go
// Email text thông thường
err := emailService.SendEmail(
    []string{"user@example.com"},
    "Test Subject",
    "Email body content",
)

// Email HTML
htmlBody := "<h1>Hello</h1><p>This is HTML email</p>"
err := emailService.SendHTMLEmail(
    []string{"user@example.com"},
    "HTML Email",
    htmlBody,
)

// Email sử dụng template
templateData := map[string]interface{}{
    "Name": "John Doe",
    "Date": time.Now().Format("02/01/2006"),
}
err := emailService.SendTemplateEmail(
    []string{"user@example.com"},
    "Welcome to Go-Shop",
    "welcome", // template name
    templateData,
)
```

## Cấu hình Environment Variables

```bash
# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=your-email@gmail.com
SMTP_FROM_NAME=Go-Shop
SMTP_USE_TLS=true
SMTP_USE_SSL=false

# Template Path (optional)
EMAIL_TEMPLATE_PATH=./templates/email
```

## Tích hợp vào các service

### User Service

```go
// trong auth_service.go
type AuthService struct {
    emailService email.EmailService
    // ... other fields
}

func (s *AuthService) RegisterUser(userData UserDTO) error {
    // ... tạo user ...
    
    // Gửi email chào mừng
    err := s.sendWelcomeEmail(user.Email, user.Name)
    if err != nil {
        log.Printf("Failed to send welcome email: %v", err)
        // Không return lỗi vì đăng ký thành công rồi
    }
    
    return nil
}

func (s *AuthService) sendWelcomeEmail(userEmail, userName string) error {
    subject := "Chào mừng bạn đến với Go-Shop!"
    
    data := map[string]interface{}{
        "Name": userName,
        "Date": time.Now().Format("02/01/2006"),
    }
    
    return s.emailService.SendTemplateEmail(
        []string{userEmail},
        subject,
        "welcome",
        data,
    )
}

func (s *AuthService) ForgotPassword(email string) error {
    // ... tạo reset token ...
    
    return s.sendPasswordResetEmail(email, resetLink)
}

func (s *AuthService) sendPasswordResetEmail(userEmail, resetLink string) error {
    subject := "Đặt lại mật khẩu Go-Shop"
    
    data := map[string]interface{}{
        "ResetLink":  resetLink,
        "ExpireTime": "24 giờ",
    }
    
    return s.emailService.SendTemplateEmail(
        []string{userEmail},
        subject,
        "password_reset",
        data,
    )
}
```

### Order Service

```go
// trong order_service.go
type OrderService struct {
    emailService email.EmailService
    // ... other fields
}

func (s *OrderService) CreateOrder(orderData OrderDTO) error {
    // ... tạo đơn hàng ...
    
    // Gửi email xác nhận
    err := s.sendOrderConfirmationEmail(order.CustomerEmail, order.ID, orderData)
    if err != nil {
        log.Printf("Failed to send order confirmation email: %v", err)
    }
    
    return nil
}

func (s *OrderService) sendOrderConfirmationEmail(customerEmail, orderID string, orderData interface{}) error {
    subject := fmt.Sprintf("Xác nhận đơn hàng #%s", orderID)
    
    data := map[string]interface{}{
        "OrderID":   orderID,
        "OrderData": orderData,
        "Date":      time.Now().Format("02/01/2006 15:04"),
    }
    
    return s.emailService.SendTemplateEmail(
        []string{customerEmail},
        subject,
        "order_confirmation",
        data,
    )
}
```

### Notification Service

```go
// trong notification_service.go
type NotificationService struct {
    emailService email.EmailService
    // ... other fields
}

func (s *NotificationService) SendNotification(notification NotificationDTO) error {
    // Gửi qua email
    return s.emailService.SendTemplateEmail(
        []string{notification.UserEmail},
        notification.Subject,
        "notification",
        map[string]interface{}{
            "Message": notification.Message,
            "Date":    time.Now().Format("02/01/2006 15:04"),
        },
    )
}
```

## Email Templates

Service hỗ trợ HTML templates với Go template engine. Đặt các file template trong thư mục được cấu hình:

- `welcome.html` - Template email chào mừng
- `password_reset.html` - Template reset password
- `order_confirmation.html` - Template xác nhận đơn hàng
- `notification.html` - Template thông báo

### Template Variables

#### Welcome Email
- `{{.Name}}` - Tên người dùng
- `{{.Date}}` - Ngày gửi email

#### Password Reset Email
- `{{.ResetLink}}` - Link reset password
- `{{.ExpireTime}}` - Thời gian hết hạn

#### Order Confirmation Email
- `{{.OrderID}}` - Mã đơn hàng
- `{{.Date}}` - Thời gian đặt hàng
- `{{.OrderData}}` - Dữ liệu đơn hàng

#### Notification Email
- `{{.Message}}` - Nội dung thông báo
- `{{.Date}}` - Thời gian gửi

## Cấu hình Gmail

Để sử dụng Gmail SMTP:

1. Bật 2-Factor Authentication cho tài khoản Gmail
2. Tạo App Password:
   - Đi đến Google Account Settings
   - Security → 2-Step Verification → App passwords
   - Tạo password mới cho ứng dụng
3. Sử dụng App Password thay vì password thường

## Testing

Để test email trong môi trường development, có thể sử dụng MailHog:

```bash
# Chạy MailHog
docker run -d -p 1025:1025 -p 8025:8025 mailhog/mailhog

# Cấu hình
config := email.DefaultSMTPConfig() // localhost:1025
```

## Lỗi thường gặp

1. **Authentication failed**: Kiểm tra username/password, với Gmail phải dùng App Password
2. **Connection refused**: Kiểm tra host và port
3. **Template not found**: Kiểm tra đường dẫn template path
4. **TLS/SSL errors**: Kiểm tra cấu hình TLS/SSL cho SMTP server

## Security

- Không bao giờ commit SMTP credentials vào code
- Sử dụng environment variables
- Với Gmail, luôn sử dụng App Password
- Cân nhắc sử dụng email service như SendGrid, AWS SES cho production
