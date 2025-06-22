package domain

import (
	"context"
	"time"
)

type Event struct {
	ID          int64     `json:"id"`
	Slug        string    `json:"slug"`
	CategoryID  int16     `json:"category_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PlaceID     int64     `json:"place_id"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type EventRepo interface {
	ListActive(ctx context.Context, limit, offset int) ([]Event, error)
	GetBySlug(ctx context.Context, slug string) (Event, error)
	DeactivatePast(ctx context.Context) (int64, error)
	IterChangedSince(ctx context.Context, lastTS time.Time, batch int) ([]Event, error)
	Create(ctx context.Context, ev *Event) error
}
