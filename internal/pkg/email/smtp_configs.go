package email

// Các cấu hình SMTP phổ biến

// GmailSMTPConfig cấu hình cho Gmail SMTP
func GmailSMTPConfig(username, password, from, fromName string) SMTPConfig {
	return SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: username,
		Password: password,
		From:     from,
		FromName: fromName,
		UseTLS:   true,
		UseSSL:   false,
	}
}

// OutlookSMTPConfig cấu hình cho Outlook/Hotmail SMTP
func OutlookSMTPConfig(username, password, from, fromName string) SMTPConfig {
	return SMTPConfig{
		Host:     "smtp.live.com",
		Port:     587,
		Username: username,
		Password: password,
		From:     from,
		FromName: fromName,
		UseTLS:   true,
		UseSSL:   false,
	}
}

// YahooSMTPConfig cấu hình cho Yahoo SMTP
func YahooSMTPConfig(username, password, from, fromName string) SMTPConfig {
	return SMTPConfig{
		Host:     "smtp.mail.yahoo.com",
		Port:     587,
		Username: username,
		Password: password,
		From:     from,
		FromName: fromName,
		UseTLS:   true,
		UseSSL:   false,
	}
}

// CustomSMTPConfig cấu hình cho SMTP server tùy chỉnh
func CustomSMTPConfig(host string, port int, username, password, from, fromName string, useTLS, useSSL bool) SMTPConfig {
	return SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		FromName: fromName,
		UseTLS:   useTLS,
		UseSSL:   useSSL,
	}
}

// DefaultSMTPConfig cấu hình mặc định (cho development)
func DefaultSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:     "localhost",
		Port:     1025, // MailHog default port
		Username: "",
		Password: "",
		From:     "noreply@go-shop.com",
		FromName: "Go-Shop",
		UseTLS:   false,
		UseSSL:   false,
	}
}
