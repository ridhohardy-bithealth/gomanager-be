# Include .env file
include .env

# Export all variables
export

# Database URL construction
DB_URL=${DATABASE_URL}

# Default number of migrations to up/down
N ?= 1

## Migration Commands
.PHONY: server migrations-create migrations-up migrations-down migrations-force migrations-version migrations-init

# Run Server
server:
	go run cmd/api/main.go
# Create a new migration file
migrations-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./db/migrations -seq $${name}

# Apply all or N migrations
migrations-up:
	migrate -path ./db/migrations -database "$(DB_URL)" up $(N)

# Rollback all or N migrations
migrations-down:
	migrate -path ./db/migrations -database "$(DB_URL)" down $(N)

# Force set version
migrations-force:
	@read -p "Enter version to force: " version; \
	migrate -path ./db/migrations -database "$(DB_URL)" force $${version}

# Show current migration version
migrations-version:
	migrate -path ./db/migrations -database "$(DB_URL)" version

# Drop everything
migrations-drop:
	migrate -path ./db/migrations -database "$(DB_URL)" drop

# Create initial migration
migrations-init:
	migrate create -ext sql -dir ./db/migrations -seq init_schema