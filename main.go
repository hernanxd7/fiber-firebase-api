// main.go
package main

import (
	"github.com/hernanxd7/fiber-firebase-api/config"
	"github.com/hernanxd7/fiber-firebase-api/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	config.InitFirebase()
	routes.Setup(app)

	// Ruta raíz para probar en el navegador
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("¡API de Fiber + Firebase funcionando!")
	})

	app.Listen(":3000")
}
