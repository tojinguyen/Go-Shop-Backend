lint:
	golangci-lint run

build:
	docker compose up -d --build 
up:
	docker compose up -d 
down:
	docker compose down

seed-users:
	@echo "🌱 Seeding user-service database..."
	@go run ./internal/services/user-service/cmd/seeder/main.go -count 20000