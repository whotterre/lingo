migrate:
	@echo "Running migrations..."
	migrate -path db/migrations -database "$DB_URL" up
sqlc_gen:
	@echo "Generating SQL code..."
	sqlc generate

