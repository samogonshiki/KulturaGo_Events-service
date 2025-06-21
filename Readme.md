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

```
[Postgres] --batched-fetch--> [Exporter] --produce--> [Kafka topic: events.public]
     ^                                              |
     |                                              v
[Redis] <--persist cursor----------------------------+
```


## Таблица `events` БД

| Поле            | Тип                | NULL | По-умолчанию                | Назначение / пример                                |
|-----------------|--------------------|------|-----------------------------|----------------------------------------------------|
| `id`            | `uuid`             | ❌   | `uuid_generate_v4()`        | Уникальный идентификатор события                   |
| `category`      | `varchar(32)`      | ❌   | —                           | Тип мероприятия: `theatre`, `show`, `museum`, …    |
| `title`         | `varchar(256)`     | ❌   | —                           | Человекочитаемое имя («Концерт The Beatles»)       |
| `starts_at`     | `timestamptz`      | ❌   | —                           | Дата/время начала                                  |
| `ends_at`       | `timestamptz`      | ❌   | —                           | Дата/время окончания                               |
| `expires_at`    | `timestamptz`      | ❌   | `ends_at + interval '4 mon'`| Когда событие “протухает” и будет удалено          |
| `created_at`    | `timestamptz`      | ❌   | `now()`                     | Вставка в БД                                       |

> **TTL-правило:** событие живёт ≈ 3-4 месяца.  
> После даты `expires_at` запись удаляется фоновым джобом.

**Индексы**

| Имя                         | Колонка(и)        | Зачем                              |
|-----------------------------|-------------------|------------------------------------|
| `idx_events_starts_at`      | `starts_at`       | Быстрый вывод афиши по датам       |
| `idx_events_expires_at`     | `expires_at`      | Поиск кандидатов на удаление       |
| `idx_events_category`       | `category`        | Фильтрация (театр, шоу, …)         |

---

### Пример содержимого

| id                                   | category | title                    | starts_at               | ends_at                 | expires_at              |
|--------------------------------------|----------|--------------------------|-------------------------|-------------------------|-------------------------|
| `ec3d…`                              | concert  | Queen Tribute Show       | 2025-07-05 20:00 +02:00 | 2025-07-05 22:30 +02:00 | 2025-11-05 22:30 +01:00 |
| `a91d…`                              | theatre  | «Горе от ума»            | 2025-08-10 19:00 +02:00 | 2025-08-10 21:45 +02:00 | 2025-12-10 21:45 +01:00 |
| `b77e…`                              | olympiad | Школьная олимпиада по ИТ | 2025-09-15 10:00 +02:00 | 2025-09-15 16:00 +02:00 | 2026-01-15 16:00 +01:00 |

---

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