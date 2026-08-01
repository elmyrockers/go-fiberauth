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
	loginURL string
}

func New(config ...Config) *Auth {
	// Configure
		var cfg Config
		if len(config) > 1 { cfg = config[0] }

	// Set up Auth attributes, then return it
		loginURL := cfg.LoginURL
		if loginURL == "" { loginURL = "/auth/login" }

		session := fibersession.NewStore( cfg.SessionConfig )
		return &Auth{
			sessionStore: session,
			dbAdapter: cfg.DatabaseAdapter,
			mailAdapter: cfg.MailAdapter,
			loginURL: loginURL,
		}
}

// HELPER -------------------------------------------------------------------------
func wantsJSON(c fiber.Ctx) bool {
	if c.Get("X-Requested-With") == "XMLHttpRequest" {
		return true
	}
	return c.Accepts("html", "json") == "json"
}

// Redirects user to login url if unauthenticated
func (a *Auth) unauthenticated(c fiber.Ctx) error {
	if wantsJSON(c) {
		return fiber.ErrUnauthorized
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


// ---------------------------------------------------------------------------------