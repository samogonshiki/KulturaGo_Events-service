package dto

import "time"

type Category struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type Place struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Photo struct {
	URL     string  `json:"url"`
	AltText *string `json:"alt_text,omitempty"`
	IsMain  bool    `json:"is_main"`
}

type Legal struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

type Tag struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type Person struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	Tag  Tag    `json:"tag"`
}

type PublicEvent struct {
	ID   int64  `json:"id"`
	Slug string `json:"slug"`

	Category    Category `json:"category"`
	Title       string   `json:"title"`
	Description string   `json:"description"`

	CreatedAt time.Time `json:"-"`

	Place Place `json:"place"`

	StartsAt FlexTime `json:"starts_at"`
	EndsAt   FlexTime `json:"ends_at"`

	Photos    []Photo  `json:"photos"`
	LegalInfo []Legal  `json:"legal_info"`
	People    []Person `json:"people"`
}

type CreateEventInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	StartsAt    FlexTime `json:"starts_at"`
	EndsAt      FlexTime `json:"ends_at"`
	IsActive    bool     `json:"is_active"`

	Category Category `json:"category"`
	Place    Place    `json:"place"`

	Photos    []Photo  `json:"photos"`
	LegalInfo []Legal  `json:"legal_info"`
	People    []Person `json:"people"`
}
