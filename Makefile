PG_CONNECTION_STRING ?= 'postgres://user:password@localhost:5432/go_clean_arch?sslmode=disable'
MIGRATE := go run -tags='postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0
MIGRATIONS_PATH ?= './internal/gateways/postgres/migrations'

GATEWAY_PORTS_PATH=./internal/domain/ports
GATEWAY_PORTS_MOCKS_PATH=./internal/domain/ports/mocks

GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.1
SWAGGER := docker run --rm -e GOPATH=$$(go env GOPATH):/go -v $$(pwd):$$(pwd) -w $$(pwd) quay.io/goswagger/swagger:v0.30.4
MOCKGEN := go run go.uber.org/mock/mockgen@v0.4.0

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run --fix

.PHONY: build
build:
	go build cmd/main.go

.PHONY: deps/start
deps/start:
	docker compose up -d
	until docker exec postgres pg_isready; do echo 'Waiting for postgres server...' && sleep 1; done
	make migration/up

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

.PHONY: swagger/generate
swagger/generate:
	$(GOSWAGGER) generate spec -o ./swagger.yaml --scan-models

mocks/generate:
	$(MOCKGEN) -source=${GATEWAY_PORTS_PATH}/user_gateways.go \
			   -destination=${GATEWAY_PORTS_MOCKS_PATH}/user_gateways_mock.go \
			   -package=mock

test/run:
	go test ./... -cover