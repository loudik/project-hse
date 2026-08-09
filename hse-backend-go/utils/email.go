package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(toEmail, subject, htmlBody string) error {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return nil // not configured - caller should fall back to another way of surfacing the content (e.g. dev response field)
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = username
	}
	fromName := os.Getenv("SMTP_FROM_NAME")
	if fromName == "" {
		fromName = "HSE Approval System"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", username, password, host)

	msg := fmt.Sprintf(
		"From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		fromName, fromEmail, toEmail, subject, htmlBody,
	)

	return smtp.SendMail(addr, auth, fromEmail, []string{toEmail}, []byte(msg))
}
