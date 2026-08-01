package session

import (
	"github.com/gofiber/fiber/v3"
	fibersession "github.com/gofiber/fiber/v3/middleware/session"

	// "github.com/davecgh/go-spew/spew"
)


type Auth struct {
	sessionStore *fibersession.Store
	dbAdapter DatabaseAdapter
	mailAdapter MailAdapter
	loginUrl string
}

func New(config ...Config) *Auth {
	// Configure
		var cfg Config
		if len(config) > 1 { cfg = config[0] }

	// Set up Auth attributes, then return it
		loginURL := cfg.LoginURL
		if loginURL == nil { loginURL = "/auth/login" }

		session := fibersession.NewStore( cfg.SessionConfig )
		return &Auth{
			sessionStore: session,
			dbAdapter: cfg.DatabaseAdapter,
			mailAdapter: cfg.MailAdapter,
			loginURL: loginURL,
		}
}

// Redirects user to login url if unauthenticated
func (a *Auth) unauthenticated(c fiber.Ctx) error {
	if c.Accepts("json") == "json" {
		return fiber.ErrUnauthorized
	}

	return c.Redirect(a.LoginURL)
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


// ---------------------------------------------------------------------------------