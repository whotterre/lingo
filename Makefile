migrate:
	@echo "Running migrations..."
	@if [ -z "$DB_URL" ]; then \
		echo "Error: DB_URL environment variable is not set."; \
		exit 1; \
	fi
	cd ./internal/db/migrations && migrate -path . -database $DB_URL up
sqlc_gen:
	@echo "Generating SQL code..."
	sqlc generate
run:
	@echo "Running application..."
	cd ./cmd/api && go run .
