package routes

import (
	handlerhttp "kulturaGo/events-service/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRoutes(event *handlerhttp.EventHandler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/events", event.ListActive)
		api.Get("/events/{slug}", event.GetBySlug)

		api.Post("/events", event.Create)
	})

	return r
}
