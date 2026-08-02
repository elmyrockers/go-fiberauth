package controller

import (
	"github.com/gofiber/fiber/v3"
)

type AuthController struct {}

func NewAuth() *AuthController {
	return &AuthController{}
}

func (ac *AuthController) Register(c fiber.Ctx) error {
	return c.Render("register", fiber.Map{
				"Title": "Register Form",
			})
}

func (ac *AuthController) Login(c fiber.Ctx) error {
	return c.Render("login", fiber.Map{
				"Title": "Login Form",
			})
}