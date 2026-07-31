package session

import "time"

type User interface {
	// ID
	GetID() string
	SetID(id string)

	// Credentials
	GetEmail() string
	SetEmail(email string)

	GetPassword() string
	SetPassword(password string)

	// Email verification
	GetEmailVerifiedAt() *time.Time
	SetEmailVerifiedAt(t *time.Time)

	// Remember me
	GetRememberToken() string
	SetRememberToken(token string)

	// Two factor authentication
	GetTwoFactorSecret() string
	SetTwoFactorSecret(secret string)

	GetTwoFactorRecoveryCodes() []string
	SetTwoFactorRecoveryCodes(codes []string)

	GetTwoFactorConfirmedAt() *time.Time
	SetTwoFactorConfirmedAt(t *time.Time)
}

type PasswordResetToken interface {
	// ID
	GetID() string
	SetID(id string)

	// User relation
	GetUserID() string
	SetUserID(id string)

	// Token
	GetTokenHash() string
	SetTokenHash(hash string)

	// Expiry
	GetExpiresAt() time.Time
	SetExpiresAt(t time.Time)

	// Usage
	GetUsedAt() *time.Time
	SetUsedAt(t *time.Time)
}

type DatabaseAdapter interface {
	// User
	FindUserByEmail(email string) (User, error)
	FindUserByID(id string) (User, error)
	CreateUser(user User) error
	UpdateUser(user User) error

	// Password reset
	CreatePasswordResetToken(token PasswordResetToken) error
	FindPasswordResetToken(tokenHash string) (PasswordResetToken, error)
	DeletePasswordResetToken(id string) error

	// Cleanup
	DeleteExpiredPasswordResetTokens() error
}