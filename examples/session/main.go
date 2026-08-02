package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/template/jet/v3"
	fiberauth "github.com/elmyrockers/go-fiberauth/session"

	"github.com/elmyrockers/go-fiberauth/examples/session/controller"

	// "github.com/davecgh/go-spew/spew"
)

func Hello(c fiber.Ctx) error {
	return c.SendString("Hello, Welcome 👋!")
}

func main(){
	auth := fiberauth.New()
	_ = auth

	authController := controller.NewAuth()

	engine := jet.New( "views", ".jet" )
	app := fiber.New(fiber.Config{ Views: engine })
	app.Use(csrf.New())


	app.Get("/auth/login", authController.LoginPage )
	app.Get("/auth/register", authController.RegisterPage )
	app.Get("/auth/forgot-password", authController.ForgotPasswordPage )
	app.Get("/auth/reset-password", authController.ResetPasswordPage )



	// app.Post("/auth/login", authController.Login)
	// app.Get("/member/index", auth.AuthRequired, Hello)
	// app.Post("/auth/register", auth.GuestOnly, auth.RateLimiter(5, time.Minute), authController.register )



	app.Listen(":3000")
}