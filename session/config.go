package session

import (
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage"
)

type Config struct {
	SessionStorage storage.Storage
}