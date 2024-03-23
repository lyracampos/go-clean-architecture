PG_CONNECTION_STRING ?= 'postgres://user:password@localhost:5433/go_clean_arch?sslmode=disable'

MIGRATE := go run -tags='postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0
MIGRATIONS_PATH ?= './internal/gateways/postgres/migrations'

.PHONY: build
build:
	go build cmd/main.go

.PHONY: deps/start
deps/start:
	docker compose up -d
	until docker exec postgres pg_isready; do echo 'Waiting for postgres server...' && sleep 1; done

.PHONY: deps/stop
deps/stop:
	docker compose down

.PHONY: api/start
api/start:
	go run cmd/main.go -e api -c ./config/config.yaml

.PHONY: worker/start
worker/start:
	go run cmd/main.go -e worker -c ./config/config.yaml

.PHONY: migration/create
migration/create:
	$(MIGRATE) create -seq -ext sql -dir $(MIGRATIONS_PATH) $(MIGRATION_NAME)

.PHONY: migration/up
migration/up:
	$(MIGRATE) -path $(MIGRATIONS_PATH) --database $(PG_CONNECTION_STRING) up
