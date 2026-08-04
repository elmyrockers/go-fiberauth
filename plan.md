# Auth Routes Plan

All routes grouped under `/auth` prefix.

## Auth pages
| Method | Route | Purpose |
| ------ | ---------------------- | ------------------------------- |
| GET | /auth/login | Show login page |
| GET | /auth/register | Show registration page |
| GET | /auth/forgot-password | Show forgot password page |
| GET | /auth/reset-password | Show reset password page |
| GET | /auth/password/confirm | Show password confirmation page |
| GET | /auth/2fa/challenge | Show two-factor challenge page |

## Auth core
| Method | Route | Purpose |
|--------|-----------------|---------------------|
| POST | /auth/login | Log in (accepts optional `remember: bool`) |
| POST | /auth/logout | Log out |
| POST | /auth/register | Register new user |

## Profile & password management
| Method | Route | Purpose |
|--------|-----------------|-------------------------------------------------|
| PUT | /auth/profile | Update name/email (authenticated) |
| PUT | /auth/password | Change password (authenticated, requires current password) |

## Password reset
| Method | Route | Purpose |
|--------|--------------------------|-----------------------------------|
| POST | /auth/forgot-password | Request reset link (send email) |
| POST | /auth/reset-password | Submit new password + token |

## Email verification
| Method | Route | Purpose |
|--------|---------------------------|------------------------------|
| POST | /auth/email/verify-send | Send verification email |
| GET | /auth/email/verify | Verify email using signed URL token |

## Password confirmation (optional)
| Method | Route | Purpose |
|--------|--------------------------------|-----------------------------------|
| GET | /auth/password/confirm-status | Check if recently confirmed |
| POST | /auth/password/confirm | Re-enter password to confirm |

## Two-factor auth (optional)
| Method | Route | Purpose |
|--------|----------------------------|------------------------------------|
| POST | /auth/2fa/enable | Generate secret + QR code |
| POST | /auth/2fa/confirm | Confirm setup with a code |
| POST | /auth/2fa/disable | Turn off 2FA |
| POST | /auth/2fa/challenge | Verify code or recovery code during login (pending) |
| GET | /auth/2fa/recovery-codes | View recovery codes |
| POST | /auth/2fa/recovery-codes | Regenerate recovery codes |

## Passkeys / WebAuthn (optional)
| Method | Route | Purpose |
|--------|---------------------------------------|-------------------------------------------|
| GET | /auth/passkeys | List registered passkeys |
| POST | /auth/passkeys/register/options | Generate WebAuthn registration challenge |
| POST | /auth/passkeys/register | Store new passkey credential |
| PATCH | /auth/passkeys/{id} | Rename an existing passkey |
| POST | /auth/passkeys/login/options | Generate WebAuthn login challenge |
| POST | /auth/passkeys/login | Verify passkey signature & log in |
| DELETE | /auth/passkeys/{id} | Delete/revoke a passkey |


# Middleware Plan

Middleware acts strictly as gatekeepers (reading session state or cookies, no DB writes).

## Middleware List
| # | Middleware | Purpose |
|---|---------------------------------|----------------------------------------------------------------------------|
| 1 | GuestOnly | Blocks access if user is already logged in |
| 2 | AuthRequired | Blocks access if not logged in; falls back to remember-token cookie |
| 3 | RateLimiter | Throttles requests (reusable, configurable per route: max + window) |
| 4 | RequiresTwoFactorPending | Only allows access if login is mid-2FA-challenge (temp session state) |
| 5 | RequiresPasswordConfirmed | Blocks sensitive actions unless password was recently re-confirmed |
| 6 | RequiresVerifiedEmail | Blocks access if the authenticated user has not verified their email (applied externally to protected app routes like dashboards, or optionally wired here) |
| 7 | CSRF | Protects state-changing requests (Fiber's built-in `csrf.New()`) |

## Route → Middleware Mapping
| Route | Middleware |
| --------------------------------- | --------------------------------------------- |
| GET /auth/login | GuestOnly |
| GET /auth/register | GuestOnly |
| GET /auth/forgot-password | GuestOnly |
| GET /auth/reset-password | GuestOnly |
| GET /auth/password/confirm | AuthRequired |
| GET /auth/password/confirm-status | AuthRequired |
| GET /auth/email/verify | (none) |
| GET /auth/2fa/challenge | RequiresTwoFactorPending |
| GET /auth/2fa/recovery-codes | AuthRequired |
| GET /auth/passkeys | AuthRequired |
| POST /auth/login | GuestOnly, RateLimiter, CSRF |
| POST /auth/logout | AuthRequired, CSRF |
| POST /auth/register | GuestOnly, RateLimiter, CSRF |
| PUT /auth/profile | AuthRequired, CSRF |
| PUT /auth/password | AuthRequired, RequiresPasswordConfirmed, CSRF |
| POST /auth/forgot-password | GuestOnly, RateLimiter, CSRF |
| POST /auth/reset-password | GuestOnly, CSRF |
| POST /auth/email/verify-send | AuthRequired, RateLimiter, CSRF |
| POST /auth/password/confirm | AuthRequired, CSRF |
| POST /auth/2fa/enable | AuthRequired, RequiresPasswordConfirmed, CSRF |
| POST /auth/2fa/confirm | AuthRequired, CSRF |
| POST /auth/2fa/disable | AuthRequired, RequiresPasswordConfirmed, CSRF |
| POST /auth/2fa/challenge | RequiresTwoFactorPending, RateLimiter, CSRF |
| POST /auth/2fa/recovery-codes | AuthRequired, RequiresPasswordConfirmed, RateLimiter, CSRF |
| POST /auth/passkeys/register/options | AuthRequired, CSRF |
| POST /auth/passkeys/register | AuthRequired, CSRF |
| PATCH /auth/passkeys/{id} | AuthRequired, CSRF |
| POST /auth/passkeys/login/options | GuestOnly, RateLimiter, CSRF |
| POST /auth/passkeys/login | GuestOnly, RateLimiter, CSRF |
| DELETE /auth/passkeys/{id} | AuthRequired, RequiresPasswordConfirmed, CSRF |

> **Note on Ownership Checks:** Handlers for resource-scoped endpoints (such as `PATCH /auth/passkeys/{id}` and `DELETE /auth/passkeys/{id}`) must explicitly enforce ownership validation in application logic—ensuring the targeted passkey record belongs strictly to the currently authenticated session's user ID.


# Actions Plan

Actions are methods on the `Auth` struct — containing business logic (DB/session writes), called from inside controllers/handlers. 

## Actions List
| # | Action | Touches |
|---|-----------------------------------|-------------------------------------|
| 1 | VerifyCredentials | DB (read `users`) |
| 2 | CreateSession | Session (write `user_id`) |
| 3 | CreatePendingTwoFactorSession | Session (write `pending_2fa_user_id`) |
| 4 | CreateRememberToken | DB (write `users.remember_token`) + sets cookie |
| 5 | ValidateRememberTokenAndRestoreSession | DB (read/validate token) + Session (write `user_id`) |
| 6 | DestroySession | Session (+ clears remember token/cookie) |
| 7 | CreateUser | DB (write `users`) |
| 8 | GeneratePasswordResetToken | DB (`password_reset_tokens`) |
| 9 | VerifyPasswordResetToken | DB (`password_reset_tokens`) |
| 10| UpdatePassword (via reset token) | DB (write `users`) |
| 11| UpdatePasswordAuthenticated | DB (read+write `users`) |
| 12| UpdateProfileInfo | DB (write `users`) |
| 13| SendEmailVerificationNotification | Generates cryptographically signed URL & dispatches email |
| 14| VerifyEmailTokenAndMarkVerified | Validates signed URL token + DB (write `users.email_verified_at`) |
| 15| ConfirmPassword | DB (read `users`) + Session (write `confirmed_at`) |
| 16| GenerateTotpSecret | DB (write `users`, unconfirmed) |
| 17| VerifyTotpCode | DB (read secret) |
| 18| VerifyRecoveryCode | DB (read & consume hash from `users.two_factor_recovery_codes`) |
| 19| EnableTwoFactor | DB (write `users`) |
| 20| DisableTwoFactor | DB (write `users`) |
| 21| GenerateRecoveryCodes | DB (write `users`) |
| 22| GeneratePasskeyRegistrationOptions | WebAuthn challenge generation |
| 23| StorePasskey | DB (write `passkeys`) |
| 24| UpdatePasskeyName | DB (write `passkeys` with ownership scoping) |
| 25| GeneratePasskeyLoginOptions | WebAuthn login challenge generation |
| 26| VerifyPasskeyLogin | DB (read `passkeys`) + Session (write `user_id`) |
| 27| DeletePasskey | DB (delete `passkeys` with ownership scoping) |

## Handler → Action Mapping
| Handler | Actions called |
|---------------------------|------------------------------------------------------------------|
| Login | VerifyCredentials → CreateSession (or CreatePendingTwoFactor) → CreateRememberToken |
| Logout | DestroySession |
| Register | CreateUser → CreateSession |
| UpdateProfile | UpdateProfileInfo |
| UpdatePassword | UpdatePasswordAuthenticated |
| ForgotPassword | GeneratePasswordResetToken |
| ResetPassword | VerifyPasswordResetToken → UpdatePassword |
| EmailVerifySend | SendEmailVerificationNotification |
| EmailVerify | VerifyEmailTokenAndMarkVerified |
| PasswordConfirm | ConfirmPassword |
| Enable2FA | GenerateTotpSecret |
| Confirm2FA | VerifyTotpCode → EnableTwoFactor → GenerateRecoveryCodes |
| Disable2FA | DisableTwoFactor |
| Challenge2FA (login) | VerifyTotpCode (or VerifyRecoveryCode) → CreateSession → CreateRememberToken |
| RecoveryCodes (regen) | GenerateRecoveryCodes |
| PasskeysRegisterOptions | GeneratePasskeyRegistrationOptions |
| PasskeysRegister | StorePasskey |
| PasskeysRename | UpdatePasskeyName |
| PasskeysLoginOptions | GeneratePasskeyLoginOptions |
| PasskeysLogin | VerifyPasskeyLogin → CreateSession → CreateRememberToken |
| PasskeysDelete | DeletePasskey |

**Total: 27 distinct actions**


# Controller → Action Mapping

Controllers act as the HTTP layer: they extract request inputs, validate them, call the corresponding `Auth` actions, and return HTTP responses (JSON or redirects).

| Controller / Handler | Actions Called |
|-------------------------------------|------------------------------------------------------------------|
| **LoginController** (POST) | `VerifyCredentials` → `CreateSession` (or `CreatePendingTwoFactorSession`) → `CreateRememberToken` |
| **LogoutController** (POST) | `DestroySession` |
| **RegisterController** (POST) | `CreateUser` → `CreateSession` |
| **ProfileController** (PUT) | `UpdateProfileInfo` |
| **PasswordController** (PUT) | `UpdatePasswordAuthenticated` |
| **ForgotPasswordController** (POST) | `GeneratePasswordResetToken` |
| **ResetPasswordController** (POST) | `VerifyPasswordResetToken` → `UpdatePassword` |
| **EmailVerificationController** (POST) | `SendEmailVerificationNotification` |
| **EmailVerificationController** (GET) | `VerifyEmailTokenAndMarkVerified` |
| **ConfirmPasswordController** (POST) | `ConfirmPassword` |
| **TwoFactorEnableController** (POST) | `GenerateTotpSecret` |
| **TwoFactorConfirmController** (POST) | `VerifyTotpCode` → `EnableTwoFactor` → `GenerateRecoveryCodes` |
| **TwoFactorDisableController** (POST) | `DisableTwoFactor` |
| **TwoFactorChallengeController** (POST) | `VerifyTotpCode` (or `VerifyRecoveryCode`) → `CreateSession` → `CreateRememberToken` |
| **RecoveryCodesController** (GET) | *(None - reads existing recovery codes from DB)* |
| **RecoveryCodesController** (POST) | `GenerateRecoveryCodes` |
| **PasskeyRegisterOptionsController** (POST) | `GeneratePasskeyRegistrationOptions` |
| **PasskeyRegisterController** (POST) | `StorePasskey` |
| **PasskeyRenameController** (PATCH) | `UpdatePasskeyName` |
| **PasskeyLoginOptionsController** (POST) | `GeneratePasskeyLoginOptions` |
| **PasskeyLoginController** (POST) | `VerifyPasskeyLogin` → `CreateSession` → `CreateRememberToken` |
| **PasskeyDeleteController** (DELETE) | `DeletePasskey` |


# Database Table

```sql
-- users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    email_verified_at TIMESTAMPTZ,
    password VARCHAR(255) NOT NULL,
    remember_token VARCHAR(100),
    two_factor_secret TEXT,          -- encrypted
    two_factor_recovery_codes TEXT,  -- JSON array, hashed entries
    two_factor_confirmed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_email ON users(email);

-- password_reset_tokens
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX idx_password_reset_tokens_hash ON password_reset_tokens(token_hash);

-- passkeys
CREATE TABLE passkeys (
    id VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    public_key TEXT NOT NULL,
    credential_id TEXT NOT NULL UNIQUE,
    counter BIGINT NOT NULL,
    transports TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_passkeys_user_id ON passkeys(user_id);
```

## Usage

```go
import (
    "time"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/csrf"
    fiberauth "github.com/elmyrockers/go-fiberauth/session"
    "app/controller"
)

func main() {
    auth := fiberauth.New()
    authController := controller.NewAuth(auth)

    app := fiber.New()
    app.Use(csrf.New())

    // Auth pages
    app.Get("/auth/login", auth.GuestOnly, authController.loginPage)
    app.Get("/auth/register", auth.GuestOnly, authController.registerPage)
    app.Get("/auth/forgot-password", auth.GuestOnly, authController.forgotPasswordPage)
    app.Get("/auth/reset-password", auth.GuestOnly, authController.resetPasswordPage)
    app.Get("/auth/password/confirm", auth.AuthRequired, authController.passwordConfirmPage)
    app.Get("/auth/2fa/challenge", auth.RequiresTwoFactorPending, authController.challengePage)

    // Authentication
    app.Post("/auth/register", auth.GuestOnly, auth.RateLimiter(5, time.Minute), authController.register)
    app.Post("/auth/login", auth.GuestOnly, auth.RateLimiter(5, time.Minute), authController.login)
    app.Post("/auth/logout", auth.AuthRequired, authController.logout)

    // Profile & password
    app.Put("/auth/profile", auth.AuthRequired, authController.updateProfile)
    app.Put("/auth/password", auth.AuthRequired, auth.RequiresPasswordConfirmed, authController.updatePassword)

    // Password reset
    app.Post("/auth/forgot-password", auth.GuestOnly, auth.RateLimiter(3, time.Hour), authController.forgotPassword)
    app.Post("/auth/reset-password", auth.GuestOnly, authController.resetPassword)

    // Email verification
    app.Post("/auth/email/verify-send", auth.AuthRequired, auth.RateLimiter(3, time.Hour), authController.emailVerifySend)
    app.Get("/auth/email/verify", authController.emailVerify)

    // Password confirmation
    app.Get("/auth/password/confirm-status", auth.AuthRequired, authController.passwordConfirmStatus)
    app.Post("/auth/password/confirm", auth.AuthRequired, authController.passwordConfirm)

    // Two-factor authentication
    app.Post("/auth/2fa/enable", auth.AuthRequired, auth.RequiresPasswordConfirmed, authController.enable2FA)
    app.Post("/auth/2fa/confirm", auth.AuthRequired, authController.confirm2FA)
    app.Post("/auth/2fa/disable", auth.AuthRequired, auth.RequiresPasswordConfirmed, authController.disable2FA)
    app.Post("/auth/2fa/challenge", auth.RequiresTwoFactorPending, auth.RateLimiter(5, time.Minute), authController.challenge2FA)

    app.Get("/auth/2fa/recovery-codes", auth.AuthRequired, authController.recoveryCodes)
    app.Post("/auth/2fa/recovery-codes", auth.AuthRequired, auth.RequiresPasswordConfirmed, auth.RateLimiter(5, time.Minute), authController.regenerateRecoveryCodes)

    // Passkeys / WebAuthn
    app.Get("/auth/passkeys", auth.AuthRequired, authController.listPasskeys)
    app.Post("/auth/passkeys/register/options", auth.AuthRequired, authController.passkeyRegisterOptions)
    app.Post("/auth/passkeys/register", auth.AuthRequired, authController.passkeyRegister)
    app.Patch("/auth/passkeys/{id}", auth.AuthRequired, authController.passkeyRename)
    app.Post("/auth/passkeys/login/options", auth.GuestOnly, auth.RateLimiter(5, time.Minute), authController.passkeyLoginOptions)
    app.Post("/auth/passkeys/login", auth.GuestOnly, auth.RateLimiter(5, time.Minute), authController.passkeyLogin)
    app.Delete("/auth/passkeys/{id}", auth.AuthRequired, auth.RequiresPasswordConfirmed, authController.passkeyDelete)

    app.Listen(":3000")
}
```
