package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
)

// SendEmail sends a MIME-HTML email via SMTP
func SendEmail(to, subject, htmlBody string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com" // Google SMTP server
	}
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587" // TLS/STARTTLS port
	}
	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		smtpUser = "mysticadityo@gmail.com"
	}
	smtpPass := os.Getenv("SMTP_PASS")
	if smtpPass == "" {
		smtpPass = "kemk rwrs uotr ymth"
	}

	from := "no-reply@ypai.or.id"
	
	// Create headers and body
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")

	// Base64 encode the body to prevent encoding errors
	encodedBody := base64.StdEncoding.EncodeToString([]byte(htmlBody))
	msg.WriteString(encodedBody)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
