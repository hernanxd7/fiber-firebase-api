// handlers/auth_handler.go
package handlers

import (
	"context"
	"github.com/hernanxd7/fiber-firebase-api/config"
    "github.com/hernanxd7/fiber-firebase-api/models"
    "github.com/hernanxd7/fiber-firebase-api/utils"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Hashear contraseña
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al encriptar contraseña"})
	}
	user.Password = hash

	doc, _, err := config.FirestoreClient.Collection("users").Add(context.Background(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al registrar usuario"})
	}

	user.ID = doc.ID
	user.Password = "" // nunca retornar contraseña
	return c.Status(201).JSON(user)
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	iter := config.FirestoreClient.Collection("users").Where("Email", "==", input.Email).Limit(1).Documents(context.Background())
	doc, err := iter.Next()
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Email o contraseña incorrectos"})
	}

	var user models.User
	doc.DataTo(&user)

	if !utils.CheckPassword(user.Password, input.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Email o contraseña incorrectos"})
	}

	token, err := utils.GenerateJWT(doc.Ref.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error generando token"})
	}

	// Devolver más información para depuración
	return c.JSON(fiber.Map{
		"token": token,
		"user_id": doc.Ref.ID,
		"expires_in": "10 minutos",
	})
}
