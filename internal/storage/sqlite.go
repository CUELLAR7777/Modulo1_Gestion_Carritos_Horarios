package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"uleam-transporte/internal/dominio"
)

type SQLiteAlmacen struct {
	db *sql.DB
}

func NewSQLiteAlmacen(db *sql.DB) *SQLiteAlmacen {
	return &SQLiteAlmacen{db: db}
}

func (s *SQLiteAlmacen) Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS horarios (
			id_horario INTEGER PRIMARY KEY AUTOINCREMENT,
			turno TEXT NOT NULL CHECK(turno IN ('Matutina', 'Vespertina')),
			hora_inicio TEXT NOT NULL,
			hora_fin TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS carritos (
			numero_carrito INTEGER PRIMARY KEY,
			estado_carrito TEXT NOT NULL DEFAULT 'Disponible' CHECK(estado_carrito IN ('Disponible', 'En Viaje', 'Mantenimiento')),
			capacidad_pasajeros INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS carrito_horario (
			numero_carrito INTEGER NOT NULL,
			id_horario INTEGER NOT NULL,
			hora_asignacion TEXT NOT NULL DEFAULT (datetime('now')),
			PRIMARY KEY (numero_carrito, id_horario),
			FOREIGN KEY (numero_carrito) REFERENCES carritos(numero_carrito) ON DELETE CASCADE,
			FOREIGN KEY (id_horario) REFERENCES horarios(id_horario) ON DELETE CASCADE
		)`,
	}
	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return fmt.Errorf("migracion: %w", err)
		}
	}
	return nil
}

func (s *SQLiteAlmacen) Seed() error {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM carritos").Scan(&count); err != nil {
		return fmt.Errorf("seed count: %w", err)
	}
	if count > 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	horarios := []dominio.Horario{
		{Turno: dominio.TurnoMatutino, HoraInicio: "07:00", HoraFin: "12:00"},
		{Turno: dominio.TurnoVespertino, HoraInicio: "13:00", HoraFin: "18:00"},
	}
	for _, h := range horarios {
		if _, err := tx.Exec("INSERT INTO horarios (turno, hora_inicio, hora_fin) VALUES (?, ?, ?)",
			h.Turno, h.HoraInicio, h.HoraFin); err != nil {
			return fmt.Errorf("seed horario: %w", err)
		}
	}

	carritos := []dominio.Carrito{
		{Numero: 101, Estado: dominio.EstadoDisponible, CapacidadPasajeros: 15},
		{Numero: 102, Estado: dominio.EstadoEnViaje, CapacidadPasajeros: 20},
		{Numero: 103, Estado: dominio.EstadoMantenimiento, CapacidadPasajeros: 12},
		{Numero: 104, Estado: dominio.EstadoDisponible, CapacidadPasajeros: 18},
	}
	for _, c := range carritos {
		if _, err := tx.Exec("INSERT INTO carritos (numero_carrito, estado_carrito, capacidad_pasajeros) VALUES (?, ?, ?)",
			c.Numero, c.Estado, c.CapacidadPasajeros); err != nil {
			return fmt.Errorf("seed carrito: %w", err)
		}
	}

	asignaciones := []struct {
		carrito int
		horario int
	}{
		{101, 1}, {102, 1}, {102, 2}, {104, 2},
	}
	for _, a := range asignaciones {
		if _, err := tx.Exec("INSERT OR IGNORE INTO carrito_horario (numero_carrito, id_horario) VALUES (?, ?)",
			a.carrito, a.horario); err != nil {
			return fmt.Errorf("seed asignacion: %w", err)
		}
	}

	return tx.Commit()
}

// ---- CarritoRepository ----

func (s *SQLiteAlmacen) GetAllCarritos() ([]dominio.Carrito, error) {
	rows, err := s.db.Query("SELECT numero_carrito, estado_carrito, capacidad_pasajeros FROM carritos ORDER BY numero_carrito")
	if err != nil {
		return nil, fmt.Errorf("get all carritos: %w", err)
	}
	defer rows.Close()

	var carritos []dominio.Carrito
	for rows.Next() {
		var c dominio.Carrito
		if err := rows.Scan(&c.Numero, &c.Estado, &c.CapacidadPasajeros); err != nil {
			return nil, fmt.Errorf("scan carrito: %w", err)
		}
		carritos = append(carritos, c)
	}
	if carritos == nil {
		carritos = []dominio.Carrito{}
	}
	return carritos, rows.Err()
}

func (s *SQLiteAlmacen) GetCarritoByNumero(numero int) (dominio.Carrito, error) {
	var c dominio.Carrito
	err := s.db.QueryRow(
		"SELECT numero_carrito, estado_carrito, capacidad_pasajeros FROM carritos WHERE numero_carrito = ?",
		numero,
	).Scan(&c.Numero, &c.Estado, &c.CapacidadPasajeros)
	if err == sql.ErrNoRows {
		return c, ErrNotFound
	}
	if err != nil {
		return c, fmt.Errorf("get carrito: %w", err)
	}
	return c, nil
}

func (s *SQLiteAlmacen) CreateCarrito(c dominio.Carrito) error {
	_, err := s.db.Exec(
		"INSERT INTO carritos (numero_carrito, estado_carrito, capacidad_pasajeros) VALUES (?, ?, ?)",
		c.Numero, c.Estado, c.CapacidadPasajeros,
	)
	if err != nil {
		if IsConstraintError(err) {
			return ErrConflict
		}
		return fmt.Errorf("create carrito: %w", err)
	}
	return nil
}

func (s *SQLiteAlmacen) UpdateCarrito(c dominio.Carrito) error {
	result, err := s.db.Exec(
		"UPDATE carritos SET estado_carrito = ?, capacidad_pasajeros = ? WHERE numero_carrito = ?",
		c.Estado, c.CapacidadPasajeros, c.Numero,
	)
	if err != nil {
		return fmt.Errorf("update carrito: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLiteAlmacen) DeleteCarrito(numero int) error {
	result, err := s.db.Exec("DELETE FROM carritos WHERE numero_carrito = ?", numero)
	if err != nil {
		return fmt.Errorf("delete carrito: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLiteAlmacen) GetHorariosByCarrito(numero int) ([]dominio.CarritoHorario, error) {
	rows, err := s.db.Query(`
		SELECT ch.numero_carrito, ch.id_horario, ch.hora_asignacion,
		       h.turno, h.hora_inicio, h.hora_fin
		FROM carrito_horario ch
		JOIN horarios h ON h.id_horario = ch.id_horario
		WHERE ch.numero_carrito = ?
		ORDER BY h.turno, h.hora_inicio
	`, numero)
	if err != nil {
		return nil, fmt.Errorf("get horarios by carrito: %w", err)
	}
	defer rows.Close()

	var result []dominio.CarritoHorario
	for rows.Next() {
		var ch dominio.CarritoHorario
		if err := rows.Scan(&ch.NumeroCarrito, &ch.IDHorario, &ch.HoraAsignacion,
			&ch.Turno, &ch.HoraInicio, &ch.HoraFin); err != nil {
			return nil, fmt.Errorf("scan horario: %w", err)
		}
		result = append(result, ch)
	}
	if result == nil {
		result = []dominio.CarritoHorario{}
	}
	return result, rows.Err()
}

// ---- HorarioRepository ----

func (s *SQLiteAlmacen) GetAllHorarios() ([]dominio.Horario, error) {
	rows, err := s.db.Query("SELECT id_horario, turno, hora_inicio, hora_fin FROM horarios ORDER BY id_horario")
	if err != nil {
		return nil, fmt.Errorf("get all horarios: %w", err)
	}
	defer rows.Close()

	var horarios []dominio.Horario
	for rows.Next() {
		var h dominio.Horario
		if err := rows.Scan(&h.ID, &h.Turno, &h.HoraInicio, &h.HoraFin); err != nil {
			return nil, fmt.Errorf("scan horario: %w", err)
		}
		horarios = append(horarios, h)
	}
	if horarios == nil {
		horarios = []dominio.Horario{}
	}
	return horarios, rows.Err()
}

func (s *SQLiteAlmacen) GetHorarioByID(id int) (dominio.Horario, error) {
	var h dominio.Horario
	err := s.db.QueryRow(
		"SELECT id_horario, turno, hora_inicio, hora_fin FROM horarios WHERE id_horario = ?",
		id,
	).Scan(&h.ID, &h.Turno, &h.HoraInicio, &h.HoraFin)
	if err == sql.ErrNoRows {
		return h, ErrNotFound
	}
	if err != nil {
		return h, fmt.Errorf("get horario: %w", err)
	}
	return h, nil
}

func (s *SQLiteAlmacen) CreateHorario(h dominio.Horario) (dominio.Horario, error) {
	result, err := s.db.Exec(
		"INSERT INTO horarios (turno, hora_inicio, hora_fin) VALUES (?, ?, ?)",
		h.Turno, h.HoraInicio, h.HoraFin,
	)
	if err != nil {
		return h, fmt.Errorf("create horario: %w", err)
	}
	id, _ := result.LastInsertId()
	h.ID = int(id)
	return h, nil
}

func (s *SQLiteAlmacen) UpdateHorario(h dominio.Horario) error {
	result, err := s.db.Exec(
		"UPDATE horarios SET turno = ?, hora_inicio = ?, hora_fin = ? WHERE id_horario = ?",
		h.Turno, h.HoraInicio, h.HoraFin, h.ID,
	)
	if err != nil {
		return fmt.Errorf("update horario: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLiteAlmacen) DeleteHorario(id int) error {
	result, err := s.db.Exec("DELETE FROM horarios WHERE id_horario = ?", id)
	if err != nil {
		return fmt.Errorf("delete horario: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLiteAlmacen) GetCarritosByHorario(id int) ([]dominio.CarritoHorario, error) {
	rows, err := s.db.Query(`
		SELECT ch.numero_carrito, ch.id_horario, ch.hora_asignacion,
		       c.estado_carrito, c.capacidad_pasajeros
		FROM carrito_horario ch
		JOIN carritos c ON c.numero_carrito = ch.numero_carrito
		WHERE ch.id_horario = ?
		ORDER BY c.numero_carrito
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get carritos by horario: %w", err)
	}
	defer rows.Close()

	var result []dominio.CarritoHorario
	for rows.Next() {
		var ch dominio.CarritoHorario
		if err := rows.Scan(&ch.NumeroCarrito, &ch.IDHorario, &ch.HoraAsignacion,
			&ch.EstadoCarrito, &ch.CapacidadPasajeros); err != nil {
			return nil, fmt.Errorf("scan carrito: %w", err)
		}
		result = append(result, ch)
	}
	if result == nil {
		result = []dominio.CarritoHorario{}
	}
	return result, rows.Err()
}

// ---- CarritoHorarioRepository ----

func (s *SQLiteAlmacen) AsignarCarritoHorario(rel dominio.CarritoHorarioRel) error {
	_, err := s.db.Exec(
		"INSERT INTO carrito_horario (numero_carrito, id_horario) VALUES (?, ?)",
		rel.NumeroCarrito, rel.IDHorario,
	)
	if err != nil {
		if IsConstraintError(err) {
			if strings.Contains(err.Error(), "PRIMARY") {
				return ErrConflict
			}
			return fmt.Errorf("violacion de llave foranea: verifique que el carrito y horario existan")
		}
		return fmt.Errorf("asignar carrito horario: %w", err)
	}
	return nil
}

func (s *SQLiteAlmacen) DesasignarCarritoHorario(rel dominio.CarritoHorarioRel) error {
	result, err := s.db.Exec(
		"DELETE FROM carrito_horario WHERE numero_carrito = ? AND id_horario = ?",
		rel.NumeroCarrito, rel.IDHorario,
	)
	if err != nil {
		return fmt.Errorf("desasignar carrito horario: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func IsConstraintError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "PRIMARY KEY")
}
