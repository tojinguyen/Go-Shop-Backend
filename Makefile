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

# Environment Setup Commands
env-setup-all:
	@echo "🚀 Setting up all environment files..."
	@if not exist "deployments\.env" ( \
		copy "deployments\.env.example" "deployments\.env" && echo "✅ Created deployments\.env" \
	) else ( \
		echo "⚠️  deployments\.env already exists" \
	)
	@if not exist "deployments\dev\.env" ( \
		copy "deployments\dev\.env.example" "deployments\dev\.env" && echo "✅ Created deployments\dev\.env" \
	) else ( \
		echo "⚠️  deployments\dev\.env already exists" \
	)
	@if not exist "deployments\prod\.env" ( \
		copy "deployments\prod\.env.example" "deployments\prod\.env" && echo "✅ Created deployments\prod\.env" \
	) else ( \
		echo "⚠️  deployments\prod\.env already exists" \
	)
	@echo "🎯 Environment setup complete!"

env-setup-base:
	@echo "🚀 Setting up base environment file..."
	@if not exist "deployments\.env" ( \
		copy "deployments\.env.example" "deployments\.env" && echo "✅ Created deployments\.env" \
	) else ( \
		echo "⚠️  deployments\.env already exists" \
	)

env-setup-dev:
	@echo "🚀 Setting up development environment file..."
	@if not exist "deployments\dev\.env" ( \
		copy "deployments\dev\.env.example" "deployments\dev\.env" && echo "✅ Created deployments\dev\.env" \
	) else ( \
		echo "⚠️  deployments\dev\.env already exists" \
	)

env-setup-prod:
	@echo "🚀 Setting up production environment file..."
	@if not exist "deployments\prod\.env" ( \
		copy "deployments\prod\.env.example" "deployments\prod\.env" && echo "✅ Created deployments\prod\.env" \
	) else ( \
		echo "⚠️  deployments\prod\.env already exists" \
	)

env-clean:
	@echo "🧹 Cleaning up environment files..."
	@if exist "deployments\.env" ( \
		del "deployments\.env" && echo "🗑️  Removed deployments\.env" \
	)
	@if exist "deployments\dev\.env" ( \
		del "deployments\dev\.env" && echo "🗑️  Removed deployments\dev\.env" \
	)
	@if exist "deployments\prod\.env" ( \
		del "deployments\prod\.env" && echo "🗑️  Removed deployments\prod\.env" \
	)
	@echo "✅ Environment cleanup complete!"

# Infrastructure (Main)
infra-up:
	docker-compose -f deployments/docker-compose.yml --env-file deployments/.env up -d

infra-down:
	docker-compose -f deployments/docker-compose.yml --env-file deployments/.env down

# Infrastructure (Production)
infra-prod-up:
	docker-compose -f deployments/docker-compose.yml -f deployments/prod/docker-compose.prod.yml --env-file deployments/prod/.env up -d

infra-prod-down:
	docker-compose -f deployments/docker-compose.yml -f deployments/prod/docker-compose.prod.yml --env-file deployments/prod/.env down

# Infrastructure (Development)
infra-dev-up:
	docker-compose -f deployments/docker-compose.yml -f deployments/dev/docker-compose.dev.yml --env-file deployments/dev/.env up -d

infra-dev-down:
	docker-compose -f deployments/docker-compose.yml -f deployments/dev/docker-compose.dev.yml --env-file deployments/dev/.env down