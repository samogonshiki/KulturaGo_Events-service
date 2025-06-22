package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"kulturaGo/events-service/internal/dto"
	"kulturaGo/events-service/internal/slugger"
	"time"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, in dto.CreateEventInput) (dto.PublicEvent, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return dto.PublicEvent{}, err
	}
	defer tx.Rollback(ctx)

	var categoryID int64
	if err := tx.QueryRow(ctx, `
        INSERT INTO event_categories (slug, name)
        VALUES ($1, $2)
        ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
        RETURNING id
    `, in.Category.Slug, in.Category.Name).Scan(&categoryID); err != nil {
		return dto.PublicEvent{}, fmt.Errorf("upsert category: %w", err)
	}

	var placeID int64
	if err := tx.QueryRow(ctx, `
        INSERT INTO places (
          title, country, region, city, street,
          house_num, postal_code, latitude, longitude
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        RETURNING id
    `,
		in.Place.Address, "", "", "", "", "", "", in.Place.Latitude, in.Place.Longitude,
	).Scan(&placeID); err != nil {
		return dto.PublicEvent{}, fmt.Errorf("insert place: %w", err)
	}

	slug, err := slugger.Generate(ctx, r.db, in.Title)
	if err != nil {
		return dto.PublicEvent{}, fmt.Errorf("generate slug: %w", err)
	}

	startsAt := time.Time(in.StartsAt)
	endsAt := time.Time(in.EndsAt)

	var eventID int64
	var createdAt time.Time
	if err := tx.QueryRow(ctx, `
        INSERT INTO events (
          slug, category_id, title, description,
          place_id, starts_at, ends_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7)
        RETURNING id, created_at
    `,
		slug,
		categoryID,
		in.Title,
		in.Description,
		placeID,
		startsAt,
		endsAt,
	).Scan(&eventID, &createdAt); err != nil {
		return dto.PublicEvent{}, fmt.Errorf("insert event: %w", err)
	}

	for idx, p := range in.People {
		var tagID int64
		if err := tx.QueryRow(ctx, `
            INSERT INTO tags (slug, name)
            VALUES ($1, $2)
            ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
            RETURNING id
        `, p.Tag.Slug, p.Tag.Name).Scan(&tagID); err != nil {
			return dto.PublicEvent{}, fmt.Errorf("upsert tag: %w", err)
		}

		var personID int64
		if err := tx.QueryRow(ctx, `
            INSERT INTO persons (slug, name, description, photo)
            VALUES ($1, $2, '', '')
            ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
            RETURNING id
        `, p.Slug, p.Name).Scan(&personID); err != nil {
			return dto.PublicEvent{}, fmt.Errorf("upsert person: %w", err)
		}

		if _, err := tx.Exec(ctx, `
            INSERT INTO event_people (event_id, person_id, tag_id, sort_order)
            VALUES ($1,$2,$3,$4)
        `, eventID, personID, tagID, idx); err != nil {
			return dto.PublicEvent{}, fmt.Errorf("insert event_person: %w", err)
		}
	}

	for _, ph := range in.Photos {
		if _, err := tx.Exec(ctx, `
            INSERT INTO event_photos (event_id, url, alt_text, is_main)
            VALUES ($1,$2,$3,$4)
        `, eventID, ph.URL, ph.AltText, ph.IsMain); err != nil {
			return dto.PublicEvent{}, fmt.Errorf("insert photo: %w", err)
		}
	}
	for _, li := range in.LegalInfo {
		if _, err := tx.Exec(ctx, `
            INSERT INTO legal_information (event_id, info_key, info_text)
            VALUES ($1,$2,$3)
        `, eventID, li.Key, li.Text); err != nil {
			return dto.PublicEvent{}, fmt.Errorf("insert legal_info: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.PublicEvent{}, err
	}

	return r.GetPublicBySlug(ctx, slug)
}
