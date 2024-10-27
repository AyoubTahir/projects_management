# Makefile

# Set the database URL
DATABASE_URL=postgres://postgres:root@localhost:5432/testapi?sslmode=disable

# Define migration path
MIGRATION_PATH=migrations

# Command to apply migrations
start:
	go run cmd/api/main.go
