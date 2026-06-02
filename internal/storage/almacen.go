package storage

import (
	"errors"

	"uleam-transporte/internal/dominio"
)

var (
	ErrNotFound = errors.New("recurso no encontrado")
	ErrConflict = errors.New("el recurso ya existe")
)

type Almacen interface {
	dominio.CarritoRepository
	dominio.HorarioRepository
	dominio.CarritoHorarioRepository
}
