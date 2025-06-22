package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"kulturaGo/events-service/internal/dto"
	"time"
)

const publicListSQL = `
SELECT json_build_object(
  'id',     e.id,
  'slug',   e.slug,
  'category', json_build_object(
      'slug', COALESCE(NULLIF(c.slug, ''), lower(regexp_replace(c.name, '\s+', '-', 'g'))),
      'name', c.name
  ),
  'title',       e.title,
  'description', e.description,
  'place', json_build_object(
      'address', concat_ws(', ',
          NULLIF(p.title,       ''),
          NULLIF(p.country,     ''),
          NULLIF(p.region,      ''),
          NULLIF(p.city,        ''),
          NULLIF(p.street,      ''),
          NULLIF(p.house_num,   ''),
          NULLIF(p.postal_code, '')
      ),
      'latitude',  p.latitude,
      'longitude', p.longitude
  ),
  'starts_at', to_char(e.starts_at, 'YYYY-MM-DD"T"HH24:MI:SSOF'),
  'ends_at',   to_char(e.ends_at,   'YYYY-MM-DD"T"HH24:MI:SSOF'),
  'photos', COALESCE((
      SELECT json_agg(json_build_object(
               'url',      ep.url,
               'alt_text', ep.alt_text,
               'is_main',  ep.is_main
             ) ORDER BY ep.is_main DESC)
      FROM event_photos ep
      WHERE ep.event_id = e.id
    ), '[]'::json),
  'legal_info', COALESCE((
      SELECT json_agg(json_build_object(
               'key',  li.info_key,
               'text', li.info_text
             ))
      FROM legal_information li
      WHERE li.event_id = e.id
    ), '[]'::json),
  'people', COALESCE((
      SELECT json_agg(json_build_object(
               'slug', COALESCE(NULLIF(per.slug, ''), lower(regexp_replace(per.name, '\s+', '-', 'g'))),
               'name', per.name,
               'tag', json_build_object(
                   'slug', COALESCE(NULLIF(t.slug, ''), lower(regexp_replace(t.name, '\s+', '-', 'g'))),
                   'name', t.name
               )
             ) ORDER BY ep.sort_order)
      FROM event_people ep
      JOIN persons per ON per.id = ep.person_id
      JOIN tags    t   ON t.id   = ep.tag_id
      WHERE ep.event_id = e.id
    ), '[]'::json)
) AS event
FROM events e
JOIN event_categories c ON c.id = e.category_id
JOIN places           p ON p.id = e.place_id
WHERE e.is_active = TRUE
  AND e.ends_at  >= NOW()
ORDER BY e.starts_at
LIMIT  $1
OFFSET $2;
`

const publicBySlugSQL = `
SELECT json_build_object(
  'id',    e.id,
  'slug',  e.slug,
  'category', json_build_object(
      'slug', COALESCE(NULLIF(c.slug, ''), lower(regexp_replace(c.name, '\s+', '-', 'g'))),
      'name', c.name
  ),
  'title',       e.title,
  'description', e.description,
  'place', json_build_object(
      'address', concat_ws(', ',
          NULLIF(p.title,       ''),
          NULLIF(p.country,     ''),
          NULLIF(p.region,      ''),
          NULLIF(p.city,        ''),
          NULLIF(p.street,      ''),
          NULLIF(p.house_num,   ''),
          NULLIF(p.postal_code, '')
      ),
      'latitude',  p.latitude,
      'longitude', p.longitude
  ),
  'starts_at', to_char(e.starts_at, 'YYYY-MM-DD"T"HH24:MI:SSOF'),
  'ends_at',   to_char(e.ends_at,   'YYYY-MM-DD"T"HH24:MI:SSOF'),
  'photos', COALESCE((
      SELECT json_agg(json_build_object(
               'url',      ep.url,
               'alt_text', ep.alt_text,
               'is_main',  ep.is_main
             ) ORDER BY ep.is_main DESC)
      FROM event_photos ep
      WHERE ep.event_id = e.id
    ), '[]'::json),
  'legal_info', COALESCE((
      SELECT json_agg(json_build_object(
               'key',  li.info_key,
               'text', li.info_text
             ))
      FROM legal_information li
      WHERE li.event_id = e.id
    ), '[]'::json),
  'people', COALESCE((
      SELECT json_agg(json_build_object(
               'slug', COALESCE(NULLIF(per.slug, ''), lower(regexp_replace(per.name, '\s+', '-', 'g'))),
               'name', per.name,
               'tag',  json_build_object(
                    'slug', COALESCE(NULLIF(t.slug, ''), lower(regexp_replace(t.name, '\s+', '-', 'g'))),
                    'name', t.name
               )
             ) ORDER BY ep.sort_order)
      FROM event_people ep
      JOIN persons per ON per.id = ep.person_id
      JOIN tags    t   ON t.id   = ep.tag_id
      WHERE ep.event_id = e.id
    ), '[]'::json)
) AS event
FROM events e
JOIN event_categories c ON c.id = e.category_id
JOIN places           p ON p.id = e.place_id
WHERE e.is_active = TRUE
  AND e.ends_at  >= NOW()
  AND e.slug     = $1
LIMIT 1;
`

const iterSQL = `
SELECT json_build_object(
  'id',          e.id,
  'slug',        e.slug,
  'created_at',  to_char(e.created_at, 'YYYY-MM-DD"T"HH24:MI:SSOF'),
  'category',    json_build_object('slug', c.slug, 'name', c.name),
  'title',       e.title,
  'description', e.description,
  'place',       json_build_object(
      'address',   concat_ws(', ',
                        p.title,
                        p.country, p.region, p.city,
                        p.street, p.house_num, p.postal_code),
      'latitude',  p.latitude,
      'longitude', p.longitude
  ),
  'starts_at',   to_char(e.starts_at, 'YYYY-MM-DD"T"HH24:MI:SSOF'),
  'ends_at',     to_char(e.ends_at,   'YYYY-MM-DD"T"HH24:MI:SSOF'),
  'is_active',   e.is_active,
  'photos',      COALESCE((
      SELECT json_agg(json_build_object(
               'url',      ep.url,
               'alt_text', ep.alt_text,
               'is_main',  ep.is_main
             ) ORDER BY ep.is_main DESC)
      FROM event_photos ep
      WHERE ep.event_id = e.id
    ), '[]'::json),
  'legal_info',  COALESCE((
      SELECT json_agg(json_build_object(
               'key',  li.info_key,
               'text', li.info_text
             ))
      FROM legal_information li
      WHERE li.event_id = e.id
    ), '[]'::json),
  'people',      COALESCE((
      SELECT json_agg(json_build_object(
               'slug', per.slug,
               'name', per.name,
               'tag',  json_build_object('slug', t.slug, 'name', t.name)
             ) ORDER BY ep.sort_order)
      FROM event_people ep
      JOIN persons per ON per.id = ep.person_id
      JOIN tags    t   ON t.id   = ep.tag_id
      WHERE ep.event_id = e.id
    ), '[]'::json)
) AS event
FROM events e
JOIN event_categories c ON c.id = e.category_id
JOIN places p           ON p.id = e.place_id
WHERE (e.created_at > $1 OR e.updated_at > $1)
ORDER BY e.created_at
LIMIT  $2;`

func (r *EventRepository) ListPublic(ctx context.Context, limit, offset int) ([]dto.PublicEvent, error) {
	rows, err := r.db.Query(ctx, publicListSQL, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []dto.PublicEvent
	for rows.Next() {
		var raw json.RawMessage
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		var ev dto.PublicEvent
		if err := json.Unmarshal(raw, &ev); err != nil {
			return nil, err
		}
		out = append(out, ev)
	}
	return out, rows.Err()
}

func (r *EventRepository) GetPublicBySlug(ctx context.Context, slug string) (dto.PublicEvent, error) {
	var raw json.RawMessage
	if err := r.db.QueryRow(ctx, publicBySlugSQL, slug).Scan(&raw); err != nil {
		return dto.PublicEvent{}, err
	}
	var ev dto.PublicEvent
	if err := json.Unmarshal(raw, &ev); err != nil {
		return dto.PublicEvent{}, err
	}
	return ev, nil
}

func (r *EventRepository) DeactivatePast(ctx context.Context) (int64, error) {
	cmd, err := r.db.Exec(ctx,
		`UPDATE events SET is_active = false WHERE ends_at < NOW() AND is_active = true`,
	)
	return cmd.RowsAffected(), err
}

func (r *EventRepository) IterChangedSince(ctx context.Context, lastTS time.Time, batch int) ([]dto.PublicEvent, error) {
	rows, err := r.db.Query(ctx, iterSQL, lastTS, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []dto.PublicEvent
	for rows.Next() {
		var raw json.RawMessage
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		var ev dto.PublicEvent
		if err := json.Unmarshal(raw, &ev); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event: %w", err)
		}
		out = append(out, ev)
	}
	return out, rows.Err()
}
