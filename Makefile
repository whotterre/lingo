migup:
	@echo "Running migrations..."
	@if [ -z "$DB_URL" ]; then \
		echo "Error: DB_URL environment variable is not set."; \
		exit 1; \
	fi
	cd ./internal/db/migrations && migrate -path . -database "postgresql://postgres:password@172.17.0.2:5432/lingo_db?sslmode=disable" up
migdown:
	@echo "Rolling back migrations..."
	@if [ -z "$DB_URL" ]; then \
		echo "Error: DB_URL environment variable is not set."; \
		exit 1; \
	fi
	cd ./internal/db/migrations && migrate -path . -database "postgresql://postgres:password@172.17.0.2:5432/lingo_db?sslmode=disable" down
sqlc_gen:
	@echo "Generating SQL code..."
	sqlc generate
run:
	@echo "Running application..."
	cd ./cmd/api && go run .
