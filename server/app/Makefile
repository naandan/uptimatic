DB_URL=postgres://root:forward098@localhost:5432/uptimatic?sslmode=disable
MIGRATIONS_PATH=./migrations

# Apply semua migration
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

# Rollback 1 step
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

# Rollback semua
migrate-rollback:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down -all

# Re-run dari awal (reset + up)
migrate-reset: migrate-reset migrate-up

# Buat file migration baru (contoh: make migrate-new name=create_table_xyz)
migrate-new:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)


# Start Server
server:
	go run main.go server

# Start Worker
worker:
	go run main.go worker