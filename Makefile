include .env
export

env = local

.env:
	@ln -sf ./configs/.env.${env} .env

.PHONY: postgres
postgres:
	@docker-compose up -d postgres
	@until docker-compose exec -T postgres pg_isready > /dev/null 2>&1; do echo "Waiting postgres"; done

.PHONY: deps
deps: postgres

.PHONY: psql
psql:
	@docker-compose exec postgres psql -U ${POSTGRES_USER} ${POSTGRES_DB}

.PHONY: dev
dev:
	@reflex -s -- go run ./cmd/go-url-shortener/

.PHONY: build
build:
	@go build  -o ./bin/ ./cmd/...

.PHONY: test
test:
	@go test -v -cover ./...

.PHONY: lint
lint:
	@staticcheck ./...