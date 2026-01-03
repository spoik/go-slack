# include .env.dev
# export

MIGRATE_CMD=cmd/migrations/main.go
DOCKER_DEV=docker compose -f docker-compose.yml -f docker-compose.dev.yml --env-file .env.dev

dev:
	$(DOCKER_DEV) up --build
down:
	$(DOCKER_DEV) down
test:
	docker compose -f docker-compose.yml -f docker-compose.test.yml --env-file .env.test up --build --exit-code-from app-test
db-shell:
	$(DOCKER_DEV) exec -it db psql -U $(DB_USER) -d $(DB_NAME)
# Initialize database migration tables
db-init:
	$(DOCKER_DEV) run --rm app-dev go run $(MIGRATE_CMD) init
# Run database migrations
db-up:
	$(DOCKER_DEV) run --rm app-dev go run $(MIGRATE_CMD) up
# Rollback one database migration
db-down:
	$(DOCKER_DEV) run --rm app-dev go run $(MIGRATE_CMD) down
# Check database migration status
db-status:
	$(DOCKER_DEV) exec -it app-dev go run $(MIGRATE_CMD) status
db-migration:
	$(DOCKER_DEV) run --rm app-dev go run $(MIGRATE_CMD) new $(name)

