package http

import (
	"context"
	"encoding/json"
	"kulturaGo/events-service/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	Repo domain.EventRepo
}

func NewEventHandler(repo domain.EventRepo) *EventHandler { return &EventHandler{Repo: repo} }

func (h *EventHandler) ListActive(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	evs, err := h.Repo.ListActive(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(r.Context(), w, evs)
}

func (h *EventHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	ev, err := h.Repo.GetBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(r.Context(), w, ev)
}

func writeJSON(_ context.Context, w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}
