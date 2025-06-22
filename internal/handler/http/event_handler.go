package http

import (
	"context"
	"encoding/json"
	"kulturaGo/events-service/internal/domain"
	nethttp "net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	Repo domain.EventRepo
}

func NewEventHandler(repo domain.EventRepo) *EventHandler { return &EventHandler{Repo: repo} }

func (h *EventHandler) ListActive(w nethttp.ResponseWriter, r *nethttp.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	evs, err := h.Repo.ListActive(r.Context(), limit, offset)
	if err != nil {
		nethttp.Error(w, err.Error(), 500)
		return
	}
	writeJSON(r.Context(), w, evs)
}

func (h *EventHandler) GetBySlug(w nethttp.ResponseWriter, r *nethttp.Request) {
	slug := chi.URLParam(r, "slug")
	ev, err := h.Repo.GetBySlug(r.Context(), slug)
	if err != nil {
		nethttp.NotFound(w, r)
		return
	}
	writeJSON(r.Context(), w, ev)
}

func writeJSON(_ context.Context, w nethttp.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

func (h *EventHandler) Create(w nethttp.ResponseWriter, r *nethttp.Request) {
	var in struct {
		Title       string    `json:"title"`
		CategoryID  int16     `json:"category_id"`
		Description string    `json:"description"`
		PlaceID     int64     `json:"place_id"`
		StartsAt    time.Time `json:"starts_at"`
		EndsAt      time.Time `json:"ends_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusBadRequest)
		return
	}

	ev := domain.Event{
		Title:       in.Title,
		CategoryID:  in.CategoryID,
		Description: in.Description,
		PlaceID:     in.PlaceID,
		StartsAt:    in.StartsAt,
		EndsAt:      in.EndsAt,
	}
	if err := h.Repo.Create(r.Context(), &ev); err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/api/v1/events/"+ev.Slug)
	w.WriteHeader(nethttp.StatusCreated)
	writeJSON(r.Context(), w, ev)
}
