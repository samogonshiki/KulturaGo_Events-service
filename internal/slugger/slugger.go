package slugger

import (
	"context"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Generate(ctx context.Context, db *pgxpool.Pool, title string) (string, error) {
	base := slug.MakeLang(title, "ru")

	if base == "" {
		base = "event"
	}

	candidate := base
	var i int

	for {
		var exists bool
		err := db.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM events WHERE slug = $1)`,
			candidate).Scan(&exists)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		i++
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}
