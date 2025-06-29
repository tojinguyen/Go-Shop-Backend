lint:
	golangci-lint run
create-migration:
	if "$(name)" == "" ( \
		echo ❌ Thiếu tên migration. Dùng: make create-migration name=ten_migration & exit /b 1 \
	) else ( \
		goose -dir migrations/mysql create -s $(name) sql \
	)

	@echo "🧹 Cleaning up environment files..."
	@if exist "deployments\.env" ( \
		del "deployments\.env" && echo "🗑️  Removed deployments\.env" \
	)
	@if exist "deployments\dev\.env" ( \
		del "deployments\dev\.env" && echo "🗑️  Removed deployments\dev\.env" \
	)
	@echo "✅ Environment cleanup complete!"


