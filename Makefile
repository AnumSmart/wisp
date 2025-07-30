.PHONY: up down install-deps

include .env
#----------------------------------------------------------------------------------------
LOCAL_BIN := $(CURDIR)/bin
LOCAL_MIGRATION_DIR := $(MIGRATION_DIR)
LOCAL_MIGRATION_DSN := user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DATABASE_NAME) host=$(PG_HOST) port=$(PG_PORT) sslmode=$(PG_SSLMODE)

#----------------------------------------------------------------------------------------
# Установка бинарника утилиты для миграции goose
install-deps:
	@echo "Installation of goose into: $(LOCAL_BIN)..."
	@if not exist "$(LOCAL_BIN)" mkdir "$(LOCAL_BIN)"
	go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	@move "%GOPATH%\bin\goose.exe" "$(LOCAL_BIN)\" >nul 2>&1 || echo "Goose is already in right folder"
	@echo "Goose has been installed: $(LOCAL_BIN)\goose.exe"

local-migration-status:
	@"$(LOCAL_BIN)\goose.exe" -dir "$(LOCAL_MIGRATION_DIR)" postgres "$(LOCAL_MIGRATION_DSN)" status -v

local-migration-up:
	@"$(LOCAL_BIN)\goose.exe" -dir "$(LOCAL_MIGRATION_DIR)" postgres "$(LOCAL_MIGRATION_DSN)" up -v

local-migration-down:
	@"$(LOCAL_BIN)\goose.exe" -dir "$(LOCAL_MIGRATION_DIR)" postgres "$(LOCAL_MIGRATION_DSN)" down -v

#----------------------------------------------------------------------------------------
# Генерация новой миграции
# Пример: make migration-create name=create_users_table
migration-create:
	@"$(LOCAL_BIN)\goose.exe" -dir "$(LOCAL_MIGRATION_DIR)" create "$(name)" sql
	@echo "The Migration file was created in folder: $(LOCAL_MIGRATION_DIR)"

#----------------------------------------------------------------------------------------
# Запуск и останов контейнеров
up:
	docker-compose up -d
down:
	docker-compose down

# Дополнительные команды (опционально)
start: up
stop: down
restart: down up
status:
	docker-compose ps


