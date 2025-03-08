.PHONY: setup-env .env

include .env

DB_PORT ?= 5432
BINARY=arcane.out

build:
	go build -o $(BINARY) .

up-db: down-db
	@podman run --name=$(DB_CONTAINER_NAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) -itd -p 5432:$(DB_PORT) postgres:latest
	@echo "Waiting 3 seconds before running migrations"
	@echo 3
	@sleep 1
	@echo 2
	@sleep 1
	@echo 1
	@sleep 1
	$(MAKE) migrate
	$(MAKE) seed

down-db:
	@podman stop $(DB_CONTAINER_NAME) && podman rm $(DB_CONTAINER_NAME) || true

reset-db:
	@podman exec -it $(DB_CONTAINER_NAME) psql -U postgres -d $(DB_NAME) -c "DROP SCHEMA public CASCADE;"
	@podman exec -it $(DB_CONTAINER_NAME) psql -U postgres -d $(DB_NAME) -c "CREATE SCHEMA public;"
	$(MAKE) migrate

into-db:
	@podman exec -it $(DB_CONTAINER_NAME) bash

migrate:
	@podman cp ./schema.sql $(DB_CONTAINER_NAME):/tmp/schema.sql
	@podman exec -it $(DB_CONTAINER_NAME) psql -U postgres -d $(DB_NAME) -f /tmp/schema.sql

seed:
	@podman cp ./seed.sql $(DB_CONTAINER_NAME):/tmp/seed.sql
	@podman exec -it $(DB_CONTAINER_NAME) psql -U postgres -d $(DB_NAME) -f /tmp/seed.sql

