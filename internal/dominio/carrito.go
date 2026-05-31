
package dominio

import "errors"

type EstadoCarrito string

const (
	EstadoDisponible    EstadoCarrito = "Disponible"
	EstadoEnViaje       EstadoCarrito = "En Viaje"
	EstadoMantenimiento EstadoCarrito = "Mantenimiento"
)

func (e EstadoCarrito) Valid() bool {
	switch e {
	case EstadoDisponible, EstadoEnViaje, EstadoMantenimiento:
		return true
	}
	return false
}

type Carrito struct {
	Numero             int           `json:"numero_carrito"`
	Estado             EstadoCarrito `json:"estado_carrito"`
	CapacidadPasajeros int           `json:"capacidad_pasajeros"`
}

func (c *Carrito) Validate() error {
	if c.Numero <= 0 {
		return errors.New("numero_carrito debe ser un entero positivo")
	}
	if !c.Estado.Valid() {
		return errors.New("estado_carrito debe ser 'Disponible', 'En Viaje' o 'Mantenimiento'")
	}
	if c.CapacidadPasajeros <= 0 {
		return errors.New("capacidad_pasajeros debe ser un entero positivo")
	}
	return nil
}

type CarritoRepository interface {
	GetAllCarritos() ([]Carrito, error)
	GetCarritoByNumero(numero int) (Carrito, error)
	CreateCarrito(c Carrito) error
	UpdateCarrito(c Carrito) error
	DeleteCarrito(numero int) error
	GetHorariosByCarrito(numero int) ([]CarritoHorario, error)
}
