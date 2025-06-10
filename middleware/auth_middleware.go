package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hernanxd7/fiber-firebase-api/utils"
)

// Protected es un middleware que verifica si el token JWT es válido
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obtener el token del encabezado Authorization
		authHeader := c.Get("Authorization")
		
		// Verificar si el encabezado existe y tiene el formato correcto
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "Se requiere token de autenticación"})
		}
		
		// Extraer el token (quitar el prefijo "Bearer ")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Validar el token
		userID, err := utils.ValidateJWT(tokenString)
		if err != nil {
			// Agregar log para depuración
			return c.Status(401).JSON(fiber.Map{"error": "Token inválido o expirado"})
		}
		
		// Almacenar el ID del usuario en el contexto para usarlo en los handlers
		c.Locals("userID", userID)
		
		// Continuar con la siguiente función en la cadena
		return c.Next()
	}
}