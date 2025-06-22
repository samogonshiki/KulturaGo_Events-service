# KulturaGo_Events-service


## Структура

```
events-service/
├── cmd/
│   └── events-service/
│       └── main.go
├── internal/
│   ├── config/
│   ├── domain/
│   ├── repository/
│   │   ├── postgres/
│   │   └── redis/
│   ├── usecase/
│   │   └── exporter.go
│   ├── broker/
│   │   └── kafka/
│   │       ├── producer.go
│   │       └── health.go
│   ├── transport/
│   │   └── http/
│   ├── scheduler/
│   ├── logger/
│   └── metrics/
├── pkg/
├── migrations/
├── deploy/
│   ├── docker/
│   └── k8s/
│       ├── deployment.yaml 
│       └── service.yaml
├── .env.example
├── Makefile
└── go.mod
```

---

## 🌐 Архитектурная диаграмма

```mermaid
flowchart LR
    PG[PostgreSQL]
    Exp[Exporter]
    Kafka[(Kafka topic: events.public)]
    Redis[(Redis)]

    PG -- batched-fetch --> Exp
    Exp -- produce --> Kafka
    Kafka -- store cursor --> Redis
    Redis -- last_event_ts --> Exp
```

## HTTP API


| Method | Path                                 | Description                                  | Request Body (JSON)          | Response Code  |
|:------:|:-------------------------------------|:---------------------------------------------|:-----------------------------|:---------------|
| GET    | `/api/v1/events?limit={limit}&offset={offset}` | Список активных событий (постранично)       | —                             | `200 OK`       |
| GET    | `/api/v1/events/{slug}`              | Детали события по его `slug`                | —                             | `200 OK`       |
| POST   | `/api/v1/events`                     | Создать новый объект события                | `CreateEventInput` (см. выше) | `201 Created`  |

**Query Parameters for GET `/api/v1/events`**
- `limit`: `int` (опционально) — сколько записей вернуть (по умолчанию 50).
- `offset`: `int` (опционально) — смещение для пагинации.

**Пример запроса**

```
{
  "category": {
    "name": "выставка"
  },
  "title": "Выставка искусств Севера",
  "description": "Приглашаем всех!",
  "place": {
    "address": "ул. Ленина, 1, Москва",
    "latitude": 55.751244,
    "longitude": 37.618423
  },
  "starts_at": "2025-11-01T10:00:00Z",
  "ends_at":   "2025-11-01T13:00:00Z",
  "photos": [
    {
      "url": "https://example.com/photos/vsosh-main.jpg",
      "alt_text": "Логотип выставки",
      "is_main": true
    },
    {
      "url": "https://example.com/photos/vsosh-participants.jpg",
      "alt_text": "все",
      "is_main": false
    }
  ],
  "legal_info": [],
  "people": [
    {
      "name": "Фетюков Никита",
      "tag": {
        "name": "Организатор"
      }
    }
  ]
}
```


## Таблица `event_categories`

| Поле  | Тип             | NULL | По умолчанию | Описание                         |
|-------|-----------------|------|--------------|----------------------------------|
| id    | `smallserial`   | ❌   | —            | Первичный ключ                   |
| slug  | `varchar(32)`   | ❌   | —            | Уникальный человекочитаемый slug |
| name  | `varchar(64)`   | ❌   | —            | Название категории               |

**Индексы и ограничения**
- `event_categories_pkey` (PK по `id`)
- `event_categories_slug_key` (UNIQUE по `slug`)

---

## Таблица `tags`

| Поле  | Тип             | NULL | По умолчанию | Описание                         |
|-------|-----------------|------|--------------|----------------------------------|
| id    | `smallserial`   | ❌   | —            | Первичный ключ                   |
| slug  | `varchar(32)`   | ❌   | —            | Уникальный человекочитаемый slug |
| name  | `varchar(64)`   | ❌   | —            | Название тега                    |

**Индексы и ограничения**
- `tags_pkey` (PK по `id`)
- `tags_slug_key` (UNIQUE по `slug`)

---

## Таблица `places`

| Поле         | Тип                | NULL | По умолчанию | Описание               |
|--------------|--------------------|------|--------------|------------------------|
| id           | `bigserial`        | ❌   | —            | Первичный ключ         |
| title        | `varchar(128)`     | ❌   | —            | Название/заголовок     |
| country      | `varchar(64)`      | ✅   | —            | Страна                 |
| region       | `varchar(64)`      | ✅   | —            | Регион/область         |
| city         | `varchar(128)`     | ❌   | —            | Город                  |
| street       | `varchar(256)`     | ✅   | —            | Улица                  |
| house_num    | `varchar(16)`      | ✅   | —            | Номер дома             |
| postal_code  | `varchar(16)`      | ✅   | —            | Почтовый индекс        |
| longitude    | `numeric(9,6)`     | ✅   | —            | Долгота                |
| latitude     | `numeric(9,6)`     | ✅   | —            | Широта                 |
| created_at   | `timestamptz`      | ❌   | `now()`      | Дата создания записи   |

**Индексы**
- `idx_places_city` (`city`)
- `idx_places_country` (`country`)

---

## Таблица `events`

| Поле         | Тип                | NULL | По умолчанию        | Описание                                         |
|--------------|--------------------|------|---------------------|--------------------------------------------------|
| id           | `bigserial`        | ❌   | —                   | Первичный ключ                                   |
| slug         | `varchar(128)`     | ❌   | —                   | Уникальный человекочитаемый идентификатор события|
| category_id  | `smallint`         | ❌   | —                   | FK → `event_categories(id)`                      |
| title        | `varchar(64)`      | ❌   | —                   | Название события                                 |
| description  | `varchar(4096)`    | ❌   | —                   | Описание события                                 |
| place_id     | `bigint`           | ❌   | —                   | FK → `places(id)` (локация)                      |
| starts_at    | `timestamptz`      | ❌   | —                   | Время начала                                     |
| ends_at      | `timestamptz`      | ❌   | —                   | Время окончания                                  |
| is_active    | `boolean`          | ❌   | `TRUE`              | Флаг активности                                  |
| created_at   | `timestamptz`      | ❌   | `now()`             | Время создания записи                            |
| updated_at   | `timestamptz`      | ❌   | `now()`             | Время последнего обновления (триггер)            |

**Индексы и триггер**
- `events_pkey` (PK по `id`)
- `events_slug_key` (UNIQUE по `slug`)
- `idx_events_active` (`is_active`)
- `idx_events_category` (`category_id`)
- `idx_events_starts_at` (`starts_at`)
- `idx_events_ends_at` (`ends_at`)
- `idx_events_place` (`place_id`)
- Триггер `events_set_updated_at` обновляет `updated_at` перед каждым `UPDATE events`.

---

## Таблица `legal_information`

| Поле      | Тип             | NULL | По умолчанию | Описание                             |
|-----------|-----------------|------|--------------|--------------------------------------|
| id        | `bigserial`     | ❌   | —            | Первичный ключ                       |
| event_id  | `bigint`        | ❌   | —            | FK → `events(id)`                    |
| info_key  | `varchar(64)`   | ❌   | —            | Ключ информационного поля            |
| info_text | `text`          | ❌   | —            | Текст правовой/дополнительной информации |

**Индексы**
- `idx_legal_info_event` (`event_id`)

---

## Таблица `event_photos`

| Поле       | Тип             | NULL | По умолчанию | Описание                  |
|------------|-----------------|------|--------------|---------------------------|
| id         | `bigserial`     | ❌   | —            | Первичный ключ            |
| event_id   | `bigint`        | ❌   | —            | FK → `events(id)`         |
| url        | `text`          | ❌   | —            | Ссылка на фотографию      |
| alt_text   | `varchar(256)`  | ✅   | —            | Альтернативный текст      |
| is_main    | `boolean`       | ❌   | `FALSE`      | Флаг “главной” фотографии |
| created_at | `timestamptz`   | ❌   | `now()`      | Время добавления          |

**Индексы**
- `idx_photos_event` (`event_id`)

---

## Таблица `persons`

| Поле        | Тип              | NULL | По умолчанию | Описание                      |
|-------------|------------------|------|--------------|-------------------------------|
| id          | `bigserial`      | ❌   | —            | Первичный ключ                |
| slug        | `varchar(128)`   | ❌   | —            | Уникальный человекочитаемый slug |
| name        | `varchar(256)`   | ❌   | —            | Имя/ФИО человека              |
| description | `text`           | ✅   | —            | Описание персоны              |
| photo       | `text`           | ❌   | —            | Ссылка на фото                |
| created_at  | `timestamptz`    | ❌   | `now()`      | Время создания записи         |

**Индексы**
- `idx_people_slug` (`slug`)

---

## Таблица `event_people`

| Поле       | Тип           | NULL | По умолчанию | Описание                                 |
|------------|---------------|------|--------------|------------------------------------------|
| id         | `bigserial`   | ❌   | —            | Первичный ключ                           |
| event_id   | `bigint`      | ❌   | —            | FK → `events(id)`                        |
| person_id  | `bigint`      | ❌   | —            | FK → `persons(id)`                       |
| tag_id     | `smallint`    | ❌   | —            | FK → `tags(id)`                          |
| sort_order | `smallint`    | ❌   | `0`          | Порядок сортировки (для вывода людей)    |

**Ограничения**
- `UNIQUE (event_id, person_id, tag_id)`

---

## Таблица `export_cursors`

| Поле           | Тип            | NULL | По умолчанию | Описание                                |
|----------------|----------------|------|--------------|-----------------------------------------|
| consumer       | `text`         | ❌   | —            | Имя потребителя                        |
| last_event_ts  | `timestamptz`  | ❌   | —            | Время последнего экспортированного события |

## Команды для makefile

- локально (development)
```shell
make migrate-up
```
- откатить на один шаг
```shell
make migrate-down
```


- выполнить миграции в контейнере stage/prod
```shell
GOOSE_ENV=production DATABASE_URL="postgres://user:pass@postgres:5432/prod?sslmode=disable" \
make migrate-up
```

**by Finnik**