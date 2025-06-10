package handlers

import (
    "context"
    "log"
    "time"
    // "errors" // <--- Remove this import

    "github.com/gofiber/fiber/v2"
    "github.com/go-playground/validator/v10"
    "github.com/hernanxd7/fiber-firebase-api/config"
    "github.com/hernanxd7/fiber-firebase-api/models"
    "cloud.google.com/go/firestore"
    "google.golang.org/api/iterator"
)

// GetAllTasks retrieves all tasks from the Firestore database
func GetAllTasks(c *fiber.Ctx) error {
    ctx := context.Background() // Obtiene un contexto

    // Obtiene un iterador para la colección "tasks"
    iter := config.FirestoreClient.Collection("tasks").Documents(ctx)
    defer iter.Stop() // Asegura que el iterador se detenga al finalizar la función

    var tasks []fiber.Map // Slice para almacenar los datos de las tareas

    // Itera sobre los documentos
    for {
        doc, err := iter.Next()
        if err != nil {
            // Verifica si el error es específicamente iterator.Done (fin de la colección)
            if err == iterator.Done {
                break // Sale del bucle
            }
            // Si es otro tipo de error, retorna un 500
            log.Printf("Error iterating tasks: %v", err)
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
        }

        var task models.Task // Estructura para mapear los datos del documento
        if err := doc.DataTo(&task); err != nil {
            // Maneja errores al convertir datos del documento a la estructura Task
            log.Printf("Error mapping task data for ID %s: %v", doc.Ref.ID, err)
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
        }

        // Crea un mapa con los datos de la tarea, incluyendo el ID del documento
        taskData := fiber.Map{
            "id":               doc.Ref.ID, // Usa el ID del documento de Firestore
            "titulo":           task.Titulo,
            "descripcion":      task.Descripcion,
            "fecha_inicio":     task.FechaInicio,
            "deadline":         task.Deadline,
            "usuario_id":       task.UsuarioID,
            "created_at":       task.CreatedAt,
            "updated_at":       task.UpdatedAt,
            // Añade otros campos si tu modelo Task los tiene
        }

        tasks = append(tasks, taskData) // Añade el mapa de datos al slice
    }

    // Retorna la lista de tareas en formato JSON
    log.Printf("Successfully retrieved %d tasks", len(tasks))
    return c.JSON(tasks)
}

// CreateTask creates a new task in the Firestore database
func CreateTask(c *fiber.Ctx) error {
    ctx := context.Background()
    task := new(models.Task) // Crea una nueva instancia de Task

    // 1. Parsear el cuerpo de la solicitud
    if err := c.BodyParser(task); err != nil {
        log.Printf("Error parsing task body: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task data"})
    }

    // 2. Validar los datos
    validate := validator.New()
    if err := validate.Struct(task); err != nil {
        log.Printf("Task validation error: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    // 3. Establecer marcas de tiempo
    task.CreatedAt = time.Now()
    task.UpdatedAt = time.Now()
    // El ID del documento se generará automáticamente por Firestore

    // 4. Guardar la tarea en Firestore
    ref, _, err := config.FirestoreClient.Collection("tasks").Add(ctx, task)
    if err != nil {
        log.Printf("Error adding task to Firestore: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating task"})
    }

    // 5. Devolver respuesta de éxito
    // Opcionalmente, puedes obtener el documento recién creado para devolverlo completo
    // doc, err := ref.Get(ctx)
    // if err != nil {
    //     log.Printf("Error getting created task document: %v", err)
    //     // Puedes decidir si esto es un error fatal o si devuelves solo el ID
    // }
    // var createdTask models.Task
    // if doc != nil {
    //     doc.DataTo(&createdTask)
    //     createdTask.ID = doc.Ref.ID // Asigna el ID del documento
    // }


    log.Printf("Task created successfully with ID: %s", ref.ID)
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Task created successfully",
        "id":      ref.ID, // Devuelve el ID generado por Firestore
        // "task":    createdTask, // Si decides obtener y devolver el objeto completo
    })
}

// GetTask retrieves a task by ID from the Firestore database
func GetTask(c *fiber.Ctx) error {
    id := c.Params("id") // Obtiene el ID del parámetro de la URL
    ctx := context.Background()

    // Obtiene la referencia al documento por su ID
    doc, err := config.FirestoreClient.Collection("tasks").Doc(id).Get(ctx)
    if err != nil {
        // Si el error es que el documento no existe, retorna 404
        if err != nil && err.Error() == "firestore: document not found" { // <--- Modified check
            log.Printf("Task not found with ID: %s", id)
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
        }
        // Si es otro tipo de error, retorna 500
        log.Printf("Error getting task %s from Firestore: %v", id, err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    var task models.Task // Estructura para mapear los datos del documento
    if err := doc.DataTo(&task); err != nil {
        // Maneja errores al convertir datos del documento a la estructura Task
        log.Printf("Error mapping task data for ID %s: %v", doc.Ref.ID, err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    // Asigna el ID del documento a la estructura Task
    task.ID = doc.Ref.ID

    // Retorna los datos de la tarea en formato JSON
    log.Printf("Successfully retrieved task with ID: %s", id)
    return c.JSON(task)
}

// UpdateTask updates a task's information in the Firestore database
func UpdateTask(c *fiber.Ctx) error {
    id := c.Params("id") // Obtiene el ID del parámetro de la URL
    ctx := context.Background()

    var updateData map[string]interface{} // Mapa para recibir los datos de actualización
    // 1. Parsear el cuerpo de la solicitud en un mapa
    if err := c.BodyParser(&updateData); err != nil {
        log.Printf("Error parsing update task body: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
    }

    // Opcional: Evitar que se actualicen ciertos campos si es necesario
    // Por ejemplo, si no quieres que se cambie el usuario_id o created_at
    delete(updateData, "usuario_id")
    delete(updateData, "created_at")
    // Puedes añadir otros campos aquí si no quieres que se actualicen

    // 2. Añadir o actualizar la marca de tiempo de actualización
    updateData["updated_at"] = time.Now()

    // 3. Realizar la actualización en Firestore
    // Usamos Set con firestore.MergeAll para actualizar solo los campos proporcionados
    _, err := config.FirestoreClient.Collection("tasks").Doc(id).Set(ctx, updateData, firestore.MergeAll)
    if err != nil {
        log.Printf("Error updating task %s in Firestore: %v", id, err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating task"})
    }

    // 4. Devolver respuesta de éxito
    log.Printf("Task %s updated successfully", id)
    return c.JSON(fiber.Map{"message": "Task updated successfully"})
}

// DeleteTask removes a task from the Firestore database
func DeleteTask(c *fiber.Ctx) error {
    id := c.Params("id") // Obtiene el ID del parámetro de la URL
    ctx := context.Background()

    // Intenta eliminar el documento por su ID
    _, err := config.FirestoreClient.Collection("tasks").Doc(id).Delete(ctx)
    if err != nil {
        // Si el error es que el documento no existe, retorna 404
        if err != nil && err.Error() == "firestore: document not found" { // <--- Modified check
            log.Printf("Task not found for deletion with ID: %s", id)
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
        }
        // Si es otro tipo de error, retorna 500
        log.Printf("Error deleting task %s from Firestore: %v", id, err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting task"})
    }

    // Retorna respuesta de éxito
    log.Printf("Task %s deleted successfully", id)
    return c.JSON(fiber.Map{"message": "Task deleted successfully"})
}
