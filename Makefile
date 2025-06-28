# =======================================================
# GO-SHOP MAKEFILE - Updated for new structure
# =======================================================

.PHONY: help dev-up dev-down prod-up prod-down clean logs build test env-setup

# Default target
help:
	@echo "Available commands:"
	@echo "  dev-up      - Start development environment"
	@echo "  dev-down    - Stop development environment"
	@echo "  prod-up     - Start production environment"
	@echo "  prod-down   - Stop production environment"
	@echo "  clean       - Remove all containers, volumes, and networks"
	@echo "  logs        - Show logs for all services"
	@echo "  build       - Build all services"
	@echo "  test        - Run tests for all services"
	@echo "  env-setup   - Setup environment files"

# Development environment
dev-up:
	@echo "🚀 Starting development environment..."
	@copy .env.development .env
	@docker-compose -f deployments/docker-compose.yml -f deployments/dev/docker-compose.dev.yml up -d

dev-down:
	@echo "⏹️  Stopping development environment..."
	@docker-compose -f deployments/docker-compose.yml -f deployments/dev/docker-compose.dev.yml down

# Production environment  
prod-up:
	@echo "🚀 Starting production environment..."
	@copy .env.production .env
	@docker-compose -f deployments/docker-compose.yml -f deployments/prod/docker-compose.prod.yml up -d

prod-down:
	@echo "⏹️  Stopping production environment..."
	@docker-compose -f deployments/docker-compose.yml -f deployments/prod/docker-compose.prod.yml down

# Clean everything
clean:
	@echo "🧹 Cleaning up containers, volumes, and networks..."
	@docker-compose -f deployments/docker-compose.yml -f deployments/dev/docker-compose.dev.yml down -v --remove-orphans
	@docker-compose -f deployments/docker-compose.yml -f deployments/prod/docker-compose.prod.yml down -v --remove-orphans
	@docker system prune -f

# Show logs
logs:
	@docker-compose -f deployments/docker-compose.yml logs -f

# Build all services
build:
	@echo "� Building all services..."
	@docker-compose -f deployments/docker-compose.yml build

# Run tests
test:
	@echo "🧪 Running tests for all services..."
	@docker-compose -f deployments/docker-compose.yml exec user-service go test ./... || echo "❌ user-service tests failed"
	@docker-compose -f deployments/docker-compose.yml exec shop-service go test ./... || echo "❌ shop-service tests failed"
	@docker-compose -f deployments/docker-compose.yml exec order-service go test ./... || echo "❌ order-service tests failed"
	@docker-compose -f deployments/docker-compose.yml exec payment-service go test ./... || echo "❌ payment-service tests failed"
	@docker-compose -f deployments/docker-compose.yml exec shipping-service go test ./... || echo "❌ shipping-service tests failed"
	@docker-compose -f deployments/docker-compose.yml exec review-service go test ./... || echo "❌ review-service tests failed"
	@docker-compose -f deployments/docker-compose.yml exec notification-service go test ./... || echo "❌ notification-service tests failed"

# Setup environment files
env-setup:
	@echo "📝 Setting up environment files..."
	@if not exist ".env" copy ".env.development" ".env" && echo "✅ Created .env from .env.development"
	@echo "🎯 Environment setup complete!"

# Legacy commands for backward compatibility
dev:
	go run ./cmd/server/main.go

rand-user:
	go run ./cmd/seed/main.go

lint:
	golangci-lint run

docker build:
	docker-compose up -d --build

docker up:
	docker-compose up -d

docker down:
	docker-compose down

create-migration:
	if "$(name)" == "" ( \
		echo ❌ Thiếu tên migration. Dùng: make create-migration name=ten_migration & exit /b 1 \
	) else ( \
		goose -dir migrations/mysql create -s $(name) sql \
	)

# Infrastructure (Main)
up:
	docker-compose -f deployments/docker-compose.yml up -d

down:
	docker-compose -f deployments/docker-compose.yml down