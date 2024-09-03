.PHONY: build run init tidy migrate migration-create

run:
	go run cmd/app/main.go

init:
	go mod init pharmacy-backend

tidy:
	@echo "Running go mod tidy..."
	go mod tidy

migrate:
	go run cmd/migrate/main.go

migration-create:
	 @if [ -z "$(name)" ]; then \
		echo "Error: Migration name is required. Usage: make create-migration name=your_migration_name"; \
		exit 1; \
	 fi
	 migrate create -ext sql -dir ./migrations -seq $(name)