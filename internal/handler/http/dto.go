package http

import "time"

type Place struct {
	ID      int64    `json:"id"`
	Title   string   `json:"title"`
	Country *string  `json:"country,omitempty"`
	Region  *string  `json:"region,omitempty"`
	City    string   `json:"city"`
	Street  *string  `json:"street,omitempty"`
	House   *string  `json:"house_num,omitempty"`
	Lat     *float64 `json:"latitude,omitempty"`
	Lon     *float64 `json:"longitude,omitempty"`
}

type Event struct {
	ID          int64     `json:"id"`
	Slug        string    `json:"slug"`
	CategoryID  int16     `json:"category_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Place       Place     `json:"place"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
}
