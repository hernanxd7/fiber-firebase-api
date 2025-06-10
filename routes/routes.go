// routes/routes.go
package routes

import (
    "github.com/hernanxd7/fiber-firebase-api/handlers"
    "github.com/hernanxd7/fiber-firebase-api/middleware"
    "github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
    api := app.Group("/api")

    // Rutas públicas de autenticación
    api.Post("/register", handlers.Register)
    api.Post("/login", handlers.Login)

    // Rutas protegidas
    user := api.Group("/users", middleware.Protected())
    task := api.Group("/tasks", middleware.Protected())

    // Usar los handlers implementados
    user.Get("/", handlers.GetAllUsers)
    user.Get("/:id", handlers.GetUser)
    user.Put("/:id", handlers.UpdateUser)
    user.Delete("/:id", handlers.DeleteUser)

    task.Get("/", handlers.GetAllTasks)
    task.Post("/", handlers.CreateTask)
    task.Get("/:id", handlers.GetTask)
    task.Put("/:id", handlers.UpdateTask)
    task.Delete("/:id", handlers.DeleteTask)
}
