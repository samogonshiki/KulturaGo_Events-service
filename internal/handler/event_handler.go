package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"kulturaGo/events-service/internal/domain"
	"kulturaGo/events-service/internal/dto"
	nethttp "net/http"
	"strconv"
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

	evs, err := h.Repo.ListPublic(r.Context(), limit, offset)
	if err != nil {
		nethttp.Error(w, err.Error(), 500)
		return
	}
	writeJSON(r.Context(), w, evs)
}

func (h *EventHandler) GetBySlug(w nethttp.ResponseWriter, r *nethttp.Request) {
	slug := chi.URLParam(r, "slug")
	ev, err := h.Repo.GetPublicBySlug(r.Context(), slug)
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
	ctx := r.Context()
	var input dto.PublicEvent
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusBadRequest)
		return
	}

	pubEvent, err := h.Repo.Create(ctx, input)
	if err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
		return
	}

	location := fmt.Sprintf("/api/v1/events/%s", pubEvent.Slug)
	w.Header().Set("Location", location)

	w.WriteHeader(nethttp.StatusCreated)
	if err := json.NewEncoder(w).Encode(pubEvent); err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
	}
}
