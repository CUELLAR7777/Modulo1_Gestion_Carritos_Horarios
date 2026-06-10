package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"uleam-transporte/internal/middleware"
	"uleam-transporte/internal/storage"
)

type Server struct {
	carritoHandler *CarritoHandler
	horarioHandler *HorarioHandler
}

func NewServer(almacen storage.Almacen) *Server {
	return &Server{
		carritoHandler: NewCarritoHandler(almacen),
		horarioHandler: NewHorarioHandler(almacen),
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(middleware.CORSOptions()))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/carritos", func(r chi.Router) {
			r.Get("/", s.carritoHandler.List)
			r.Post("/", s.carritoHandler.Create)

			r.Route("/{numero}", func(r chi.Router) {
				r.Get("/", s.carritoHandler.GetByNumero)
				r.Put("/", s.carritoHandler.Update)
				r.Delete("/", s.carritoHandler.Delete)
				r.Get("/horarios", s.carritoHandler.GetHorarios)
				r.Post("/horarios", s.carritoHandler.AsignarHorario)
				r.Delete("/horarios/{idHorario}", s.carritoHandler.DesasignarHorario)
			})
		})

		r.Route("/horarios", func(r chi.Router) {
			r.Get("/", s.horarioHandler.List)
			r.Post("/", s.horarioHandler.Create)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", s.horarioHandler.GetByID)
				r.Put("/", s.horarioHandler.Update)
				r.Delete("/", s.horarioHandler.Delete)
				r.Get("/carritos", s.horarioHandler.GetCarritos)
			})
		})
	})

	return r
}
