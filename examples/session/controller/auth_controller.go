package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

type AuthController struct {}

func NewAuth() *AuthController {
	return &AuthController{}
}

func (ac *AuthController) RegisterPage(c fiber.Ctx) error {
	return c.Render("register", fiber.Map{
				"Title": "Register Form",
				"csrf": csrf.TokenFromContext(c),
			})
}

func (ac *AuthController) LoginPage(c fiber.Ctx) error {
	return c.Render("login", fiber.Map{
				"Title": "Login Form",
				"csrf": csrf.TokenFromContext(c),
			})
}

func (ac *AuthController) ForgotPasswordPage(c fiber.Ctx) error {
	return c.Render("forgot-password", fiber.Map{
				"Title": "Forgot Password Form",
				"csrf": csrf.TokenFromContext(c),
			})
}

func (ac *AuthController) ResetPasswordPage(c fiber.Ctx) error {
	return c.Render("reset-password", fiber.Map{
				"Title": "Reset Password Form",
				"csrf": csrf.TokenFromContext(c),
				"token": "",
			})
}