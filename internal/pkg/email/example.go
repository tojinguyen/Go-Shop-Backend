package email

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Example: Cách sử dụng SMTP Email Service cơ bản

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

	// 3. Sử dụng các method cơ bản

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

	// Gửi email sử dụng template
	templateData := map[string]interface{}{
		"Name": "John Doe",
		"Date": time.Now().Format("02/01/2006"),
	}
	err = emailService.SendTemplateEmail(
		[]string{"user@example.com"},
		"Welcome to Go-Shop",
		"welcome", // template name
		templateData,
	)
	if err != nil {
		log.Printf("Failed to send template email: %v", err)
	}
}

// Ví dụ implement Welcome Email trong User Service
func ExampleWelcomeEmailInUserService() {
	// Trong user service, bạn có thể implement như sau:

	config := GmailSMTPConfig(
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_FROM"),
		"Go-Shop",
	)
	config.TemplatePath = "./templates/email"

	emailService, _ := NewSMTPEmailService(config)

	// Method gửi welcome email
	sendWelcomeEmail := func(userEmail, userName string) error {
		subject := "Chào mừng bạn đến với Go-Shop!"

		// Sử dụng template nếu có
		data := map[string]interface{}{
			"Name": userName,
			"Date": time.Now().Format("02/01/2006"),
		}

		return emailService.SendTemplateEmail(
			[]string{userEmail},
			subject,
			"welcome",
			data,
		)
	}

	// Sử dụng
	err := sendWelcomeEmail("user@example.com", "John Doe")
	if err != nil {
		log.Printf("Failed to send welcome email: %v", err)
	}
}

// Ví dụ implement Password Reset Email trong Auth Service
func ExamplePasswordResetEmailInAuthService() {
	config := GmailSMTPConfig(
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_FROM"),
		"Go-Shop",
	)
	config.TemplatePath = "./templates/email"

	emailService, _ := NewSMTPEmailService(config)

	// Method gửi password reset email
	sendPasswordResetEmail := func(userEmail, resetLink string) error {
		subject := "Đặt lại mật khẩu Go-Shop"

		// Sử dụng template nếu có
		data := map[string]interface{}{
			"ResetLink":  resetLink,
			"ExpireTime": "24 giờ",
		}

		return emailService.SendTemplateEmail(
			[]string{userEmail},
			subject,
			"password_reset",
			data,
		)
	}

	// Sử dụng
	resetLink := "https://go-shop.com/reset-password?token=abc123"
	err := sendPasswordResetEmail("user@example.com", resetLink)
	if err != nil {
		log.Printf("Failed to send password reset email: %v", err)
	}
}

// Ví dụ implement Order Confirmation Email trong Order Service
func ExampleOrderConfirmationEmailInOrderService() {
	config := GmailSMTPConfig(
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_FROM"),
		"Go-Shop",
	)
	config.TemplatePath = "./templates/email"

	emailService, _ := NewSMTPEmailService(config)

	// Method gửi order confirmation email
	sendOrderConfirmationEmail := func(customerEmail, orderID string, orderData interface{}) error {
		subject := fmt.Sprintf("Xác nhận đơn hàng #%s", orderID)

		// Sử dụng template nếu có
		data := map[string]interface{}{
			"OrderID":   orderID,
			"OrderData": orderData,
			"Date":      time.Now().Format("02/01/2006 15:04"),
		}

		return emailService.SendTemplateEmail(
			[]string{customerEmail},
			subject,
			"order_confirmation",
			data,
		)
	}

	// Sử dụng
	orderData := map[string]interface{}{
		"CustomerName": "John Doe",
		"TotalAmount":  250000,
		"Items": []map[string]interface{}{
			{"Name": "Product 1", "Quantity": 2, "Price": 100000},
			{"Name": "Product 2", "Quantity": 1, "Price": 150000},
		},
	}
	err := sendOrderConfirmationEmail("user@example.com", "ORD001", orderData)
	if err != nil {
		log.Printf("Failed to send order confirmation email: %v", err)
	}
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
