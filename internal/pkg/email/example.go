package email

import (
	"log"
	"os"
)

// Example: Cách sử dụng SMTP Email Service

func ExampleUsage() {
	// 1. Tạo cấu hình SMTP
	config := GmailSMTPConfig(
		os.Getenv("SMTP_USERNAME"), // your-email@gmail.com
		os.Getenv("SMTP_PASSWORD"), // your-app-password
		os.Getenv("SMTP_FROM"),     // your-email@gmail.com
		"Go-Shop System",           // Display name
	)

	// Thêm đường dẫn template nếu có
	config.TemplatePath = "./templates/email"

	// 2. Tạo email service
	emailService, err := NewSMTPEmailService(config)
	if err != nil {
		log.Fatalf("Failed to create email service: %v", err)
	}

	// 3. Sử dụng các method khác nhau

	// Gửi email text thông thường
	err = emailService.SendEmail(
		[]string{"user@example.com"},
		"Test Email",
		"This is a test email body.",
	)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
	}

	// Gửi email HTML
	htmlBody := `
		<html>
		<body>
			<h1>Welcome to Go-Shop!</h1>
			<p>Thank you for joining us.</p>
		</body>
		</html>
	`
	err = emailService.SendHTMLEmail(
		[]string{"user@example.com"},
		"Welcome to Go-Shop",
		htmlBody,
	)
	if err != nil {
		log.Printf("Failed to send HTML email: %v", err)
	}

	// Gửi email chào mừng
	err = emailService.SendWelcomeEmail("user@example.com", "John Doe")
	if err != nil {
		log.Printf("Failed to send welcome email: %v", err)
	}

	// Gửi email reset password
	resetLink := "https://go-shop.com/reset-password?token=abc123"
	err = emailService.SendPasswordResetEmail("user@example.com", resetLink)
	if err != nil {
		log.Printf("Failed to send password reset email: %v", err)
	}

	// Gửi email xác nhận đơn hàng
	orderData := map[string]interface{}{
		"CustomerName": "John Doe",
		"TotalAmount":  250000,
		"Items": []map[string]interface{}{
			{"Name": "Product 1", "Quantity": 2, "Price": 100000},
			{"Name": "Product 2", "Quantity": 1, "Price": 150000},
		},
	}
	err = emailService.SendOrderConfirmationEmail("user@example.com", "ORD001", orderData)
	if err != nil {
		log.Printf("Failed to send order confirmation email: %v", err)
	}

	// Gửi email thông báo
	err = emailService.SendNotificationEmail(
		"user@example.com",
		"System Maintenance",
		"The system will be under maintenance from 2:00 AM to 4:00 AM.",
	)
	if err != nil {
		log.Printf("Failed to send notification email: %v", err)
	}
}

// Cách tích hợp vào từng service:

// 1. Trong user-service (auth_service.go)
func IntegrateWithUserService() {
	// config := email.GmailSMTPConfig(...)
	// emailService, _ := email.NewSMTPEmailService(config)

	// Trong method register user:
	// emailService.SendWelcomeEmail(user.Email, user.Name)

	// Trong method forgot password:
	// emailService.SendPasswordResetEmail(user.Email, resetLink)
}

// 2. Trong order-service
func IntegrateWithOrderService() {
	// config := email.GmailSMTPConfig(...)
	// emailService, _ := email.NewSMTPEmailService(config)

	// Khi tạo đơn hàng thành công:
	// emailService.SendOrderConfirmationEmail(order.CustomerEmail, order.ID, orderData)
}

// 3. Trong notification-service
func IntegrateWithNotificationService() {
	// config := email.GmailSMTPConfig(...)
	// emailService, _ := email.NewSMTPEmailService(config)

	// Khi có thông báo mới:
	// emailService.SendNotificationEmail(user.Email, notification.Subject, notification.Message)
}

// Environment Variables cần thiết:
// SMTP_HOST=smtp.gmail.com
// SMTP_PORT=587
// SMTP_USERNAME=your-email@gmail.com
// SMTP_PASSWORD=your-app-password
// SMTP_FROM=your-email@gmail.com
// SMTP_FROM_NAME=Go-Shop
// SMTP_USE_TLS=true
// SMTP_USE_SSL=false
// EMAIL_TEMPLATE_PATH=./templates/email
