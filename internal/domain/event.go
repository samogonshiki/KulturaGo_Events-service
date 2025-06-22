package domain

import (
	"context"
	"kulturaGo/events-service/internal/dto"
	"time"
)

type EventRepo interface {
	ListPublic(ctx context.Context, limit, offset int) ([]dto.PublicEvent, error)
	GetPublicBySlug(ctx context.Context, slug string) (dto.PublicEvent, error)
	Create(ctx context.Context, in dto.CreateEventInput) (dto.PublicEvent, error)
	DeactivatePast(ctx context.Context) (int64, error)
	IterChangedSince(ctx context.Context, lastTS time.Time, batch int) ([]dto.PublicEvent, error)
}
