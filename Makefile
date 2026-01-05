include .env.dev
export

# Get the current user's ID and Group ID from the host system
export UID := $(shell id -u)
export GID := $(shell id -g)

DOCKER_DEV=docker compose -f docker-compose.yml -f docker-compose.dev.yml --env-file .env.dev
MIGRATE_CONFIG=-path migrations -database $$DB_URL

dev:
	$(DOCKER_DEV) up -d --build
down:
	$(DOCKER_DEV) down
shell:
	$(DOCKER_DEV) exec app-dev sh
test:
	docker compose -f docker-compose.yml -f docker-compose.test.yml --env-file .env.test up --build --exit-code-from app-test
db-shell:
	$(DOCKER_DEV) exec -it db psql -U $(DB_USER) -d $(DB_NAME)
# Run database migrations
db-up:
	$(DOCKER_DEV) exec app-dev sh -c "migrate -path migrations -database \$$DB_URL up"
# Rollback one database migration
db-down:
	$(DOCKER_DEV) exec app-dev sh -c "migrate -path migrations -database \$$DB_URL down"
db-drop:
	$(DOCKER_DEV) exec app-dev sh -c "migrate -path migrations -database \$$DB_URL drop"
db-migration:
	$(DOCKER_DEV) exec app-dev migrate create -ext sql -dir migrations $(name)
debug:
	$(DOCKER_DEV) exec app-dev sh -c "echo \$$UID"
