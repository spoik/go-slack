# Get the current user's ID and Group ID from the host system
export UID := $(shell id -u)
export GID := $(shell id -g)

DOCKER_DEV=docker compose -p go-slack-dev -f docker-compose.yml -f docker-compose.dev.yml --env-file .env.dev
DOCKER_TEST=docker compose -p go-slack-test -f docker-compose.yml -f docker-compose.test.yml --env-file .env.test
MIGRATE_CONFIG=-path migrations -database $$DB_URL

.PHONY: help
.DEFAULT_GOAL := help

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# --- Development environment ---
dev: ## Start development docker environment
	$(DOCKER_DEV) up --build
down: ## Stop development docker environment
	$(DOCKER_DEV) down
shell:
	$(DOCKER_DEV) exec app-dev sh
db-shell:
	@export $$(cat .env.dev | xargs) && $(DOCKER_DEV) exec -it db psql -U $$DB_USER -d $$DB_NAME
db-up-dev: # Run database migrations
	$(DOCKER_DEV) exec app-dev sh -c "migrate -path migrations -database \$$DB_URL up"
db-down-dev: ## Rollback one database migration in the development environment
	$(DOCKER_DEV) exec app-dev sh -c "migrate -path migrations -database \$$DB_URL down"
db-drop-dev: ## Drop database tables in development environment
	$(DOCKER_DEV) exec app-dev sh -c "migrate -path migrations -database \$$DB_URL drop"
db-migration: ## Create a new database migration in the development environment
	$(DOCKER_DEV) exec app-dev migrate create -ext sql -dir migrations $(name)
sqlc-generate: ## Generate sqlc files
	$(DOCKER_DEV) exec app-dev go tool sqlc generate

# --- Testing environment ---
test-up:
	$(DOCKER_TEST) up --build -d
test-down:
	$(DOCKER_TEST) down
test-app:
	$(DOCKER_TEST) exec -T app-test go test ./...
test-frontend:
	$(DOCKER_TEST) exec -T frontend-test npm run test
test-all: test-app test-frontend
db-up-test: ## Run database migrations in the testing environment
	$(DOCKER_TEST) exec app-test sh -c "migrate -path migrations -database \$$DB_URL up"
db-drop-test: ## Drop database tables in the testing environment
	$(DOCKER_TEST) exec app-dev sh -c "migrate -path migrations -database \$$DB_URL drop"
