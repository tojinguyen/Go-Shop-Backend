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

.PHONY: seed-users
seed-users:
	@echo "🌱 Seeding user-service database with 50000 customers and 1000 shippers..."
	@go run ./internal/services/user-service/cmd/seeder/main.go -users=50000 -shippers=1000


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