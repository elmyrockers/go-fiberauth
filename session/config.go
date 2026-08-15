package session

import (
	"crypto/rsa"
	"github.com/gofiber/fiber/v3/middleware/session"
)

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
	// For JWT signature
	HMACKey           []byte
	RSAPrivateKeyPath string
	RSAPublicKeyPath  string

	SessionConfig session.Config
	DatabaseAdapter DatabaseAdapter
	MailAdapter MailAdapter
	MailConfig MailConfig
	LoginURL string
}

func (c *Config) LoadPrivateKey() (*rsa.PrivateKey, error) {
	if c.RSAPrivateKeyPath == "" {
		return nil, fmt.Errorf("fiberauth: RSAPrivateKeyPath is empty")
	}

	data, err := os.ReadFile(c.RSAPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("fiberauth: read private key: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("fiberauth: invalid PEM in private key file %q", c.RSAPrivateKeyPath)
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("fiberauth: parse private key %q: %w", c.RSAPrivateKeyPath, err)
	}

	rsaKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("fiberauth: %q does not contain an RSA private key", c.RSAPrivateKeyPath)
	}

	return rsaKey, nil
}

func (c *Config) LoadPublicKey() (*rsa.PublicKey, error) {
	if c.RSAPublicKeyPath == "" {
		return nil, fmt.Errorf("fiberauth: RSAPublicKeyPath is empty")
	}

	data, err := os.ReadFile(c.RSAPublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("fiberauth: read public key: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("fiberauth: invalid PEM in public key file %q", c.RSAPublicKeyPath)
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("fiberauth: parse public key %q: %w", c.RSAPublicKeyPath, err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("fiberauth: %q does not contain an RSA public key", c.RSAPublicKeyPath)
	}

	return rsaPub, nil
}