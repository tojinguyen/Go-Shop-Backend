lint:
	golangci-lint run

build:
	docker compose up -d --build 
up:
	docker compose up -d 
down:
	docker compose down

.PHONY: seed-users
seed-users:
	@echo "🌱 Seeding user-service database with 50000 customers and 1000 shippers..."
	@go run ./internal/services/user-service/cmd/seeder/main.go -users=50000 -shippers=1000