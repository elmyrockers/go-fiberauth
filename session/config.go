package session

import (
	"github.com/gofiber/fiber/v3/middleware/session"
)

type MailAddress struct {
	Name  string
	Email string
}

const (
	EmailVerification = "email_verification"
	PasswordReset = "password_reset"
)

type MailTemplate struct {
    Subject string
    Body    string
    AltBody string
}



type MailConfig struct {
	From MailAddress
	ReplyTo *MailAddress
	Templates map[string]MailTemplate
}




type Config struct {
	SessionConfig session.Config
	DatabaseAdapter DatabaseAdapter
	MailAdapter MailAdapter
	MailConfig MailConfig
}