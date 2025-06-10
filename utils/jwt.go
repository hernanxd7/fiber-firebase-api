// utils/jwt.go
package utils

import (
	"errors" // <--- Añade esta línea para importar el paquete errors
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Clave secreta para firmar los tokens JWT
var jwtSecret = []byte("tu_clave_secreta_aqui") // Cambia esto por una clave segura en producción

// GenerateJWT genera un token JWT para el ID de usuario proporcionado
func GenerateJWT(userID string) (string, error) {
	// Crear un nuevo token con el algoritmo de firma HS256
	token := jwt.New(jwt.SigningMethodHS256)
	
	// Configurar los claims (reclamaciones) del token
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix() // Token válido por 10 minutos
	claims["iat"] = time.Now().Unix() // Tiempo de emisión
	
	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidateJWT valida un token JWT y devuelve el ID de usuario si es válido
// Asegúrate de que la función ValidateJWT esté correctamente implementada
func ValidateJWT(tokenString string) (string, error) {
	// Parsear el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma sea el esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	
	if err != nil {
		return "", err
	}
	
	// Verificar que el token sea válido
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extraer el ID de usuario
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("user_id no encontrado en el token")
		}
		return userID, nil
	}
	
	return "", errors.New("token inválido")
}
