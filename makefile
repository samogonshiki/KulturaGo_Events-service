APP            := events-service
BIN            := bin/$(APP)
IMAGE          := $(APP)
GOFLAGS        := -trimpath -ldflags="-s -w"
SWAG           := $(GOPATH)/bin/swag
MIGRATE        := $(GOPATH)/bin/migrate
DATABASE_URL  ?= postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
GOOSE         ?= $(GOPATH)/bin/goose
GOOSE_DIR     := ./db/migrations
GOOSE_ENV     ?= development

.PHONY: all build swagger migrate-up migrate-down test stage prod clean

all: build

clean:
	rm -rf bin

swagger: $(SWAG)
	$(SWAG) init -g cmd/events-service/main.go -o api/docs --parseDependency

$(SWAG):
	go install github.com/swaggo/swag/cmd/swag@latest

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BIN) ./cmd/events-service

migrate-up: $(MIGRATE)
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" up

migrate-down: $(MIGRATE)
	$(MIGRATE) -path migrations -database "$(DATABASE_URL)" down 1

$(MIGRATE):
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

$(GOOSE):
	go install github.com/pressly/goose/v3/cmd/goose@latest

migrate-up: $(GOOSE)
	$(GOOSE) -dir $(GOOSE_DIR) -env $(GOOSE_ENV) up

migrate-down: $(GOOSE)
	$(GOOSE) -dir $(GOOSE_DIR) -env $(GOOSE_ENV) down


dc-up:
	docker-compose up -d --build

dc-dw:
	docker-compose down

test: export LOG_LEVEL=debug
test: build
	go test ./... -v

stage: export LOG_LEVEL=error
stage: build
	docker build -t $(IMAGE):stage .

prod: export LOG_LEVEL=error
prod: build
	docker build -t $(IMAGE):latest .