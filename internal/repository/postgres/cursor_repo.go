package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CursorRepository struct{ db *pgxpool.Pool }

func NewCursorRepository(db *pgxpool.Pool) *CursorRepository {
	return &CursorRepository{db: db}
}

func (r *CursorRepository) Get(ctx context.Context, consumer string) (time.Time, error) {
	var ts time.Time
	err := r.db.QueryRow(ctx,
		`SELECT last_event_ts
		   FROM export_cursors
		  WHERE consumer = $1`, consumer,
	).Scan(&ts)
	if err != nil {
		return time.Time{}, nil
	}
	return ts, nil
}

func (r *CursorRepository) Update(ctx context.Context, consumer string, ts time.Time) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO export_cursors (consumer, last_event_ts)
		     VALUES ($1, $2)
		ON CONFLICT (consumer) DO UPDATE
		     SET last_event_ts = EXCLUDED.last_event_ts`,
		consumer, ts,
	)
	return err
}
