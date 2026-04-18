MIGRATIONS_DIR := ./internal/database/migrations

# Optional: provides DATABASE_URL for migrate-up / migrate-down (same as the app uses).
-include .env
export DATABASE_URL

# Extra goals after "migrate-create" (e.g. make migrate-create add_items_table).
# First extra word is the migration name; any further words are no-op phony targets.
EXTRA_TARGETS := $(wordlist 2,999,$(MAKECMDGOALS))
MIGRATION_NAME := $(firstword $(EXTRA_TARGETS))

.PHONY: migrate-create migrate-up migrate-down $(EXTRA_TARGETS) run

ifneq ($(strip $(EXTRA_TARGETS)),)
$(EXTRA_TARGETS):
	@:
endif

# Usage: make migrate-create your_migration_name
# make migrate-create NAME=your_migration_name
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(if $(strip $(MIGRATION_NAME)),$(MIGRATION_NAME),$(if $(strip $(NAME)),$(NAME),$(error Migration name required. Usage: make migrate-create your_migration_name)))

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(if $(DATABASE_URL),$(DATABASE_URL),$(error DATABASE_URL is not set. Add it to .env beside this Makefile or export it.))" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(if $(DATABASE_URL),$(DATABASE_URL),$(error DATABASE_URL is not set. Add it to .env beside this Makefile or export it.))" down

run:
	go run ./cmd/api .