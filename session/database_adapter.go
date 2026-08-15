package session

import "time"

type User interface {
	// ID
	GetID() int64
	SetID(id int64)

	// Name
	GetName() string
	SetName(string)

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
	GetID() int64
	SetID(id int64)

	// User relation
	GetUserID() int64
	SetUserID(id int64)

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
	// New instance
	NewUser() User
	NewPasswordResetToken() PasswordResetToken

	// User
	FindUserByEmail(email string) (User, error)
	FindUserByID(id int64) (User, error)
	CreateUser(user User) (int64, error)
	UpdateUser(user User) error

	// Password reset
	CreatePasswordResetToken(token PasswordResetToken) error
	FindPasswordResetToken(tokenHash string) (PasswordResetToken, error)
	DeletePasswordResetToken(id int64) error

	// Cleanup
	DeleteExpiredPasswordResetTokens() error
}