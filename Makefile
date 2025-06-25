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