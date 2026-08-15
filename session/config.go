package session

import (
	"crypto/rsa"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// Constant For JWT
const (
	EmailVerification = "email_verification"
	PasswordReset = "password_reset"
)

type MailConfig struct {
	From MailAddress
	ReplyTo *MailAddress
	Templates map[string]MailTemplate
}

type Config struct {
	// JWT signing
	SharedKey      []byte // used when Algorithm is symmetric (HS256)
	PrivateKeyPath string // used when Algorithm is asymmetric (RS256)
	PublicKeyPath  string

	SessionConfig session.Config
	DatabaseAdapter DatabaseAdapter
	MailAdapter MailAdapter
	MailConfig MailConfig
	LoginURL string
}

func (c *Config) LoadPrivateKey() (*rsa.PrivateKey, error) {
	if c.PrivateKeyPath == "" {
		return nil, fmt.Errorf("fiberauth: PrivateKeyPath is empty")
	}

	data, err := os.ReadFile(c.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("fiberauth: read private key: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("fiberauth: invalid PEM in private key file %q", c.PrivateKeyPath)
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("fiberauth: parse private key %q: %w", c.PrivateKeyPath, err)
	}

	rsaKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("fiberauth: %q does not contain an RSA private key", c.PrivateKeyPath)
	}

	return rsaKey, nil
}

func (c *Config) LoadPublicKey() (*rsa.PublicKey, error) {
	if c.PublicKeyPath == "" {
		return nil, fmt.Errorf("fiberauth: PublicKeyPath is empty")
	}

	data, err := os.ReadFile(c.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("fiberauth: read public key: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("fiberauth: invalid PEM in public key file %q", c.PublicKeyPath)
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("fiberauth: parse public key %q: %w", c.PublicKeyPath, err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("fiberauth: %q does not contain an RSA public key", c.PublicKeyPath)
	}

	return rsaPub, nil
}