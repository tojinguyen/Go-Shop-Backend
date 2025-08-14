lint:
	golangci-lint run

build:
	docker compose up -d --build 
up:
	docker compose up -d 
down:
	docker compose down
recreate:
	docker compose up -d --force-recreate

# ===================================================================
# Protobuf/gRPC Generation
# ===================================================================
# Biến để dễ dàng thay đổi đường dẫn nếu cần
PROTO_DIR := proto
GEN_DIR_GO := proto/gen/go

# Automatically find all .proto files.
PROTO_FILES := $(patsubst $(PROTO_DIR)/%,%,$(wildcard $(PROTO_DIR)/*/v1/*.proto))

# Lệnh chính để generate code
.PHONY: proto-gen
proto-gen:
	@echo "Generating Go code from Protobuf definitions..."
	@echo "Processing files: $(PROTO_FILES)"
	@protoc --proto_path=$(PROTO_DIR) \
	       --go_out=paths=source_relative:$(GEN_DIR_GO) \
	       --go-grpc_out=paths=source_relative:$(GEN_DIR_GO) \
	       $(PROTO_FILES)
	@echo "Protobuf/gRPC code generated successfully."

# Lệnh để cập nhật go.mod trong thư mục generated code
.PHONY: proto-tidy
proto-tidy:
	@echo "Tidying Go modules in generated proto directory..."
	@cd $(GEN_DIR_GO) && go mod tidy
	@echo "Go modules for generated code are up to date."

# Lệnh tổng hợp: generate code và sau đó tidy go.mod
.PHONY: proto
proto: proto-gen proto-tidy ## Generate all Protobuf/gRPC code and tidy modules



# ===================================================================
# Database Seeding
# ===================================================================

# Lệnh để seed 50,0000 users với phân bố thực tế
.PHONY: seed-users
seed-users:
	@echo "Seeding user-service database with 50,0000 users (realistic e-commerce distribution)..."
	@echo "Distribution: ~87% customers, ~10% sellers, ~2.5% shippers, ~0.5% admins"
	@cd internal/services/user-service && go run ./cmd/seeder/main.go -total=500000
	@echo "50K users seeding complete!"

# Lệnh để seed 1000 shops
.PHONY: seed-shops
seed-shops:
	@echo "Seeding shop-service database with 1000 shops..."
	@cd internal/services/shop-service && go run ./cmd/seeder/main.go -shops=1000
	@echo "Shop service seeding complete."

# Lệnh để seed 1000 products
.PHONY: seed-products
seed-products:
	@echo "Seeding product-service database with 10000 products..."
	@cd internal/services/product-service && go run ./cmd/seeder/main.go -products=10000
	@echo "Product service seeding complete."
