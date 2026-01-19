SHELL := /bin/bash

BIN_DIR      ?= $(PWD)/bin
GOBIN        ?= $(BIN_DIR)
GOLANGCI_VER ?= v1.60.3
GOFUMPT_VER  ?= v0.6.0

APP_CMD      := ./cmd/tracker
HTTP_PORT    ?= 8080
DB_URI       ?= postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable

.PHONY: all lint fmt test run deps tools docker-build docker-up docker-down help

all: fmt lint test ## Форматирование, линт, тесты

help: ## Показать справку
	@grep -E '^[a-zA-Z0-9_\-]+:.*?## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

tools: ## Установить локальные инструменты (gofumpt, golangci-lint)
	@mkdir -p $(BIN_DIR)
	@[ -x "$(BIN_DIR)/gofumpt" ] || { echo "📦 Installing gofumpt $(GOFUMPT_VER)"; GOBIN=$(BIN_DIR) go install mvdan.cc/gofumpt@$(GOFUMPT_VER); }
	@[ -x "$(BIN_DIR)/golangci-lint" ] || { echo "📦 Installing golangci-lint $(GOLANGCI_VER)"; GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VER); }

fmt: tools ## Форматировать Go-код
	@$(BIN_DIR)/gofumpt -w $$(find . -name '*.go' -not -path './shared/pkg/*')

lint: tools ## Линтер Go-кода
	@$(BIN_DIR)/golangci-lint run ./...

test: ## Запустить тесты
	@go test ./...

run: ## Локальный запуск приложения
	@HTTP_PORT=$(HTTP_PORT) DB_URI=$(DB_URI) go run $(APP_CMD)

deps: ## Обновить зависимости
	@go mod tidy

docker-build: ## Собрать Docker image приложения
	@docker build -t tasktracker-app:local .

docker-up: ## Поднять приложение + postgres через docker compose
	@cd deploy && docker compose up -d

docker-down: ## Остановить docker compose
	@cd deploy && docker compose down
