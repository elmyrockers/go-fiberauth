package main



import (
	"github.com/gofiber/fiber/v3"
	fiberauth "github.com/elmyrockers/go-xvelope/middleware/fiber"
	"github.com/davecgh/go-spew/spew"
)




func main() {
	app := fiber.New()
	app.Use( fiberauth.New() )

	app.Get( "/",func( c fiber.Ctx ) error {
		auth := fiberauth.FromContext(c)
		spew.Dump( auth )
		return c.SendString("OK")
	})

	app.Listen(":3000")
}