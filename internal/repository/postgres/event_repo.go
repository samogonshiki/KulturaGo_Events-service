package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"kulturaGo/events-service/internal/domain"
	"kulturaGo/events-service/internal/slugger"
	"time"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository { return &EventRepository{db: db} }

const listActiveSQL = `
SELECT id, slug, category_id, title, description, place_id,
       starts_at, ends_at, is_active, created_at
FROM events
WHERE is_active = true
  AND ends_at >= NOW()
ORDER BY starts_at
LIMIT $1 OFFSET $2
`

func (r *EventRepository) Create(ctx context.Context, ev *domain.Event) error {
	s, err := slugger.Generate(ctx, r.db, ev.Title)
	if err != nil {
		return err
	}
	ev.Slug = s

	const q = `
        INSERT INTO events
            (slug, category_id, title, description,
             place_id, starts_at, ends_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7)
        RETURNING id, created_at`
	return r.db.QueryRow(ctx, q,
		ev.Slug, ev.CategoryID, ev.Title, ev.Description,
		ev.PlaceID, ev.StartsAt, ev.EndsAt,
	).Scan(&ev.ID, &ev.CreatedAt)
}

func (r *EventRepository) ListActive(ctx context.Context, limit, offset int) ([]domain.Event, error) {
	rows, err := r.db.Query(ctx, listActiveSQL, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Event
	for rows.Next() {
		var ev domain.Event
		if err = rows.Scan(&ev.ID, &ev.Slug, &ev.CategoryID, &ev.Title, &ev.Description,
			&ev.PlaceID, &ev.StartsAt, &ev.EndsAt, &ev.IsActive, &ev.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, ev)
	}
	return res, rows.Err()
}

func (r *EventRepository) GetBySlug(ctx context.Context, slug string) (domain.Event, error) {
	const q = `SELECT id, slug, category_id, title, description, place_id,
		         starts_at, ends_at, is_active, created_at
		       FROM events WHERE slug = $1`
	var ev domain.Event
	err := r.db.QueryRow(ctx, q, slug).Scan(&ev.ID, &ev.Slug, &ev.CategoryID, &ev.Title,
		&ev.Description, &ev.PlaceID, &ev.StartsAt, &ev.EndsAt, &ev.IsActive, &ev.CreatedAt)
	return ev, err
}

func (r *EventRepository) DeactivatePast(ctx context.Context) (int64, error) {
	cmd, err := r.db.Exec(ctx,
		`UPDATE events SET is_active = false
		  WHERE is_active = true AND ends_at < NOW()`)
	return cmd.RowsAffected(), err
}

func (r *EventRepository) IterChangedSince(ctx context.Context, ts time.Time, batch int) ([]domain.Event, error) {
	const q = `
SELECT id, slug, category_id, title, description, place_id,
       starts_at, ends_at, is_active, created_at
FROM events
WHERE created_at > $1 OR updated_at > $1
ORDER BY created_at
LIMIT $2`
	rows, err := r.db.Query(ctx, q, ts, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Event
	for rows.Next() {
		var ev domain.Event
		if err = rows.Scan(&ev.ID, &ev.Slug, &ev.CategoryID, &ev.Title,
			&ev.Description, &ev.PlaceID, &ev.StartsAt, &ev.EndsAt,
			&ev.IsActive, &ev.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, ev)
	}
	return res, rows.Err()
}
