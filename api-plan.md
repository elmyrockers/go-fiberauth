# Authentication Middleware for Golang Fiber v3

## A. Middlewares

1. **GuestOnly** — Blocks access if user is already logged in
2. **AuthRequired** — Blocks access if not logged in; falls back to remember-token cookie
3. **RateLimiter** — Throttles requests (reusable, configurable per route: max + window)
4. **RequiresTwoFactorPending** — Only allows access if login is mid-2FA-challenge (temp session state)
5. **RequiresPasswordConfirmed** — Blocks sensitive actions unless password was recently re-confirmed
6. **RequiresVerifiedEmail** — Blocks access if the authenticated user has not verified their email (applied externally to protected app routes like dashboards, or optionally wired here)
7. **CSRF** — Protects state-changing requests (Fiber's built-in `csrf.New()`)

## B. Actions

### Password-Based Authentication

#### 1. User Lookup
- `FindUserByID(id string) (User, error)`
- `FindUserByName(username string) (User, error)`
- `FindUserByEmail(email string) (User, error)`

#### 2. Registration
- `NewUser() User`
- `CreateUser(user User, password string) error`
- `GenerateEmailConfirmationToken(user User) (token string, error)`
- `SendConfirmationEmail(user User, token string) error`
- `ConfirmEmail(user User, token string) error`

#### 3. Login & Session
- `CheckPassword(user User, password string) (isValid bool, error)`
- `SignIn(user User, isPersistent bool) (Session, error)` — low level, it will just create login state only
  - `HttpContext.SignInAsync(string scheme, ClaimsPrincipal principal, AuthenticationProperties properties)`
    - OPAQUE TOKENS AND COOKIE
    - DURING LOGIN: Detect query — `?useCookies=true`
    - EVERY REQUEST: Auto-detect bearer/cookie header
- `PasswordSignIn(user User, password string, rememberMe bool) (SignInResult, error)` — signin with User instance, will call `SignIn()` underneath
- `PasswordSignInWithEmail(email, password string, rememberMe bool) (SignInResult, error)` — signin with email, will call `PasswordSignIn()` underneath
- `SignOut() error`

#### 4. Password Reset (Unauthenticated)
- `GeneratePasswordResetToken(user User) (resetToken string, error)`
- `SendPasswordResetEmail(user User, token string) error`
- `ResetPassword(user User, token, newPassword string) error`

#### 5. Profile & Password Updates (Authenticated)
- `UpdateUser(user User) error`
- `DeleteUser(user User) error`
- `GenerateChangeEmailToken(user User, newEmail string) (emailChangeToken string, error)`
- `ChangeEmail(user User, newEmail, token string) error`
- `ChangePassword(user User, currentPassword, newPassword string) error`

#### 6. Password Confirmation (Authenticated)
- `CheckPassword(user User, password string) (isValid bool, error)`