package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"uleam-transporte/internal/dominio"
	"uleam-transporte/internal/storage"
)

type HorarioHandler struct {
	almacen storage.Almacen
}

func NewHorarioHandler(almacen storage.Almacen) *HorarioHandler {
	return &HorarioHandler{almacen: almacen}
}

func (h *HorarioHandler) List(w http.ResponseWriter, r *http.Request) {
	horarios, err := h.almacen.GetAllHorarios()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, horarios)
}

func (h *HorarioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id_horario debe ser un entero")
		return
	}

	horario, err := h.almacen.GetHorarioByID(id)
	if err == storage.ErrNotFound {
		RespondError(w, http.StatusNotFound, "horario no encontrado")
		return
	}
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, horario)
}

func (h *HorarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	var horario dominio.Horario
	if err := json.NewDecoder(r.Body).Decode(&horario); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	if err := horario.Validate(); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	creado, err := h.almacen.CreateHorario(horario)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, creado)
}

func (h *HorarioHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id_horario debe ser un entero")
		return
	}

	var horario dominio.Horario
	if err := json.NewDecoder(r.Body).Decode(&horario); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON invalido")
		return
	}
	horario.ID = id

	if err := horario.Validate(); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.almacen.UpdateHorario(horario); err != nil {
		if err == storage.ErrNotFound {
			RespondError(w, http.StatusNotFound, "horario no encontrado")
			return
		}
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, horario)
}

func (h *HorarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id_horario debe ser un entero")
		return
	}

	if err := h.almacen.DeleteHorario(id); err != nil {
		if err == storage.ErrNotFound {
			RespondError(w, http.StatusNotFound, "horario no encontrado")
			return
		}
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HorarioHandler) GetCarritos(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id_horario debe ser un entero")
		return
	}

	carritos, err := h.almacen.GetCarritosByHorario(id)
	if err == storage.ErrNotFound {
		RespondError(w, http.StatusNotFound, "horario no encontrado")
		return
	}
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, carritos)
}
