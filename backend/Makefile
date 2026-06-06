# ========== Variables ==========
APP_NAME    := jiraflow-backend
LOG_FILE    := logs/app.log
BUILD_DIR   := bin
MAIN_PATH   := cmd/main.go
SWAG_MAIN   := cmd/main.go
MIGRATE_DIR := migrations
DB_URL      ?= $(shell grep DATABASE_URL .env 2>/dev/null | cut -d= -f2)
DB_HOST     ?= $(shell grep ^DB_HOST     .env 2>/dev/null | cut -d= -f2)
DB_PORT     ?= $(shell grep ^DB_PORT     .env 2>/dev/null | cut -d= -f2)
DB_NAME     ?= $(shell grep ^DB_NAME     .env 2>/dev/null | cut -d= -f2)
DB_USER     ?= $(shell grep ^DB_USER     .env 2>/dev/null | cut -d= -f2)
DB_PASSWORD ?= $(shell grep ^DB_PASSWORD .env 2>/dev/null | cut -d= -f2)
MIGRATE     := $(shell go env GOPATH)/bin/migrate
SCRIPTS_DIR := scripts

.PHONY: all build run test lint fmt vet swagger migrate-up migrate-down \
        docker-up docker-down clean deps seed seed-test api-test \
        db-reset dev logs logs-f logs-errors logs-clear help

#make api-test — scripts/test_apis.sh ni ishga tushiradi
#make dev — build qilib serverni ishga tushiradi
#make seed — seed ma'lumotlarini DB ga yuklaydi
#make db-reset — DB ni tozalab qayta migrate + seed qiladi

# ========== Build ==========
build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

# ========== Quality ==========
fmt:
	gofmt -s -w .

vet:
	go vet ./...

lint:
	golangci-lint run ./...

test:
	go test -race -coverprofile=coverage.out ./...

test-cover: test
	go tool cover -html=coverage.out -o coverage.html

# ========== Swagger ==========
swagger:
	swag init -g $(SWAG_MAIN) --output docs --parseDependency

# ========== Database ==========
migrate-up:
	$(MIGRATE) -path $(MIGRATE_DIR) -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATE_DIR) -database "$(DB_URL)" down

migrate-create:
	$(MIGRATE) create -ext sql -dir $(MIGRATE_DIR) -seq $(name)

# Test ma'lumotlarini yuklash (o'chirmasdan, ON CONFLICT qoidasi bilan)
seed:
	@echo "→ Test seed ma'lumotlarini yuklamoqda..."
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) \
	  -U $(DB_USER) -d $(DB_NAME) \
	  -f $(SCRIPTS_DIR)/seed_test_data.sql
	@echo "✓ Seed bajarildi"

# API test skriptini ishga tushirish
api-test:
	@echo "→ API testlar ishga tushmoqda..."
	bash $(SCRIPTS_DIR)/test_apis.sh

# Build qilib serverni ishga tushirish
dev: build
	./$(BUILD_DIR)/$(APP_NAME)

# DB ni migration gacha reset qilib, qayta seed yuklash
db-reset:
	@echo "→ Migratsiyalar qayta ishlatilmoqda..."
	$(MIGRATE) -path $(MIGRATE_DIR) -database "$(DB_URL)" drop -f
	$(MIGRATE) -path $(MIGRATE_DIR) -database "$(DB_URL)" up
	@echo "→ Seed yuklanmoqda..."
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) \
	  -U $(DB_USER) -d $(DB_NAME) \
	  -f $(SCRIPTS_DIR)/seed_test_data.sql
	@echo "✓ DB reset va seed bajarildi"

# ========== Logs ==========
logs:
	@test -f $(LOG_FILE) || { echo "Log fayl topilmadi: $(LOG_FILE)"; exit 1; }
	tac $(LOG_FILE) | less -R

logs-f:
	tail -f $(LOG_FILE)

logs-errors:
	@test -f $(LOG_FILE) || { echo "Log fayl topilmadi: $(LOG_FILE)"; exit 1; }
	tac $(LOG_FILE) | grep '"level":"ERROR"' | less -R

logs-clear:
	@read -p "Log faylni tozalash. Ishonchingiz komilmi? [y/N] " ans && [ "$$ans" = "y" ]
	> $(LOG_FILE)
	@echo "✓ Log fayl tozalandi"

# ========== Docker ==========
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# ========== Cleanup ==========
clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html

deps:
	go mod tidy
	go mod verify

# ========== Default ==========
all: fmt vet build

help:
	@echo "Available targets:"
	@echo "  build          - Build binary to bin/"
	@echo "  run            - Run server directly"
	@echo "  fmt            - Format code (gofmt)"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run golangci-lint"
	@echo "  test           - Run tests with race detector"
	@echo "  test-cover     - Tests + HTML coverage report"
	@echo "  swagger        - Generate Swagger docs"
	@echo "  migrate-up     - Apply all migrations"
	@echo "  migrate-down   - Rollback last migration"
	@echo "  migrate-create name=xxx - Create new migration"
	@echo "  docker-up      - Start docker-compose services"
	@echo "  docker-down    - Stop docker-compose services"
	@echo "  deps           - go mod tidy + verify"
	@echo "  clean          - Remove build artifacts"
	@echo "  seed           - Load test seed data into DB"
	@echo "  db-reset       - Drop + migrate up + seed"
	@echo "  api-test       - Run API test script (scripts/test_apis.sh)"
	@echo "  dev            - Build and run server"
	@echo "  logs           - View log file newest first (less viewer)"
	@echo "  logs-f         - Follow log file in real-time (tail -f)"
	@echo "  logs-errors    - View only ERROR logs, newest first"
	@echo "  logs-clear     - Clear the log file (with confirmation)"
