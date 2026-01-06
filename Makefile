# Get the current user's ID and Group ID from the host system
export UID := $(shell id -u)
export GID := $(shell id -g)

DOCKER_DEV=docker compose -f docker-compose.yml -f docker-compose.dev.yml --env-file .env.dev
DOCKER_TEST=docker compose -f docker-compose.yml -f docker-compose.test.yml --env-file .env.test
MIGRATE_CONFIG=-path migrations -database $$DB_URL

dev:
	$(DOCKER_DEV) up --build
down:
	$(DOCKER_DEV) down
shell:
	$(DOCKER_DEV) exec app-dev sh
test:
	$(DOCKER_TEST) up --build --exit-code-from app-test
db-shell:
	@export $$(cat .env.dev | xargs) && $(DOCKER_DEV) exec -it db psql -U $$DB_USER -d $$DB_NAME
# Run database migrations
db-up-dev:
	$(DOCKER_DEV) exec app-dev sh -c "migrate -path migrations -database \$$DB_URL up"
db-up-test:
	$(DOCKER_TEST) run --rm app-test sh -c "migrate -path migrations -database \$$DB_URL up"
# Rollback one database migration
db-down:
	$(DOCKER_DEV) exec app-dev sh -c "migrate -path migrations -database \$$DB_URL down"
db-drop:
	$(DOCKER_DEV) exec app-dev sh -c "migrate -path migrations -database \$$DB_URL drop"
db-migration:
	$(DOCKER_DEV) exec app-dev migrate create -ext sql -dir migrations $(name)
sqlc-generate:
	$(DOCKER_DEV) exec app-dev go tool sqlc generate
