package postgres

import (
	"context"
	"fmt"
)

type PublicPlace struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Country   *string  `json:"country,omitempty"`
	Region    *string  `json:"region,omitempty"`
	City      string   `json:"city"`
	Street    *string  `json:"street,omitempty"`
	HouseNum  *string  `json:"house_num,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

type PublicEvent struct {
	ID          int64       `json:"id"`
	Slug        string      `json:"slug"`
	CategoryID  int16       `json:"category_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Place       PublicPlace `json:"place"`
	StartsAt    string      `json:"starts_at"`
	EndsAt      string      `json:"ends_at"`
}

const publicSelect = `
SELECT 
  e.id, e.slug, e.category_id, e.title, e.description,
  e.starts_at, e.ends_at,
  p.id, p.title, p.country, p.region, p.city, p.street,
  p.house_num, p.latitude, p.longitude
FROM events       AS e
JOIN places       AS p ON p.id = e.place_id
WHERE e.is_active = TRUE
  AND e.ends_at  >= NOW()
`

func (r *EventRepository) ListPublic(ctx context.Context, limit, offset int) ([]PublicEvent, error) {
	sql := publicSelect + "ORDER BY e.starts_at LIMIT $1 OFFSET $2"
	rows, err := r.db.Query(ctx, sql, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PublicEvent
	for rows.Next() {
		var ev PublicEvent
		err = rows.Scan(
			&ev.ID, &ev.Slug, &ev.CategoryID, &ev.Title, &ev.Description,
			&ev.StartsAt, &ev.EndsAt,
			&ev.Place.ID, &ev.Place.Title, &ev.Place.Country, &ev.Place.Region, &ev.Place.City,
			&ev.Place.Street, &ev.Place.HouseNum, &ev.Place.Latitude, &ev.Place.Longitude,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, ev)
	}
	return out, rows.Err()
}

func (r *EventRepository) GetPublicBySlug(ctx context.Context, slug string) (PublicEvent, error) {
	sql := publicSelect + " AND e.slug = $1 LIMIT 1"
	var ev PublicEvent
	err := r.db.QueryRow(ctx, sql, slug).Scan(
		&ev.ID, &ev.Slug, &ev.CategoryID, &ev.Title, &ev.Description,
		&ev.StartsAt, &ev.EndsAt,
		&ev.Place.ID, &ev.Place.Title, &ev.Place.Country, &ev.Place.Region, &ev.Place.City,
		&ev.Place.Street, &ev.Place.HouseNum, &ev.Place.Latitude, &ev.Place.Longitude,
	)
	if err != nil {
		return PublicEvent{}, fmt.Errorf("get event by slug: %w", err)
	}
	return ev, nil
}
