# Auth Routes Plan

All routes grouped under `/auth` prefix.

## Auth core
| Method | Route           | Purpose            |
|--------|-----------------|---------------------|
| POST   | /auth/login     | Log in (accepts optional `remember: bool`) |
| POST   | /auth/logout    | Log out             |
| POST   | /auth/register  | Register new user    |

## Profile & password management
| Method | Route           | Purpose                                        |
|--------|-----------------|-------------------------------------------------|
| PUT    | /auth/profile   | Update name/email (authenticated)               |
| PUT    | /auth/password  | Change password (authenticated, requires current password) |

## Password reset
| Method | Route                  | Purpose                        |
|--------|--------------------------|-----------------------------------|
| POST   | /auth/forgot-password  | Request reset link (send email)  |
| POST   | /auth/reset-password   | Submit new password + token       |

## Email verification
| Method | Route                          | Purpose                     |
|--------|-----------------------------------|--------------------------------|
| POST   | /auth/email/verify-send        | Resend verification email    |
| GET    | /auth/email/verify/:id/:token  | Verify email via link click  |

## Password confirmation (optional)
| Method | Route                        | Purpose                       |
|--------|--------------------------------|-----------------------------------|
| GET    | /auth/password/confirm-status | Check if recently confirmed     |
| POST   | /auth/password/confirm        | Re-enter password to confirm     |

## Two-factor auth (optional)
| Method | Route                    | Purpose                          |
|--------|----------------------------|------------------------------------|
| POST   | /auth/2fa/enable          | Generate secret + QR code          |
| POST   | /auth/2fa/confirm         | Confirm setup with a code          |
| POST   | /auth/2fa/disable         | Turn off 2FA                       |
| POST   | /auth/2fa/challenge       | Verify code during login (pending) |
| GET    | /auth/2fa/recovery-codes  | View recovery codes                |
| POST   | /auth/2fa/recovery-codes  | Regenerate recovery codes          |

**Total: 17 routes** (9 if skipping password confirmation + 2FA for v1)


# Middleware Plan

## Middleware List

| # | Middleware                    | Purpose                                                                 |
|---|---------------------------------|----------------------------------------------------------------------------|
| 1 | GuestOnly                     | Blocks access if user is already logged in                              |
| 2 | AuthRequired                   | Blocks access if not logged in; falls back to remember-token cookie (if present and valid) to re-establish a session before rejecting |
| 3 | RateLimiter                    | Throttles requests (reusable, configurable per route: max + window)      |
| 4 | RequiresTwoFactorPending        | Only allows access if login is mid-2FA-challenge (temp session state)     |
| 5 | RequiresPasswordConfirmed       | Blocks sensitive actions unless password was recently re-confirmed        |
| 6 | CSRF                           | Protects state-changing requests (Fiber's built-in `csrf.New()`)           |

## Route → Middleware Mapping

| Route                              | Middleware                                      |
|--------------------------------------|-----------------------------------------------------|
| POST /auth/login                   | GuestOnly, RateLimiter, CSRF                       |
| POST /auth/logout                  | AuthRequired, CSRF                                  |
| POST /auth/register                | GuestOnly, RateLimiter, CSRF                        |
| PUT /auth/profile                  | AuthRequired, CSRF                                  |
| PUT /auth/password                 | AuthRequired, RequiresPasswordConfirmed, CSRF        |
| POST /auth/forgot-password          | GuestOnly, RateLimiter, CSRF                        |
| POST /auth/reset-password           | GuestOnly, CSRF                                     |
| POST /auth/email/verify-send        | AuthRequired, RateLimiter, CSRF                     |
| GET /auth/email/verify/:id/:token   | (none)                                              |
| GET /auth/password/confirm-status   | AuthRequired                                        |
| POST /auth/password/confirm         | AuthRequired, CSRF                                  |
| POST /auth/2fa/enable               | AuthRequired, RequiresPasswordConfirmed, CSRF        |
| POST /auth/2fa/confirm              | AuthRequired, CSRF                                  |
| POST /auth/2fa/disable              | AuthRequired, RequiresPasswordConfirmed, CSRF        |
| POST /auth/2fa/challenge            | RequiresTwoFactorPending, RateLimiter, CSRF          |
| GET /auth/2fa/recovery-codes        | AuthRequired                                        |
| POST /auth/2fa/recovery-codes       | AuthRequired, RequiresPasswordConfirmed, CSRF        |


# Actions Plan

Actions are methods on the `Auth` struct — business logic (DB/session
writes), called from inside handlers. Not to be confused with middleware,
which only does gatekeeping (session reads, no DB access).

## Actions List

| # | Action                          | Touches                        |
|---|-----------------------------------|-------------------------------------|
| 1 | VerifyCredentials                | DB (read `users`)                  |
| 2 | CreateSession                    | Session (write `user_id`)          |
| 3 | CreatePendingTwoFactorSession      | Session (write `pending_2fa_user_id`) |
| 4 | CreateRememberToken               | DB (write `users.remember_token`) + sets long-lived remember cookie (separate from session cookie) |
| 5 | DestroySession                    | Session (+ clears remember token/cookie on explicit logout) |
| 6 | CreateUser                        | DB (write `users`)                  |
| 7 | GenerateToken                     | DB (write `tokens`)                  |
| 8 | VerifyToken                       | DB (read `tokens`)                   |
| 9 | UpdatePassword (via reset token)  | DB (write `users`)                   |
| 10| UpdatePasswordAuthenticated       | DB (read+write `users`) — verifies current password before writing new one |
| 11| UpdateProfileInfo                 | DB (write `users`) — name/email; if email changes, clears `email_verified_at` and triggers reverification |
| 12| MarkEmailVerified                 | DB (write `users`)                    |
| 13| ConfirmPassword                   | DB (read `users`) + Session (write `confirmed_at`) |
| 14| GenerateTotpSecret                | DB (write `users`, unconfirmed)       |
| 15| VerifyTotpCode                    | DB (read secret)                      |
| 16| EnableTwoFactor                   | DB (write `users`)                    |
| 17| DisableTwoFactor                  | DB (write `users`)                    |
| 18| GenerateRecoveryCodes             | DB (write `users`)                     |

## Handler → Action Mapping

| Handler                | Actions called                                              |
|---------------------------|------------------------------------------------------------------|
| Login                   | VerifyCredentials → CreateSession (or CreatePendingTwoFactorSession if 2FA enabled) → CreateRememberToken (if `remember=true`) |
| Logout                  | DestroySession                                                    |
| Register                | CreateUser → CreateSession                                        |
| UpdateProfile           | UpdateProfileInfo (→ GenerateToken + send email if email changed) |
| UpdatePassword          | UpdatePasswordAuthenticated                                       |
| ForgotPassword           | GenerateToken                                                     |
| ResetPassword            | VerifyToken → UpdatePassword                                      |
| EmailVerifySend          | GenerateToken                                                     |
| EmailVerify              | VerifyToken → MarkEmailVerified                                   |
| PasswordConfirm          | ConfirmPassword                                                   |
| Enable2FA                | GenerateTotpSecret                                                |
| Confirm2FA               | VerifyTotpCode → EnableTwoFactor → GenerateRecoveryCodes           |
| Disable2FA               | DisableTwoFactor                                                  |
| Challenge2FA (login)     | VerifyTotpCode → CreateSession → CreateRememberToken (if `remember=true`) |
| RecoveryCodes (regen)    | GenerateRecoveryCodes                                              |

**Total: 18 distinct actions**, some reused across multiple handlers
(e.g. `CreateSession` used by both Login and 2FA Challenge; `VerifyTotpCode`
used by both 2FA Confirm and 2FA Challenge; `CreateRememberToken` used by
both Login and 2FA Challenge when `remember=true`).


# Database Table

```sql
-- users
CREATE TABLE users (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                        VARCHAR(255) NOT NULL,
    email                       VARCHAR(255) NOT NULL UNIQUE,
    email_verified_at           TIMESTAMPTZ,
    password                    VARCHAR(255) NOT NULL,
    remember_token              VARCHAR(100),
    two_factor_secret           TEXT,           -- encrypted
    two_factor_recovery_codes   TEXT,           -- JSON array, hashed entries
    two_factor_confirmed_at     TIMESTAMPTZ,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_email ON users(email);

-- tokens (password reset + email verification, shared table via `purpose`)
CREATE TABLE tokens (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash   VARCHAR(255) NOT NULL,
    purpose      VARCHAR(50) NOT NULL,   -- 'password_reset' | 'email_verification'
    expires_at   TIMESTAMPTZ NOT NULL,
    used_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_hash_purpose ON tokens(token_hash, purpose);

-- sessions (only needed if we're NOT using Fiber's in-memory/Redis session store as-is)
CREATE TABLE sessions (
    id                  VARCHAR(255) PRIMARY KEY,
    user_id             UUID REFERENCES users(id) ON DELETE CASCADE,  -- nullable: pending-2FA sessions
    pending_2fa_user_id UUID REFERENCES users(id) ON DELETE CASCADE,  -- nullable: only set mid-2FA-challenge
    confirmed_at        TIMESTAMPTZ,   -- for RequiresPasswordConfirmed
    data                JSONB,         -- any extra session payload
    expires_at          TIMESTAMPTZ NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
```


# Package Plan

```
Package: session
Struct: Auth
Constructor: New() function
Actions: Auth methods
            - Middleware:
                - list of middleware
            - Service:
                - list of actions
```

## Usage

```go
import (
    "github.com/gofiber/fiber/v3"
    fiberauth "github.com/elmyrockers/go-fiberauth/session"
    "app/controller"
)

func main() {
    auth := fiberauth.New()
    authController := controller.NewAuth(auth)

    app := fiber.New()

    app.Post("/auth/register", auth.GuestOnly, auth.RateLimiter(5, time.Minute), authController.register)
    app.Post("/auth/login", auth.GuestOnly, auth.RateLimiter(5, time.Minute), authController.login)
    app.Post("/auth/logout", auth.AuthRequired, authController.logout)

    app.Put("/auth/profile", auth.AuthRequired, authController.updateProfile)
    app.Put("/auth/password", auth.AuthRequired, auth.RequiresPasswordConfirmed, authController.updatePassword)

    app.Post("/auth/forgot-password", auth.GuestOnly, auth.RateLimiter(3, time.Hour), authController.forgotPassword)
    app.Post("/auth/reset-password", auth.GuestOnly, authController.resetPassword)

    app.Post("/auth/email/verify-send", auth.AuthRequired, auth.RateLimiter(3, time.Hour), authController.emailVerifySend)
    app.Get("/auth/email/verify/:id/:token", authController.emailVerify)

    app.Get("/auth/password/confirm-status", auth.AuthRequired, authController.passwordConfirmStatus)
    app.Post("/auth/password/confirm", auth.AuthRequired, authController.passwordConfirm)

    app.Post("/auth/2fa/enable", auth.AuthRequired, auth.RequiresPasswordConfirmed, authController.enable2FA)
    app.Post("/auth/2fa/confirm", auth.AuthRequired, authController.confirm2FA)
    app.Post("/auth/2fa/disable", auth.AuthRequired, auth.RequiresPasswordConfirmed, authController.disable2FA)
    app.Post("/auth/2fa/challenge", auth.RequiresTwoFactorPending, auth.RateLimiter(5, time.Minute), authController.challenge2FA)
    app.Get("/auth/2fa/recovery-codes", auth.AuthRequired, authController.recoveryCodes)
    app.Post("/auth/2fa/recovery-codes", auth.AuthRequired, auth.RequiresPasswordConfirmed, authController.regenerateRecoveryCodes)

    app.Listen(":3000")
}
```

## Controller → Action mapping

```
authController.register:
    - calls auth.CreateUser(...)
    - calls auth.CreateSession(...)

authController.login:
    - calls auth.VerifyCredentials(...)
    - calls auth.CreateSession(...) or auth.CreatePendingTwoFactorSession(...)
    - calls auth.CreateRememberToken(...) if remember=true

authController.logout:
    - calls auth.DestroySession(...)

authController.updateProfile:
    - calls auth.UpdateProfileInfo(...)
    - calls auth.GenerateToken(...) + sends email if email changed

authController.updatePassword:
    - calls auth.UpdatePasswordAuthenticated(...)

authController.forgotPassword:
    - calls auth.GenerateToken(...)

authController.resetPassword:
    - calls auth.VerifyToken(...)
    - calls auth.UpdatePassword(...)

authController.emailVerifySend:
    - calls auth.GenerateToken(...)

authController.emailVerify:
    - calls auth.VerifyToken(...)
    - calls auth.MarkEmailVerified(...)

authController.passwordConfirmStatus:
    - reads session state directly (no action needed — just checks confirmed_at)

authController.passwordConfirm:
    - calls auth.ConfirmPassword(...)

authController.enable2FA:
    - calls auth.GenerateTotpSecret(...)

authController.confirm2FA:
    - calls auth.VerifyTotpCode(...)
    - calls auth.EnableTwoFactor(...)
    - calls auth.GenerateRecoveryCodes(...)

authController.disable2FA:
    - calls auth.DisableTwoFactor(...)

authController.challenge2FA:
    - calls auth.VerifyTotpCode(...)
    - calls auth.CreateSession(...)
    - calls auth.CreateRememberToken(...) if remember=true

authController.recoveryCodes:
    - reads recovery codes directly from DB (read-only, no action needed)

authController.regenerateRecoveryCodes:
    - calls auth.GenerateRecoveryCodes(...)
```