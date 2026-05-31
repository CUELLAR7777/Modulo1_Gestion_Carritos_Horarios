package dominio

import "errors"

type Turno string

const (
	TurnoMatutino   Turno = "Matutina"
	TurnoVespertino Turno = "Vespertina"
)

func (t Turno) Valid() bool {
	switch t {
	case TurnoMatutino, TurnoVespertino:
		return true
	}
	return false
}

type Horario struct {
	ID        int    `json:"id_horario"`
	Turno     Turno  `json:"turno"`
	HoraInicio string `json:"hora_inicio"`
	HoraFin   string `json:"hora_fin"`
}

func (h *Horario) Validate() error {
	if !h.Turno.Valid() {
		return errors.New("turno debe ser 'Matutina' o 'Vespertina'")
	}
	if len(h.HoraInicio) != 5 || h.HoraInicio[2] != ':' {
		return errors.New("hora_inicio debe tener formato HH:MM")
	}
	if len(h.HoraFin) != 5 || h.HoraFin[2] != ':' {
		return errors.New("hora_fin debe tener formato HH:MM")
	}
	return nil
}

type CarritoHorarioRel struct {
	NumeroCarrito int `json:"numero_carrito"`
	IDHorario     int `json:"id_horario"`
}

type CarritoHorario struct {
	NumeroCarrito     int    `json:"numero_carrito,omitempty"`
	IDHorario         int    `json:"id_horario,omitempty"`
	HoraAsignacion    string `json:"hora_asignacion,omitempty"`
	Turno             string `json:"turno,omitempty"`
	HoraInicio        string `json:"hora_inicio,omitempty"`
	HoraFin           string `json:"hora_fin,omitempty"`
	EstadoCarrito     string `json:"estado_carrito,omitempty"`
	CapacidadPasajeros int   `json:"capacidad_pasajeros,omitempty"`
}

type HorarioRepository interface {
	GetAllHorarios() ([]Horario, error)
	GetHorarioByID(id int) (Horario, error)
	CreateHorario(h Horario) (Horario, error)
	UpdateHorario(h Horario) error
	DeleteHorario(id int) error
	GetCarritosByHorario(id int) ([]CarritoHorario, error)
}

type CarritoHorarioRepository interface {
	AsignarCarritoHorario(rel CarritoHorarioRel) error
	DesasignarCarritoHorario(rel CarritoHorarioRel) error
}
