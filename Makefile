.PHONY: build run

build:
	@go build -o build/ecom-go cmd/main.go

run: build
	@./build/ecom-go

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

%:
	@:

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

migrate-force:
	@go run cmd/migrate/main.go force $(filter-out $@,$(MAKECMDGOALS))
