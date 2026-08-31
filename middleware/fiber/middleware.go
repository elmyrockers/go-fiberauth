package fiber

import (
	"github.com/gofiber/fiber/v3"
	"github.com/elmyrockers/go-xvelope"
)

// Private type prevents external key collisions
type contextKey struct{}

// Unexported key instance (zero memory allocation)
var authKey = contextKey{}

type Auth struct {
	context xvelope.HttpContext
}

// Config is reserved for future options.
type Config struct{}

// New() creates the auth middleware.
func New(config ...Config) fiber.Handler {
	// Create Auth instance
		auth := &Auth{
			context: &xvelope.FastHttpContext{},
		}

	// Set httpcontext then store auth instance
		return func(c fiber.Ctx) error {
			auth.context.SetHttpContext( c.Context() )
			c.Locals(authKey, auth)

			return c.Next()
		}
}

// FromContext(c) retrieves the *Auth instance from Fiber Ctx
func FromContext(c fiber.Ctx) *Auth {
	auth, _ := c.Locals(authKey).(*Auth)
	return auth
}