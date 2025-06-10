// models/task.go
package models

import "time"

type Task struct {
	ID          string    `json:"id,omitempty"`
	Titulo      string    `json:"titulo" validate:"required"`
	Descripcion string    `json:"descripcion"`
	FechaInicio time.Time `json:"fecha_inicio" validate:"required"`
	Deadline    time.Time `json:"deadline" validate:"required"`
	UsuarioID   string    `json:"usuario_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
