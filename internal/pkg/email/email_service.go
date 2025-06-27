package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
	"strings"
)

// EmailService định nghĩa interface cho email service
type EmailService interface {
	SendEmail(to []string, subject, body string) error
	SendHTMLEmail(to []string, subject, htmlBody string) error
	SendTemplateEmail(to []string, subject, templateName string, data interface{}) error
}

// SMTPConfig cấu hình cho SMTP server
type SMTPConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	From         string `json:"from"`
	FromName     string `json:"from_name"`
	UseTLS       bool   `json:"use_tls"`
	UseSSL       bool   `json:"use_ssl"`
	TemplatePath string `json:"template_path"`
}

// SMTPEmailService implementation của EmailService sử dụng SMTP
type SMTPEmailService struct {
	config    SMTPConfig
	auth      smtp.Auth
	templates map[string]*template.Template
}

// Email struct để định nghĩa cấu trúc email
type Email struct {
	To      []string
	CC      []string
	BCC     []string
	Subject string
	Body    string
	IsHTML  bool
}

// NewSMTPEmailService tạo instance mới của SMTPEmailService
func NewSMTPEmailService(config SMTPConfig) (*SMTPEmailService, error) {
	service := &SMTPEmailService{
		config:    config,
		templates: make(map[string]*template.Template),
	}

	// Thiết lập authentication
	if config.Username != "" && config.Password != "" {
		service.auth = smtp.PlainAuth("", config.Username, config.Password, config.Host)
	}

	// Load email templates nếu có
	if config.TemplatePath != "" {
		err := service.loadTemplates()
		if err != nil {
			return nil, fmt.Errorf("failed to load email templates: %w", err)
		}
	}

	return service, nil
}

// loadTemplates load tất cả email templates từ thư mục template
func (s *SMTPEmailService) loadTemplates() error {
	templateFiles := []string{
		"welcome.html",
		"password_reset.html",
		"order_confirmation.html",
		"notification.html",
	}

	for _, fileName := range templateFiles {
		templatePath := filepath.Join(s.config.TemplatePath, fileName)
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			// Log warning nhưng không fail, template có thể không tồn tại
			continue
		}
		templateName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		s.templates[templateName] = tmpl
	}

	return nil
}

// SendEmail gửi email text thông thường
func (s *SMTPEmailService) SendEmail(to []string, subject, body string) error {
	email := Email{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	}
	return s.sendEmail(email)
}

// SendHTMLEmail gửi email HTML
func (s *SMTPEmailService) SendHTMLEmail(to []string, subject, htmlBody string) error {
	email := Email{
		To:      to,
		Subject: subject,
		Body:    htmlBody,
		IsHTML:  true,
	}
	return s.sendEmail(email)
}

// SendTemplateEmail gửi email sử dụng template
func (s *SMTPEmailService) SendTemplateEmail(to []string, subject, templateName string, data interface{}) error {
	tmpl, exists := s.templates[templateName]
	if !exists {
		return fmt.Errorf("template %s not found", templateName)
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return s.SendHTMLEmail(to, subject, buf.String())
}

// sendEmail method chính để gửi email
func (s *SMTPEmailService) sendEmail(email Email) error {
	// Tạo message
	message := s.buildMessage(email)

	// Địa chỉ server
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Gửi email
	if s.config.UseSSL {
		return s.sendEmailSSL(addr, message, email.To)
	} else if s.config.UseTLS {
		return s.sendEmailTLS(addr, message, email.To)
	} else {
		return smtp.SendMail(addr, s.auth, s.config.From, email.To, []byte(message))
	}
}

// buildMessage xây dựng message email
func (s *SMTPEmailService) buildMessage(email Email) string {
	var message strings.Builder

	// Headers
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.From))
	message.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ",")))

	if len(email.CC) > 0 {
		message.WriteString(fmt.Sprintf("CC: %s\r\n", strings.Join(email.CC, ",")))
	}

	message.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	message.WriteString("MIME-Version: 1.0\r\n")

	if email.IsHTML {
		message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	message.WriteString("\r\n")
	message.WriteString(email.Body)

	return message.String()
}

// sendEmailTLS gửi email với TLS
func (s *SMTPEmailService) sendEmailTLS(addr, message string, to []string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Start TLS
	tlsConfig := &tls.Config{
		ServerName: s.config.Host,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Auth
	if s.auth != nil {
		if err = client.Auth(s.auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Set sender
	if err = client.Mail(s.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	_, err = writer.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// sendEmailSSL gửi email với SSL
func (s *SMTPEmailService) sendEmailSSL(addr, message string, to []string) error {
	tlsConfig := &tls.Config{
		ServerName: s.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect with SSL: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Auth
	if s.auth != nil {
		if err = client.Auth(s.auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Set sender
	if err = client.Mail(s.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	_, err = writer.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// ValidateConfig kiểm tra cấu hình SMTP
func ValidateConfig(config SMTPConfig) error {
	if config.Host == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if config.Port <= 0 {
		return fmt.Errorf("SMTP port must be greater than 0")
	}
	if config.From == "" {
		return fmt.Errorf("from address is required")
	}
	return nil
}
