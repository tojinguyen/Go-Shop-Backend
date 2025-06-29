lint:
	golangci-lint run
create-migration:
	if "$(name)" == "" ( \
		echo ❌ Thiếu tên migration. Dùng: make create-migration name=ten_migration & exit /b 1 \
	) else ( \
		goose -dir migrations/mysql create -s $(name) sql \
	)

build:
	docker compose up -d --build 
up:
	docker compose up -d 
down:
	docker compose down

