package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hernanxd7/fiber-firebase-api/config"
	"github.com/hernanxd7/fiber-firebase-api/routes"
	"github.com/stretchr/testify/assert"
)

// Configuración del servidor para pruebas
func setupTestApp() *fiber.App {
	app := fiber.New()
	config.InitFirebase()
	routes.Setup(app)
	return app
}

// Prueba de la ruta raíz
func TestRootRoute(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "API de Fiber + Firebase funcionando")
}

// Prueba de registro de usuario
func TestRegisterUser(t *testing.T) {
	app := setupTestApp()

	// Datos de prueba para registro
	userData := `{
		"nombre": "Usuario",
		"apellidos": "De Prueba",
		"email": "test@example.com",
		"password": "password123",
		"fecha_nacimiento": "1990-01-01T00:00:00Z",
		"pregunta_secreta": "¿Nombre de tu mascota?",
		"respuesta_secreta": "Rex"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/register",
		strings.NewReader(userData),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	
	// Verificamos que el registro sea exitoso (código 201)
	// Nota: Esta prueba fallará si el usuario ya existe
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

// Prueba de login
func TestLogin(t *testing.T) {
	app := setupTestApp()

	// Datos de prueba para login
	loginData := `{
		"email": "test@example.com",
		"password": "password123"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		strings.NewReader(loginData),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	
	// Verificamos que el login sea exitoso
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verificamos que se devuelva un token
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "token")
}