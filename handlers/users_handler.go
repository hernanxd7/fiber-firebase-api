package handlers

import (
    "context"
    "log" // <--- Importa el paquete log
    "github.com/gofiber/fiber/v2"
    "github.com/hernanxd7/fiber-firebase-api/config"
    "github.com/hernanxd7/fiber-firebase-api/models"
    "cloud.google.com/go/firestore"
    "google.golang.org/api/iterator"
)

// GetUser retrieves a user by ID from the Firestore database
func GetUser(c *fiber.Ctx) error {
    id := c.Params("id") // Obtiene el ID del parámetro de la URL
    ctx := context.Background()

    log.Printf("Attempting to get user with ID: %s", id) // <--- Log para ver el ID

    doc, err := config.FirestoreClient.Collection("users").Doc(id).Get(ctx)
    if err != nil {
        log.Printf("Error getting user %s: %v", id, err) // <--- Log para ver el error exacto de Firestore
        return c.Status(404).JSON(fiber.Map{"error": "User not found"})
    }

    var user models.User
    if err := doc.DataTo(&user); err != nil {
        log.Printf("Error mapping user data for ID %s: %v", id, err) // <--- Log para errores de mapeo
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    // Do not return the password
    userData := fiber.Map{
        "id":               doc.Ref.ID,
        "nombre":           user.Nombre,
        "apellidos":        user.Apellidos,
        "email":            user.Email,
        "fecha_nacimiento": user.FechaNacimiento,
    }

    log.Printf("Successfully retrieved user with ID: %s", id) // <--- Log si tiene éxito
    return c.JSON(userData)
}

// GetAllUsers retrieves all users from the Firestore database
// GetAllUsers retrieves all users from the Firestore database
func GetAllUsers(c *fiber.Ctx) error {
    ctx := context.Background()

    iter := config.FirestoreClient.Collection("users").Documents(ctx)
    defer iter.Stop()

    var users []fiber.Map

    for {
        doc, err := iter.Next()
        if err != nil {
            // Verifica si el error es específicamente iterator.Done
            if err == iterator.Done { // <--- Cambia la condición aquí
                break // Sale del bucle cuando no hay más documentos
            }
            // Si es otro tipo de error, retorna un 500
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }

        var user models.User
        if err := doc.DataTo(&user); err != nil {
            // Maneja errores al convertir datos del documento
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }

        // No devolver la contraseña
        userData := fiber.Map{
            "id":               doc.Ref.ID, // <--- Aquí obtienes el ID real del documento
            "nombre":           user.Nombre,
            "apellidos":        user.Apellidos,
            "email":            user.Email,
            "fecha_nacimiento": user.FechaNacimiento,
        }

        users = append(users, userData)
    }

    return c.JSON(users) // Retorna la lista de usuarios con el ID correcto
}

// UpdateUser updates a user's information in the Firestore database
func UpdateUser(c *fiber.Ctx) error {
    id := c.Params("id")
    ctx := context.Background()

    var updateData map[string]interface{}
    if err := c.BodyParser(&updateData); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid data"})
    }

    // Never update the password directly
    delete(updateData, "password")

    _, err := config.FirestoreClient.Collection("users").Doc(id).Set(ctx, updateData, firestore.MergeAll)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error updating user"})
    }

    return c.JSON(fiber.Map{"message": "User updated successfully"})
}

// DeleteUser removes a user from the Firestore database
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id") // Obtiene el ID del usuario del parámetro de la URL
	ctx := context.Background()

	log.Printf("Attempting to delete user with ID: %s", id) // Log para ver el ID

	// Intenta eliminar el documento del usuario por su ID en la colección "users"
	_, err := config.FirestoreClient.Collection("users").Doc(id).Delete(ctx)
	if err != nil {
		// Si el error es que el documento no existe, retorna 404
		if err != nil && err.Error() == "firestore: document not found" { // <--- Corrected check
			log.Printf("User not found for deletion with ID: %s", id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		// Si es otro tipo de error, retorna 500
		log.Printf("Error deleting user %s from Firestore: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting user"})
	}

	// Retorna respuesta de éxito
	log.Printf("User %s deleted successfully", id)
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}
