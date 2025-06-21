CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE event_categories (
                                  id   smallserial  PRIMARY KEY,
                                  slug varchar(32)  NOT NULL UNIQUE,
                                  name varchar(64)  NOT NULL
);

INSERT INTO event_categories (slug, name) VALUES
                                              ('theatre','Театр'), ('show','Шоу'), ('museum','Музей'),
                                              ('concert','Концерт'), ('contest','Конкурс'),
                                              ('olympiad','Олимпиада'), ('other','Другое');

CREATE TABLE places (
                        id          bigserial     PRIMARY KEY,
                        title       varchar(128)  NOT NULL,
                        country     varchar(64),
                        region      varchar(64),
                        city        varchar(128)  NOT NULL,
                        street      varchar(256),
                        house_num   varchar(16),
                        postal_code varchar(16),
                        longitude   numeric(9,6),
                        latitude    numeric(9,6),
                        created_at  timestamptz   NOT NULL DEFAULT now()
);

CREATE TABLE events (
                        id          bigserial     PRIMARY KEY,
                        slug        varchar(128)  NOT NULL UNIQUE,
                        category_id smallint      NOT NULL REFERENCES event_categories(id),
                        title       varchar(64)   NOT NULL,
                        description varchar(4096) NOT NULL,
                        place_id    bigint        NOT NULL REFERENCES places(id) ON DELETE RESTRICT,
                        starts_at   timestamptz   NOT NULL,
                        ends_at     timestamptz   NOT NULL,
                        is_active   boolean       NOT NULL DEFAULT TRUE,
                        created_at  timestamptz   NOT NULL DEFAULT now()
);

CREATE TABLE legal_information (
                                   id        bigserial     PRIMARY KEY,
                                   event_id  bigint        NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                   info_key  varchar(64)   NOT NULL,
                                   info_text text          NOT NULL
);

CREATE TABLE event_photos (
                              id         bigserial    PRIMARY KEY,
                              event_id   bigint       NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                              url        text         NOT NULL,
                              alt_text   varchar(256),
                              is_main    boolean      NOT NULL DEFAULT FALSE,
                              created_at timestamptz  NOT NULL DEFAULT now()
);

CREATE TABLE persons (
                         id          bigserial     PRIMARY KEY,
                         slug        varchar(128)  NOT NULL UNIQUE,
                         name        varchar(256)  NOT NULL,
                         description text,
                         photo       text          NOT NULL,
                         created_at  timestamptz   NOT NULL DEFAULT now()
);

CREATE TABLE tags (
                      id   smallserial  PRIMARY KEY,
                      slug varchar(32)  NOT NULL UNIQUE,
                      name varchar(64)  NOT NULL
);

INSERT INTO tags (slug, name) VALUES
                                  ('speaker','Спикер'), ('performer','Исполнитель'),
                                  ('organizer','Организатор'), ('guest','Гость'),
                                  ('jury','Жюри'), ('other','Другое');

CREATE TABLE event_people (
                              id         bigserial PRIMARY KEY,
                              event_id   bigint    NOT NULL REFERENCES events(id)   ON DELETE CASCADE,
                              person_id  bigint    NOT NULL REFERENCES persons(id)  ON DELETE CASCADE,
                              tag_id     smallint  NOT NULL REFERENCES tags(id),
                              sort_order smallint  NOT NULL DEFAULT 0,
                              UNIQUE (event_id, person_id, tag_id)
);

CREATE INDEX idx_people_slug       ON persons(slug);
CREATE INDEX idx_photos_event      ON event_photos(event_id);
CREATE INDEX idx_places_city       ON places(city);
CREATE INDEX idx_places_country    ON places(country);

CREATE INDEX idx_legal_info_event  ON legal_information(event_id);

CREATE INDEX idx_events_active     ON events(is_active);
CREATE INDEX idx_events_category   ON events(category_id);
CREATE INDEX idx_events_starts_at  ON events(starts_at);
CREATE INDEX idx_events_ends_at    ON events(ends_at);
CREATE INDEX idx_events_place      ON events(place_id);

CREATE TABLE IF NOT EXISTS export_cursors (
                                              consumer       text PRIMARY KEY,
                                              last_event_ts  timestamptz NOT NULL
);