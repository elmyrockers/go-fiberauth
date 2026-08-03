package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

type AuthController struct {}

func NewAuth() *AuthController {
	return &AuthController{}
}

// HTTP METHOD: GET
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

// HTTP METHOD: POST, PUT
func (ac *AuthController) Register(c fiber.Ctx) error {
	// return c.Redirect().To( "/auth/register" )
	return c.SendString( "Process Register" )
}

func (ac *AuthController) Login(c fiber.Ctx) error {
	// return c.Redirect().To( "/auth/login" )
	return c.SendString( "Process Login" )
}

func (ac *AuthController) ForgotPassword(c fiber.Ctx) error {
	// return c.Redirect().To( "/auth/login" )
	return c.SendString( "Process ForgotPassword" )
}