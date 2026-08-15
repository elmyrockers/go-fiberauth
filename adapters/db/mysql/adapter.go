// Package mysql implements session.DatabaseAdapter on top of MariaDB/MySQL,
// using database/sql and the go-sql-driver/mysql driver.
package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Adapter struct {
	db *sql.DB
}

// New() opens a connection pool from cfg and verifies it with a ping.
func New(cfg Config) (*Adapter, error) {
	db, err := sql.Open("mysql", cfg.dsn())
	if err != nil {
		return nil, fmt.Errorf("mysql: open: %w", err)
	}

	maxOpen := cfg.MaxOpenConns
	if maxOpen == 0 {
		maxOpen = 10
	}
	maxIdle := cfg.MaxIdleConns
	if maxIdle == 0 {
		maxIdle = 5
	}
	connMaxLifetime := cfg.ConnMaxLifetime
	if connMaxLifetime == 0 {
		connMaxLifetime = 30 * time.Minute
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(connMaxLifetime)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("mysql: ping: %w", err)
	}

	return &Adapter{db: db}, nil
}

func (a *Adapter) NewUser() *User {
	return &User{}
}

func (a *Adapter) NewPasswordResetToken() *PasswordResetToken {
	return &PasswordResetToken{}
}

func (a *Adapter) FindUserByEmail(email string) (User, error) {
	const q = `
		SELECT id, name, email, email_verified_at, password, remember_token,
		       two_factor_secret, two_factor_recovery_codes, two_factor_confirmed_at,
		       created_at, updated_at
		FROM users
		WHERE email = ?
		LIMIT 1`
	return a.scanUser(a.db.QueryRow(q, email))
}

func (a *Adapter) FindUserByID(id string) (User, error) {
	const q = `
		SELECT id, name, email, email_verified_at, password, remember_token,
		       two_factor_secret, two_factor_recovery_codes, two_factor_confirmed_at,
		       created_at, updated_at
		FROM users
		WHERE id = ?
		LIMIT 1`
	return a.scanUser(a.db.QueryRow(q, id))
}

func (a *Adapter) scanUser(row *sql.Row) (User, error) {
	var (
		u                      User
		name                   sql.NullString
		emailVerifiedAt        sql.NullTime
		rememberToken          sql.NullString
		twoFactorSecret        sql.NullString
		twoFactorRecoveryCodes sql.NullString
		twoFactorConfirmedAt   sql.NullTime
	)

	err := row.Scan(
		&u.ID, &name, &u.Email, &emailVerifiedAt, &u.Password, &rememberToken,
		&twoFactorSecret, &twoFactorRecoveryCodes, &twoFactorConfirmedAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, err
	}

	u.Name = name.String
	u.RememberToken = rememberToken.String
	u.TwoFactorSecret = twoFactorSecret.String

	if emailVerifiedAt.Valid {
		t := emailVerifiedAt.Time
		u.EmailVerifiedAt = &t
	}
	if twoFactorConfirmedAt.Valid {
		t := twoFactorConfirmedAt.Time
		u.TwoFactorConfirmedAt = &t
	}

	codes, err := unmarshalRecoveryCodes(twoFactorRecoveryCodes.String)
	if err != nil {
		return User{}, err
	}
	u.TwoFactorRecoveryCodes = codes

	return u, nil
}

func (a *Adapter) CreateUser(user User) (int64, error) {
	codesJSON, err := marshalRecoveryCodes(user.TwoFactorRecoveryCodes)
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()

	const q = `
		INSERT INTO users (
			name, email, email_verified_at, password, remember_token,
			two_factor_secret, two_factor_recovery_codes, two_factor_confirmed_at,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := a.db.Exec(q,
		user.Name,
		user.Email,
		nullableTime(user.EmailVerifiedAt),
		user.Password,
		nullableString(user.RememberToken),
		nullableString(user.TwoFactorSecret),
		codesJSON,
		nullableTime(user.TwoFactorConfirmedAt),
		now,
		now,
	)
	if isDuplicateEntry(err) {
		return 0, ErrDuplicateEmail
	}
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (a *Adapter) UpdateUser(user User) error {
	codesJSON, err := marshalRecoveryCodes(user.TwoFactorRecoveryCodes)
	if err != nil {
		return err
	}

	const q = `
		UPDATE users SET
			name = ?,
			email = ?,
			email_verified_at = ?,
			password = ?,
			remember_token = ?,
			two_factor_secret = ?,
			two_factor_recovery_codes = ?,
			two_factor_confirmed_at = ?,
			updated_at = ?
		WHERE id = ?`

	res, err := a.db.Exec(q,
		user.Name,
		user.Email,
		nullableTime(user.EmailVerifiedAt),
		user.Password,
		nullableString(user.RememberToken),
		nullableString(user.TwoFactorSecret),
		codesJSON,
		nullableTime(user.TwoFactorConfirmedAt),
		time.Now().UTC(),
		user.ID,
	)
	if isDuplicateEntry(err) {
		return ErrDuplicateEmail
	}
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (a *Adapter) CreatePasswordResetToken(token PasswordResetToken) error {
	if token.ID == "" {
		token.ID = uuid.NewString()
	}

	const q = `
		INSERT INTO password_reset_tokens (
			id, user_id, token_hash, expires_at, used_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := a.db.Exec(q,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt.UTC(),
		nullableTime(token.UsedAt),
		time.Now().UTC(),
	)
	return err
}

func (a *Adapter) FindPasswordResetToken(tokenHash string) (PasswordResetToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token_hash = ?
		LIMIT 1`

	var (
		t         PasswordResetToken
		usedAt    sql.NullTime
		createdAt time.Time
	)

	err := a.db.QueryRow(q, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &usedAt, &createdAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return PasswordResetToken{}, ErrNotFound
	}
	if err != nil {
		return PasswordResetToken{}, err
	}

	t.CreatedAt = createdAt
	if usedAt.Valid {
		u := usedAt.Time
		t.UsedAt = &u
	}

	return t, nil
}

func (a *Adapter) DeletePasswordResetToken(id string) error {
	res, err := a.db.Exec(`DELETE FROM password_reset_tokens WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (a *Adapter) DeleteExpiredPasswordResetTokens() error {
	_, err := a.db.Exec(`DELETE FROM password_reset_tokens WHERE expires_at < ?`, time.Now().UTC())
	return err
}

// Close closes the underlying connection pool.
func (a *Adapter) Close() error {
	return a.db.Close()
}

func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nullableTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC()
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}