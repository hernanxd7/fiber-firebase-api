// models/user.go
package models

import "time"

type User struct {
	ID               string    `json:"id,omitempty"`
	Nombre           string    `json:"nombre" validate:"required"`
	Apellidos        string    `json:"apellidos" validate:"required"`
	Email            string    `json:"email" validate:"required,email"`
	Password         string    `json:"password,omitempty" validate:"required,min=6"`
	FechaNacimiento  time.Time `json:"fecha_nacimiento" validate:"required"`
	PreguntaSecreta  string    `json:"pregunta_secreta" validate:"required"`
	RespuestaSecreta string    `json:"respuesta_secreta" validate:"required"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}
