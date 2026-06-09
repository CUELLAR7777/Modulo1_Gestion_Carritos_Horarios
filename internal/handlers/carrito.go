package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"uleam-transporte/internal/dominio"
	"uleam-transporte/internal/storage"
)

type CarritoHandler struct {
	almacen storage.Almacen
}

func NewCarritoHandler(almacen storage.Almacen) *CarritoHandler {
	return &CarritoHandler{almacen: almacen}
}

func (h *CarritoHandler) List(w http.ResponseWriter, r *http.Request) {
	carritos, err := h.almacen.GetAllCarritos()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, carritos)
}

func (h *CarritoHandler) GetByNumero(w http.ResponseWriter, r *http.Request) {
	num, err := strconv.Atoi(chi.URLParam(r, "numero"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "numero_carrito debe ser un entero")
		return
	}

	carrito, err := h.almacen.GetCarritoByNumero(num)
	if err == storage.ErrNotFound {
		RespondError(w, http.StatusNotFound, "carrito no encontrado")
		return
	}
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, carrito)
}

func (h *CarritoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c dominio.Carrito
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if err := c.Validate(); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.almacen.CreateCarrito(c); err != nil {
		if err == storage.ErrConflict {
			RespondError(w, http.StatusConflict, "el carrito ya existe")
			return
		}
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, c)
}

func (h *CarritoHandler) Update(w http.ResponseWriter, r *http.Request) {
	num, err := strconv.Atoi(chi.URLParam(r, "numero"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "numero_carrito debe ser un entero")
		return
	}

	var c dominio.Carrito
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	c.Numero = num

	if err := c.Validate(); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.almacen.UpdateCarrito(c); err != nil {
		if err == storage.ErrNotFound {
			RespondError(w, http.StatusNotFound, "carrito no encontrado")
			return
		}
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, c)
}

func (h *CarritoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	num, err := strconv.Atoi(chi.URLParam(r, "numero"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "numero_carrito debe ser un entero")
		return
	}

	if err := h.almacen.DeleteCarrito(num); err != nil {
		if err == storage.ErrNotFound {
			RespondError(w, http.StatusNotFound, "carrito no encontrado")
			return
		}
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CarritoHandler) GetHorarios(w http.ResponseWriter, r *http.Request) {
	num, err := strconv.Atoi(chi.URLParam(r, "numero"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "numero_carrito debe ser un entero")
		return
	}

	horarios, err := h.almacen.GetHorariosByCarrito(num)
	if err == storage.ErrNotFound {
		RespondError(w, http.StatusNotFound, "carrito no encontrado")
		return
	}
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, horarios)
}

func (h *CarritoHandler) AsignarHorario(w http.ResponseWriter, r *http.Request) {
	num, err := strconv.Atoi(chi.URLParam(r, "numero"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "numero_carrito debe ser un entero")
		return
	}

	var body struct {
		IDHorario int `json:"id_horario"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	rel := dominio.CarritoHorarioRel{
		NumeroCarrito: num,
		IDHorario:     body.IDHorario,
	}

	if err := h.almacen.AsignarCarritoHorario(rel); err != nil {
		if err == storage.ErrConflict {
			RespondError(w, http.StatusConflict, "la asignacion ya existe")
			return
		}
		if err == storage.ErrNotFound {
			RespondError(w, http.StatusNotFound, "carrito u horario no encontrado")
			return
		}
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, rel)
}

func (h *CarritoHandler) DesasignarHorario(w http.ResponseWriter, r *http.Request) {
	num, err := strconv.Atoi(chi.URLParam(r, "numero"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "numero_carrito debe ser un entero")
		return
	}

	idHorario, err := strconv.Atoi(chi.URLParam(r, "idHorario"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id_horario debe ser un entero")
		return
	}

	rel := dominio.CarritoHorarioRel{
		NumeroCarrito: num,
		IDHorario:     idHorario,
	}

	if err := h.almacen.DesasignarCarritoHorario(rel); err != nil {
		if err == storage.ErrNotFound {
			RespondError(w, http.StatusNotFound, "asignacion no encontrada")
			return
		}
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
