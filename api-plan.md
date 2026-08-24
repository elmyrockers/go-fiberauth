# Authentication Middleware for Golang Fiber v3
## A. Setup
- `New( config ) *FiberAuth`
    - ***Usage:***
        - `app := fiber.New()`
        - `app.Use( fiberauth.New(config) )`
- `FromContext( context fiber.Ctx ) *FiberAuth`

## B. Middlewares

1. **GuestOnly** — Blocks access if user is already logged in
2. **AuthRequired** — Blocks access if not logged in; falls back to remember-token cookie
3. **RateLimiter** — Throttles requests (reusable, configurable per route: max + window)
4. **RequiresTwoFactorPending** — Only allows access if login is mid-2FA-challenge (temp session state)
5. **RequiresPasswordConfirmed** — Blocks sensitive actions unless password was recently re-confirmed
6. **RequiresVerifiedEmail** — Blocks access if the authenticated user has not verified their email (applied externally to protected app routes like dashboards, or optionally wired here)
7. **CSRF** — Protects state-changing requests (Fiber's built-in `csrf.New()`)

## C. Actions

### Password-Based Authentication

#### 1. User Lookup
- `FindUserByID(id string) (User, error)`
- `FindUserByName(username string) (User, error)`
- `FindUserByEmail(email string) (User, error)`

#### 2. Registration
- `NewUser() User`
- `CreateUser(user User, password string) error`
- `GenerateEmailConfirmationToken(user User) (token string, error)`
- `SendConfirmationEmail(user User, token string) error` - Add email delivery support (non-Identity)
- `ConfirmEmail(user User, token string) error`

#### 3. Login & Session
- `CheckPassword(user User, password string) (isValid bool, error)`
- `SignIn(scheme string, user User, isPersistent bool) (Session, error)` — low level, it will just create login state only
  - `IssueAuthTicket(scheme string, user User, isPersistent bool) (string, error)`
        - Go equivalent of `HttpContext.SignInAsync(string scheme, ClaimsPrincipal principal, AuthenticationProperties properties)` - Identity method
  - `buildClaims(user)`, `issueCookieTicket(claims, isPersistent)`, `issueBearerToken(claims, isPersistent)`
        - OPAQUE TOKENS AND COOKIE
        - DURING LOGIN: Detect query — `?useCookies=true`
  - `GuestOnly` and `AuthRequired` middlewares — auto-detect `bearer or cookie header` on every request
- `PasswordSignIn(user User, password string, isPersistence bool, lockoutOnFailure bool) (SignInResult, error)` — signin with User instance, will call `SignIn()` underneath
- `PasswordSignInWithEmail(email, password string, isPersistence bool, lockoutOnFailure bool) (SignInResult, error)` — signin with email, will call `PasswordSignIn()` underneath
- `SignOut() error`

#### 4. Password Reset (Unauthenticated)
- `GeneratePasswordResetToken(user User) (resetToken string, error)`
- `SendPasswordResetEmail(user User, token string) error` - Add email delivery support (non-Identity)
- `ResetPassword(user User, token, newPassword string) error`

#### 5. Profile & Password Updates (Authenticated)
- `UpdateUser(user User) error`
- `DeleteUser(user User) error`
- `GenerateChangeEmailToken(user User, newEmail string) (emailChangeToken string, error)`
- `ChangeEmail(user User, newEmail, token string) error`
- `ChangePassword(user User, currentPassword, newPassword string) error`

#### 6. Password Confirmation (Authenticated)
- `CheckPassword(user User, password string) (isValid bool, error)`