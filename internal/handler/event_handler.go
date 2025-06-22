package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"kulturaGo/events-service/internal/domain"
	"kulturaGo/events-service/internal/dto"
	"net/http"
	"strconv"
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

	evs, err := h.Repo.ListPublic(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(r.Context(), w, evs)
}

func (h *EventHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	ev, err := h.Repo.GetPublicBySlug(r.Context(), slug)
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

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pubEvent, err := h.Repo.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/api/v1/events/"+pubEvent.Slug)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(pubEvent); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
