DB_URL ?= mysql://root:secret@tcp(localhost:3306)/coffee_pos
MIGRATIONS_PATH ?= migrations

.PHONY: run tidy migrate-up migrate-down migrate-down-all migrate-version migrate-create

run:
	go run cmd/api/main.go

tidy:
	go mod tidy

migrate-up:
	@output=$$(migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up 2>&1) || status=$$?; \
	if [ "$${status:-0}" -ne 0 ]; then \
		if echo "$$output" | grep -q "no change"; then \
			echo "Database is already up to date (no change)."; \
		elif echo "$$output" | grep -q "file does not exist"; then \
			echo "No migration files found in $(MIGRATIONS_PATH)/. Create one using 'make migrate-create name=<name>'."; \
		else \
			echo "$$output"; \
			exit $$status; \
		fi \
	else \
		echo "$$output"; \
	fi

migrate-down:
	@output=$$(migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1 2>&1) || status=$$?; \
	if [ "$${status:-0}" -ne 0 ]; then \
		if echo "$$output" | grep -q "no change"; then \
			echo "No migrations to roll back (no change)."; \
		else \
			echo "$$output"; \
			exit $$status; \
		fi \
	else \
		echo "$$output"; \
	fi

migrate-down-all:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down -all

migrate-version:
	@output=$$(migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version 2>&1) || status=$$?; \
	if [ "$${status:-0}" -ne 0 ]; then \
		if echo "$$output" | grep -q "no migration"; then \
			echo "No migrations applied yet (version: 0)."; \
		else \
			echo "$$output"; \
			exit $$status; \
		fi \
	else \
		echo "Current migration version: $$output"; \
	fi

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: variable 'name' is required. Example: make migrate-create name=create_users_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
