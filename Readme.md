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