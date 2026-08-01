package session

import (
	// "github.com/gofiber/fiber/v3"
	fibersession "github.com/gofiber/fiber/v3/middleware/session"

	// "github.com/davecgh/go-spew/spew"
)


type Auth struct {
	sessionStore *fibersession.Store
	dbAdapter DatabaseAdapter
	mailAdapter MailAdapter
}

func New(config ...Config) *Auth {
	// Configure
		var cfg Config
		if len(config) > 1 { cfg = config[0] }

	// Set up Auth attributes, then return it
		session := fibersession.NewStore( cfg.SessionConfig )
		return &Auth{
			sessionStore: session,
			dbAdapter: cfg.DatabaseAdapter,
			mailAdapter: cfg.MailAdapter,
		}
}