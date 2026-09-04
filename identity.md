















SigninManager's Public Methods:
// --- Public Endpoints (Called by Application Controllers / Middleware) ---
// --- Password Sign-In Hot Path ---
	public virtual Task<SignInResult> PasswordSignInAsync(string userName, string password, bool isPersistent, bool lockoutOnFailure)
	public virtual Task<SignInResult> PasswordSignInAsync(TUser user, string password, bool isPersistent, bool lockoutOnFailure)
	public virtual Task<SignInResult> CheckPasswordSignInAsync(TUser user, string password, bool lockoutOnFailure)

// --- Direct Claims & Cookie / Token Issuance ---
	public virtual Task SignInAsync(TUser user, bool isPersistent, string? authenticationMethod = null)
	public virtual Task SignInAsync(TUser user, AuthenticationProperties authenticationProperties, string? authenticationMethod = null)
	public virtual Task SignInWithClaimsAsync(TUser user, bool isPersistent, IEnumerable<Claim> additionalClaims)
	public virtual Task SignInWithClaimsAsync(TUser user, AuthenticationProperties authenticationProperties, IEnumerable<Claim> additionalClaims)

// --- Session Lifecycle & Claims Inspection ---
	public virtual Task RefreshSignInAsync(TUser user)
	public virtual Task SignOutAsync()
	public virtual bool IsSignedIn(ClaimsPrincipal principal)
	public virtual Task<bool> CanSignInAsync(TUser user)
	public virtual Task<ClaimsPrincipal> CreateUserPrincipalAsync(TUser user)

// --- Security Stamp Validation ---
	public virtual Task<bool> ValidateSecurityStampAsync(TUser user, string? securityStamp)
	public virtual Task<TUser?> ValidateSecurityStampAsync(ClaimsPrincipal? principal)

// --- Protected Virtual Hooks (Internal Pipeline Extension Points) ---
	protected virtual Task<SignInResult?> PreSignInCheck(TUser user)
	protected virtual Task<bool> IsLockedOut(TUser user)
	protected virtual Task ResetLockout(TUser user)
	protected virtual Task<SignInResult> LockedOut(TUser user)



UserManager's Public Methods:
// --- Account CRUD & Querying ---
	public virtual Task<IdentityResult> CreateAsync(TUser user)
	public virtual Task<IdentityResult> CreateAsync(TUser user, string password)
	public virtual Task<IdentityResult> UpdateAsync(TUser user)
	public virtual Task<IdentityResult> DeleteAsync(TUser user)
	public virtual Task<TUser?> FindByIdAsync(string userId)
	public virtual Task<TUser?> FindByNameAsync(string userName)
	public virtual Task<TUser?> FindByEmailAsync(string email)

// --- Password Management ---
	public virtual Task<bool> CheckPasswordAsync(TUser user, string password)
	public virtual Task<bool> HasPasswordAsync(TUser user)
	public virtual Task<IdentityResult> AddPasswordAsync(TUser user, string password)
	public virtual Task<IdentityResult> ChangePasswordAsync(TUser user, string currentPassword, string newPassword)
	public virtual Task<IdentityResult> RemovePasswordAsync(TUser user)
	public virtual Task<string> GeneratePasswordResetTokenAsync(TUser user)
	public virtual Task<IdentityResult> ResetPasswordAsync(TUser user, string token, string newPassword)

// --- Email Confirmation ---
	public virtual Task<string> GenerateEmailConfirmationTokenAsync(TUser user)
	public virtual Task<IdentityResult> ConfirmEmailAsync(TUser user, string token)
	public virtual Task<bool> IsEmailConfirmedAsync(TUser user)

// --- Brute-Force Protection & Lockout ---
	public virtual Task<IdentityResult> AccessFailedAsync(TUser user)
	public virtual Task<IdentityResult> ResetAccessFailedCountAsync(TUser user)
	public virtual Task<int> GetAccessFailedCountAsync(TUser user)
	public virtual Task<bool> IsLockedOutAsync(TUser user)
	public virtual Task<bool> GetLockoutEnabledAsync(TUser user)
	public virtual Task<IdentityResult> SetLockoutEnabledAsync(TUser user, bool enabled)
	public virtual Task<DateTimeOffset?> GetLockoutEndDateAsync(TUser user)
	public virtual Task<IdentityResult> SetLockoutEndDateAsync(TUser user, DateTimeOffset? lockoutEnd)

// --- Security & Concurrency Stamps ---
	public virtual Task<string> GetSecurityStampAsync(TUser user)
	public virtual Task<IdentityResult> UpdateSecurityStampAsync(TUser user)
	public virtual Task<string> GenerateConcurrencyStampAsync(TUser user)
It has over 80 public methods in total.



UserStore:
	Interfaces
		Core (mandatory)
			IUserStore<TUser> — CRUD + id/name (~10 methods) ✅
				CreateAsync(TUser, CancellationToken)
				DeleteAsync(TUser, CancellationToken)
				UpdateAsync(TUser, CancellationToken)
				FindByIdAsync(string, CancellationToken)
				FindByNameAsync(string, CancellationToken)
				GetUserIdAsync(TUser, CancellationToken)
				GetUserNameAsync(TUser, CancellationToken)
				SetUserNameAsync(TUser, string, CancellationToken)
				GetNormalizedUserNameAsync(TUser, CancellationToken)
				SetNormalizedUserNameAsync(TUser, string, CancellationToken)
		Optional capability interfaces (all implemented by the EF Core UserStore)
			IUserPasswordStore<TUser> — password hash get/set/has (~3 methods) ✅
			IUserEmailStore<TUser> — email get/set/confirmed (~7 methods) ✅
			IUserLockoutStore<TUser> — lockout end date, enabled, failed count (~7 methods) ✅
			IUserClaimStore<TUser> — add/remove/replace/get claims (~5 methods) ✅
			IUserLoginStore<TUser> — external logins add/remove/find (~4 methods) ✅
			IUserRoleStore<TUser> — add/remove/get roles, users-in-role (~5 methods)
			IUserPhoneNumberStore<TUser> — phone number get/set/confirmed (~3 methods)
			IUserTwoFactorStore<TUser> — 2FA enabled get/set (~2 methods)
			IUserSecurityStampStore<TUser> — security stamp get/set (~2 methods)
			IUserAuthenticationTokenStore<TUser> — auth tokens (~3 methods)
			IUserAuthenticatorKeyStore<TUser> — authenticator key (~2 methods)
			IUserTwoFactorRecoveryCodeStore<TUser> — recovery codes (~3 methods)
			IQueryableUserStore<TUser> — just the Users IQueryable property
			IProtectedUserStore<TUser> — marker interface, no new members
			IUserPasskeyStore<TUser> — passkeys (newer, ~4 methods)


AUTHENTICATIONHANDLER:
	Cookie / Token-Based Handlers:
		- Cookie (encrypted http-only cookie payload )
		- Session (session-id with server store)
		- JWT
		- PASETO
		- Opaque (encrypted self-contain payload)
// --------------------------------------------------
	CookieAuthenticationHandler / BearerTokenHandler
	└─ extends SignInAuthenticationHandler<TOptions> (implements IAuthenticationSignInHandler)
	└─ extends SignOutAuthenticationHandler<TOptions> (implements IAuthenticationSignOutHandler)
	└─ extends AuthenticationHandler<TOptions>
	└─ implements IAuthenticationHandler

// Both handlers expose these identical public signatures:
	IAuthenticationHandler
	    Task InitializeAsync(AuthenticationScheme scheme, HttpContext context);
	    Task<AuthenticateResult> AuthenticateAsync();
	    Task ChallengeAsync(AuthenticationProperties? properties);
	    Task ForbidAsync(AuthenticationProperties? properties);
	IAuthenticationSignInHandler
		Task SignInAsync(ClaimsPrincipal user, AuthenticationProperties? properties);
	IAuthenticationSignOutHandler
		Task SignOutAsync(AuthenticationProperties? properties);



```go
opaqueHandler := opaque.NewHandler(opaque.Options{
    TokenStore: redisTokenStore, // Interface to validate & fetch claims
})

pasetoHandler := paseto.NewHandler(paseto.Options{
    PublicKey: myEd25519PublicKey,
    Issuer:    "https://auth.example.com",
    Audience:  "https://api.example.com",
})

//Register as distinct schemes in xvelope
authConfig := &auth.Config{
    Schemes: []*auth.Scheme{
        {
            Name:    "OpaqueBearer",
            Handler: opaqueHandler,
        },
        {
            Name:    "PasetoBearer",
            Handler: pasetoHandler,
        },
    },
}
```





