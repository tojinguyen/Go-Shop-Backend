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
infra-up:
	docker-compose -f deployments/docker-compose.yml up -d

infra-down:
	docker-compose -f deployments/docker-compose.yml down

# Infrastructure (Production)
infra-prod-up:
	docker-compose -f deployments/prod/docker-compose.prod.yml up -d

infra-prod-down:
	docker-compose -f deployments/prod/docker-compose.prod.yml down

# Infrastructure (Development)
infra-dev-up:
	docker-compose -f deployments/dev/docker-compose.dev.yml up -d

infra-dev-down:
	docker-compose -f deployments/dev/docker-compose.dev.yml down