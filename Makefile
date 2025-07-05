# Variables
SERVICES := user shop cart order payment review shipping

# Service migration directory mapping
USER_MIGRATION_DIR := internal/services/user-service/internal/db/migrations
SHOP_MIGRATION_DIR := internal/services/product-service/internal/db/migrations
CART_MIGRATION_DIR := internal/services/cart-service/internal/db/migrations
ORDER_MIGRATION_DIR := internal/services/order-service/internal/db/migrations
PAYMENT_MIGRATION_DIR := internal/services/payment-service/internal/db/migrations
REVIEW_MIGRATION_DIR := internal/services/review-service/internal/db/migrations
SHIPPING_MIGRATION_DIR := internal/services/shipping-service/internal/db/migrations

lint:
	golangci-lint run

# Generic migration function
define run_migration
	@echo "$(2) $(1) Service..."
	@$(eval SERVICE_UPPER := $(shell echo $(1) | tr '[:lower:]' '[:upper:]'))
	@set GOOSE_DRIVER=$$($(SERVICE_UPPER)_SERVICE_GOOSE_DRIVER) && \
	 set GOOSE_DBSTRING=$$($(SERVICE_UPPER)_SERVICE_GOOSE_DBSTRING) && \
	 goose -dir $$($(SERVICE_UPPER)_MIGRATION_DIR) $(3)
endef

# Create migration for specific service
define create_migration_service
create-migration-$(1):
	@if "$(name)" == "" ( \
		echo ❌ Thiếu tên migration. Dùng: make create-migration-$(1) name=ten_migration & exit /b 1 \
	) else ( \
		goose -dir $$($(shell echo $(1) | tr '[:lower:]' '[:upper:]')_MIGRATION_DIR) create -s $(name) sql \
	)
endef

# Generate targets for all services
$(foreach service,$(SERVICES),$(eval $(call create_migration_service,$(service))))

# Migration commands
migrate:
	@if "$(service)" == "" ( \
		echo ❌ Thiếu tên service. Dùng: make migrate service=ten_service action=up/down/status & exit /b 1 \
	) else if "$(action)" == "" ( \
		echo ❌ Thiếu action. Dùng: make migrate service=ten_service action=up/down/status & exit /b 1 \
	) else ( \
		$(call run_migration,$(service),🚀 Migration,$(action)) \
	)

# Individual service migration commands
migrate-up-%:
	$(call run_migration,$*,🚀 Chạy migration cho,up)

migrate-down-%:
	$(call run_migration,$*,⬇️ Rollback migration cho,down)

migrate-status-%:
	$(call run_migration,$*,📊 Kiểm tra status migration cho,status)

# Bulk operations
migrate-up-all:
	@echo "� Chạy migration cho tất cả services..."
	@$(foreach service,$(SERVICES),make migrate-up-$(service);)

migrate-down-all:
	@echo "⬇️ Rollback migration cho tất cả services..."
	@$(foreach service,$(SERVICES),make migrate-down-$(service);)

migrate-status-all:
	@echo "📊 Kiểm tra status migration cho tất cả services..."
	@$(foreach service,$(SERVICES),make migrate-status-$(service);)

build:
	docker compose up -d --build 
up:
	docker compose up -d 
down:
	docker compose down
