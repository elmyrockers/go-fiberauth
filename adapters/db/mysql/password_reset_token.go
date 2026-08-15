package mysql

import "time"

type PasswordResetToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func (t *PasswordResetToken) GetID() int64   { return t.ID }
func (t *PasswordResetToken) SetID(id int64) { t.ID = id }

func (t *PasswordResetToken) GetUserID() int64   { return t.UserID }
func (t *PasswordResetToken) SetUserID(id int64) { t.UserID = id }

func (t *PasswordResetToken) GetTokenHash() string     { return t.TokenHash }
func (t *PasswordResetToken) SetTokenHash(hash string) { t.TokenHash = hash }

func (t *PasswordResetToken) GetExpiresAt() time.Time   { return t.ExpiresAt }
func (t *PasswordResetToken) SetExpiresAt(tm time.Time) { t.ExpiresAt = tm }

func (t *PasswordResetToken) GetUsedAt() *time.Time   { return t.UsedAt }
func (t *PasswordResetToken) SetUsedAt(tm *time.Time) { t.UsedAt = tm }