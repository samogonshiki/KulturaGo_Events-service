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


**by Finnik**