lint:
	golangci-lint run
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

env-clean:
	@echo "🧹 Cleaning up environment files..."
	@if exist "deployments\.env" ( \
		del "deployments\.env" && echo "🗑️  Removed deployments\.env" \
	)
	@if exist "deployments\dev\.env" ( \
		del "deployments\dev\.env" && echo "🗑️  Removed deployments\dev\.env" \
	)
	@echo "✅ Environment cleanup complete!"



# Infrastructure (Development)
dev-up:
	docker-compose -f deployments/docker-compose.yml -f deployments/dev/docker-compose.dev.yml up -d

dev-down:
	docker-compose -f deployments/docker-compose.yml -f deployments/dev/docker-compose.dev.yml down