package mysql

import (
	"encoding/json"
	"time"
)

// User is the concrete, MariaDB/MySQL-backed implementation of session.User.
type User struct {
	ID                     int64
	Name                   string
	Email                  string
	EmailVerifiedAt        *time.Time
	Password               string
	RememberToken          string
	TwoFactorSecret        string
	TwoFactorRecoveryCodes []string
	TwoFactorConfirmedAt   *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func (u *User) GetID() int64   { return u.ID }
func (u *User) SetID(id int64) { u.ID = id }

func (u *User) GetName() string   { return u.Name }
func (u *User) SetName(name string) { u.Name = name }

func (u *User) GetEmail() string            { return u.Email }
func (u *User) SetEmail(email string)       { u.Email = email }
func (u *User) GetPassword() string         { return u.Password }
func (u *User) SetPassword(password string) { u.Password = password }

func (u *User) GetEmailVerifiedAt() *time.Time  { return u.EmailVerifiedAt }
func (u *User) SetEmailVerifiedAt(t *time.Time) { u.EmailVerifiedAt = t }

func (u *User) GetRememberToken() string      { return u.RememberToken }
func (u *User) SetRememberToken(token string) { u.RememberToken = token }

func (u *User) GetTwoFactorSecret() string       { return u.TwoFactorSecret }
func (u *User) SetTwoFactorSecret(secret string) { u.TwoFactorSecret = secret }

func (u *User) GetTwoFactorRecoveryCodes() []string      { return u.TwoFactorRecoveryCodes }
func (u *User) SetTwoFactorRecoveryCodes(codes []string) { u.TwoFactorRecoveryCodes = codes }

func (u *User) GetTwoFactorConfirmedAt() *time.Time  { return u.TwoFactorConfirmedAt }
func (u *User) SetTwoFactorConfirmedAt(t *time.Time) { u.TwoFactorConfirmedAt = t }

func marshalRecoveryCodes(codes []string) (string, error) {
	if codes == nil {
		codes = []string{}
	}
	b, err := json.Marshal(codes)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func unmarshalRecoveryCodes(raw string) ([]string, error) {
	if raw == "" {
		return []string{}, nil
	}
	var codes []string
	if err := json.Unmarshal([]byte(raw), &codes); err != nil {
		return nil, err
	}
	return codes, nil
}