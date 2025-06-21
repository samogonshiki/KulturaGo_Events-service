CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE event_category AS ENUM (
    'theatre',
    'show',
    'museum',
    'concert',
    'contest',
    'olympiad',
    'other'
);

CREATE TABLE events (
                        id          uuid            PRIMARY KEY DEFAULT uuid_generate_v4(),
                        category    event_category  NOT NULL,
                        title       varchar(256)    NOT NULL,

                        starts_at   timestamptz     NOT NULL,
                        ends_at     timestamptz     NOT NULL,

                        expires_at  timestamptz     GENERATED ALWAYS AS (ends_at + interval '4 months') STORED,

    created_at  timestamptz     NOT NULL DEFAULT now()
);

CREATE INDEX idx_events_starts_at   ON events (starts_at);
CREATE INDEX idx_events_expires_at  ON events (expires_at);
CREATE INDEX idx_events_category    ON events (category);


CREATE TABLE IF NOT EXISTS export_cursors (
                                              consumer       text         PRIMARY KEY,
                                              last_event_ts  timestamptz  NOT NULL
);