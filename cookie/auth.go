package cookie

import (
	"fmt"
	"time"
	"crypto/rsa"
	"crypto/rand"
	"errors"

	"github.com/gofiber/fiber/v3"
	fibersession "github.com/gofiber/fiber/v3/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"

	// "github.com/davecgh/go-spew/spew"
)

type Auth struct {
	sharedKey []byte
	privateKey *rsa.PrivateKey
	publicKey *rsa.PublicKey

	sessionStore *fibersession.Store
	dbAdapter DatabaseAdapter
	mailAdapter MailAdapter
	mailTemplates map[string]MailTemplate
	loginURL string
	appURL string
}

func New(config ...Config) *Auth {
	// Configure
		var cfg Config
		if len(config) > 0 { cfg = config[0] }

	// Set shared-key and rsa-keys for JWT
		sharedKey := cfg.SharedKey
		privateKeyPath := cfg.PrivateKeyPath
		publicKeyPath := cfg.PublicKeyPath
		var err error

		// if both are empty, default to hmac auto-generated
			if sharedKey == nil && privateKeyPath == "" {
				sharedKey, err = generateHMACKey()
				if err != nil {
					return nil
				}
			}

		// if rsa-keys pem are provided
			var privateKey *rsa.PrivateKey
			var publicKey *rsa.PublicKey
			if privateKeyPath!="" && publicKeyPath!="" {
				privateKey, err = cfg.LoadPrivateKey()
				if err != nil {
					return nil
				}
				publicKey, err = cfg.LoadPublicKey()
				if err != nil {
					return nil
				}
			}
	// Set up Mail config
		mailAdapter := cfg.MailAdapter
		mailConfig := cfg.MailConfig
		if mailConfig.From.Email != "" {
			mailAdapter.SetFrom( mailConfig.From )
		}
		if mailConfig.ReplyTo != nil && mailConfig.ReplyTo.Email != "" {
			mailAdapter.SetReplyTo( *mailConfig.ReplyTo )
		}

	// Set up Auth attributes, then return it
		loginURL := cfg.LoginURL
		if loginURL == "" { loginURL = "/auth/login" }

		session := fibersession.NewStore( cfg.SessionConfig )
		return &Auth{
			sharedKey: sharedKey,
			privateKey: privateKey,
			publicKey: publicKey,

			sessionStore: session,
			dbAdapter: cfg.DatabaseAdapter,
			mailAdapter: mailAdapter,
			mailTemplates: mailConfig.Templates,
			loginURL: loginURL,
		}
}

// HELPER -------------------------------------------------------------------------

// GenerateHMACKey generates a 32-byte (256-bit) cryptographically secure key
func generateHMACKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return key, nil
}

func wantsJSON(c fiber.Ctx) bool {
	if c.Get("X-Requested-With") == "XMLHttpRequest" {
		return true
	}
	return c.Accepts("html", "json") == "json"
}

// Redirects user to login url if unauthenticated
func (a *Auth) unauthenticated(c fiber.Ctx) error {
	if wantsJSON(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated.",
		})
	}
	return c.Redirect().To(a.loginURL)
}

// MIDDLEWARE ----------------------------------------------------------------------
func (a *Auth) AuthRequired(c fiber.Ctx) error {
	// Make sure user_id exists
		sess, err := a.sessionStore.Get(c)
		if err != nil { return err }

		userID := sess.Get("user_id")
		if userID == nil {
			return a.unauthenticated(c)
		}

	return c.Next()
}


// ACTIONS --------------------------------------------------------------------------
func (a *Auth) CreateUser(name, email, password string) (int64, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := a.dbAdapter.NewUser()
	user.SetName(name)
	user.SetEmail(email)
	user.SetPassword(string(hashed))

	userID, err := a.dbAdapter.CreateUser(user)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (a *Auth) buildEmailVerificationLink(userID int64, email, verifyURL string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"purpose": EmailVerification,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	var signingMethod jwt.SigningMethod
	var signingKey any
	signingMethod = jwt.SigningMethodHS256
	signingKey = a.sharedKey
	if a.privateKey != nil {
		signingMethod = jwt.SigningMethodRS256
		signingKey = a.privateKey
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	signed, err := token.SignedString(signingKey)
	if err != nil { return "", err }

	return fmt.Sprintf("%s?token=%s", verifyURL, signed), nil
}

func (a *Auth) SendEmailVerificationNotification(userID int64, name, email, verifyURL string) error {
	// Build verification link
		link, err := a.buildEmailVerificationLink( userID, email, verifyURL )
		if err != nil { return err }

	// Send notification email
		mail := a.mailAdapter
		mail.SetTo(MailAddress{Name: name, Email: email})
		mail.SetTemplate(MailTemplate{
			Subject: "Verify your email address",
			Body:    fmt.Sprintf("Click the link to verify your email: %s", link),
			AltBody: fmt.Sprintf("Verify your email: %s", link),
		})
		return mail.Send()
}

func (a *Auth) VerifyEmail(tokenString string) (int64, error) {
	// Parse JWT Token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if a.publicKey != nil {
				if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return a.publicKey, nil
			}
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return a.sharedKey, nil
		})

	// Make sure token is valid
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrTokenExpired
		}
		if err != nil || !token.Valid {
			return 0, ErrInvalidToken
		}

	// Get claims
		claims, ok := token.Claims.(jwt.MapClaims)

		// Check purpose
			if !ok || claims["purpose"] != EmailVerification {
				return 0, ErrInvalidToken
			}

		// Get User ID
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				return 0, ErrInvalidToken
			}
			userID := int64(userIDFloat)

		// Get email from claims
			tokenEmail, ok := claims["email"].(string)
			if !ok || tokenEmail == "" {
				return 0, ErrInvalidToken
			}

		// Confirm the token's email still matches the user's current email
			user, err := a.dbAdapter.FindUserByID(userID)
			if err != nil {
				return 0, err
			}
			if user.GetEmail() != tokenEmail {
				return 0, ErrInvalidToken // email changed since link was issued — reject
			}

	return userID, nil
}

func (a *Auth) CreateSession(c fiber.Ctx, userID string) error {
	// Get session store and regenerate session id
		sess, err := a.sessionStore.Get(c)
		if err != nil { return err }

		if err := sess.Regenerate(); err != nil { return err }

	// Store user_id into session
		sess.Set("user_id", userID)
		sess.Delete("pending_2fa_user_id")

	return sess.Save()
}

