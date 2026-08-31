package fiber

import (
	"github.com/gofiber/fiber/v3"
	"github.com/elmyrockers/go-xvelope"
)

// Private type prevents external key collisions
type contextKey struct{}

// Unexported key instance (zero memory allocation)
var authKey = contextKey{}

// New() creates the auth middleware.
func New(config ...xvelope.Config) fiber.Handler {
	// Resolve config slice
		var cfg xvelope.Config
		if len(config) > 0 {
			cfg = config[0]
		}

	// Create Auth instance
		auth := xvelope.New( cfg )

	// Set httpcontext then store auth instance
		return func(c fiber.Ctx) error {
			httpCtx := &xvelope.FastHttpContext{}
			httpCtx.SetContext( c.RequestCtx() )
			auth.SetHttpContext( httpCtx )

			c.Locals(authKey, auth)
			return c.Next()
		}
}

// FromContext(c) retrieves the *Auth instance from Fiber Ctx
func FromContext(c fiber.Ctx) *xvelope.Auth {
	auth, _ := c.Locals(authKey).(*xvelope.Auth)
	return auth
}