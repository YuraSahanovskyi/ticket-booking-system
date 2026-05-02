DB_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
MIGRATIONS_PATH=internal/db/migrations/

.PHONY: migrate-up migrate-down migrate-force migrate-create

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate-force:
	@read -p "Enter version: " v; \
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $$v

db-schema:
	docker exec -t postgres pg_dump -s -U postgres postgres > schema.sql