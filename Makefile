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

# Infrastructure
infra up:
 	docker-compose up -d

infra down:
  	docker-compose down

# Development
dev-up:
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

dev-down:
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml down

# Production  
prod-up:
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d