# Глобальные переменные проекта
GO_VERSION = 1.24
GOLANGCI_LINT_VERSION = v2.1.5
GCI_VERSION = v0.13.6
GOFUMPT_VERSION = v0.8.0
BUF_VERSION = 1.53.0
PROTOC_GEN_GO_VERSION = v1.36.6
PROTOC_GEN_GO_GRPC_VERSION = v1.5.1
OGEN_VERSION = v1.12.0
YQ_VERSION = v4.45.2
GRPCURL_VERSION = v1.9.3

ROOT_DIR = $(shell pwd)
BIN_DIR = $(ROOT_DIR)/bin
GOLANGCI_LINT = $(BIN_DIR)/golangci-lint
GCI = $(BIN_DIR)/gci
GOFUMPT = $(BIN_DIR)/gofumpt
BUF = $(BIN_DIR)/buf
OGEN = $(BIN_DIR)/ogen
YQ = $(BIN_DIR)/yq
PROTOC_GEN_GO = $(BIN_DIR)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(BIN_DIR)/protoc-gen-go-grpc
GRPCURL = $(BIN_DIR)/grpcurl

NODE_MODULES_DIR = $(ROOT_DIR)/node_modules/.bin
REDOCLY = $(NODE_MODULES_DIR)/redocly

# Директории для сгенерированных файлов
PROTO_GEN_DIR = $(ROOT_DIR)/shared/pkg/proto
OPENAPI_GEN_DIR = $(ROOT_DIR)/shared/pkg/openapi
OPENAPI_SCHEMAS_DIR = $(ROOT_DIR)/shared/api

# Исходные файлы для генерации
PROTO_SOURCES = $(shell find shared/proto -name '*.proto' 2>/dev/null || true)
OPENAPI_MAIN_FILE = $(ROOT_DIR)/shared/api/order/v1/order.openapi.yaml
OPENAPI_BUNDLE_FILE = $(ROOT_DIR)/shared/api/bundles/order.openapi.v1.bundle.yaml

MODULES = inventory order payment

.PHONY: help install-formatters format install-golangci-lint lint
.PHONY: install-buf proto-install-plugins proto-lint proto-gen-dir
.PHONY: redocly-cli-install redocly-cli-bundle
.PHONY: ogen-install openapi-gen-dir create-dirs
.PHONY: yq-install grpcurl-install test-api clean all deps-update gen

.DEFAULT_GOAL := help

help:  ## Показать справку по всем командам
	@echo "Доступные команды:"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'

install-formatters:  ## Устанавливает форматтеры gci и gofumpt в ./bin
	@echo "📦 Проверяем установку форматтеров..."
	@if [ ! -f $(GOFUMPT) ]; then \
		echo '📦 Устанавливаем gofumpt $(GOFUMPT_VERSION)...'; \
		GOBIN=$(BIN_DIR) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION); \
	fi
	@if [ ! -f $(GCI) ]; then \
		echo '📦 Устанавливаем gci $(GCI_VERSION)...'; \
		GOBIN=$(BIN_DIR) go install github.com/daixiang0/gci@$(GCI_VERSION); \
	fi

format: install-formatters  ## Форматирует весь проект gofumpt + gci, исключая mocks
	@echo "🧼 Форматируем через gofumpt ..."
	@for module in $(MODULES); do \
		if [ -d "$$module" ]; then \
			echo "🧼 Форматируем $$module"; \
			find "$$module" -type f -name '*.go' ! -path '*/mocks/*' -exec $(GOFUMPT) -extra -w {} +; \
		fi; \
	done
	@echo "🎯 Сортируем импорты через gci ..."
	@for module in $(MODULES); do \
		if [ -d "$$module" ]; then \
			echo "🎯 Сортируем импорты в $$module"; \
			find "$$module" -type f -name '*.go' ! -path '*/mocks/*' -exec $(GCI) write -s standard -s default -s "prefix(github.com/olezhek28/microservices-course-olezhek-solution)" {} +; \
		fi; \
	done

install-golangci-lint:  ## Устанавливает golangci-lint в каталог bin
	@if [ ! -f $(GOLANGCI_LINT) ]; then \
		mkdir -p $(BIN_DIR); \
		echo "📦 Устанавливаем golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	fi



install-buf:  ## Устанавливает Buf в каталог bin
	@if [ ! -f $(BUF) ]; then \
		mkdir -p $(BIN_DIR) tmp-buf; \
		echo "📦 Устанавливаем Buf $(BUF_VERSION)..."; \
		curl -sSL \
			"https://github.com/bufbuild/buf/releases/download/v$(BUF_VERSION)/buf-$(shell uname -s)-$(shell uname -m).tar.gz" \
			| tar -xz -C tmp-buf; \
		mv tmp-buf/buf/bin/buf $(BUF); \
		rm -rf tmp-buf; \
		chmod +x $(BUF); \
	fi

proto-install-plugins:  ## Устанавливает protoc плагины в каталог bin
	@if [ ! -f $(PROTOC_GEN_GO) ]; then \
		echo '📦 Installing protoc-gen-go...'; \
		GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION); \
	fi
	@if [ ! -f $(PROTOC_GEN_GO_GRPC) ]; then \
		echo '📦 Installing protoc-gen-go-grpc...'; \
		GOBIN=$(BIN_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION); \
	fi

proto-lint: install-buf proto-install-plugins  ## Проверка .proto-файлов на соответствие стилю
	@cd shared/proto && $(BUF) lint

create-dirs:  ## Создает директории для сгенерированных файлов
	@mkdir -p $(PROTO_GEN_DIR)
	@mkdir -p $(OPENAPI_GEN_DIR)
	@mkdir -p $(BIN_DIR)

# Проверяем наличие исходных файлов
check-proto-sources:
	@if [ -z "$(PROTO_SOURCES)" ]; then \
		echo "❌ Нет .proto файлов в shared/proto/"; \
		echo "Доступные файлы:"; \
		find shared/proto -name '*.proto' 2>/dev/null | while read file; do echo "   - $$file"; done || echo "   не найдено"; \
		exit 1; \
	else \
		echo "✅ Найдены .proto файлы:"; \
		for file in $(PROTO_SOURCES); do \
			echo "   - $$file"; \
		done; \
	fi

check-openapi-sources:
	@echo "🔍 Ищем OpenAPI файлы в shared/api/..."
	@if [ -f "$(OPENAPI_MAIN_FILE)" ]; then \
		echo "✅ Найден главный OpenAPI файл: $(OPENAPI_MAIN_FILE)"; \
	else \
		echo "❌ Главный OpenAPI файл не найден: $(OPENAPI_MAIN_FILE)"; \
		echo "Доступные файлы в shared/api/:"; \
		find shared/api -type f \( -name '*.yaml' -o -name '*.yml' \) 2>/dev/null | while read file; do echo "   - $$file"; done || echo "   не найдено"; \
		exit 1; \
	fi

proto-gen-dir: create-dirs install-buf proto-install-plugins proto-lint check-proto-sources  ## Генерация Go-кода из .proto в shared/pkg/proto
	@echo "🚀 Генерация Go-кода из .proto файлов в $(PROTO_GEN_DIR)..."
	@echo "📁 Текущая директория: $(shell pwd)"
	@echo "📁 PROTO_GEN_DIR: $(PROTO_GEN_DIR)"
	@cd shared/proto && $(BUF) generate --output $(PROTO_GEN_DIR)
	@echo "✅ Proto код сгенерирован в $(PROTO_GEN_DIR)"

redocly-cli-install:  ## Установить локально Redocly CLI
	@if [ ! -f $(REDOCLY) ]; then \
		echo "📦 Устанавливаем Redocly CLI..."; \
		npm ci; \
	fi

redocly-cli-bundle: create-dirs redocly-cli-install  ## Собрать все схемы OpenAPI в общие файлы через локальный redocly
	@echo "📦 Бандлим OpenAPI схемы..."
	@if [ -f "$(OPENAPI_MAIN_FILE)" ]; then \
		bundle_file="shared/api/bundles/order.openapi.v1.bundle.yaml"; \
		echo "📦 Бандлим $(OPENAPI_MAIN_FILE) -> $$bundle_file"; \
		mkdir -p "$$(dirname $$bundle_file)"; \
		$(REDOCLY) bundle "$(OPENAPI_MAIN_FILE)" -o "$$bundle_file" || echo "⚠️  Ошибка бандлинга $(OPENAPI_MAIN_FILE)"; \
	else \
		echo "❌ Главный OpenAPI файл не найден: $(OPENAPI_MAIN_FILE)"; \
	fi

ogen-install:  ## Скачивает ogen в папку bin
	@if [ ! -f $(OGEN) ]; then \
		mkdir -p $(BIN_DIR); \
		echo "📦 Устанавливаем ogen $(OGEN_VERSION)..."; \
		GOBIN=$(BIN_DIR) go install github.com/ogen-go/ogen/cmd/ogen@$(OGEN_VERSION); \
	fi

yq-install:  ## Устанавливает yq в bin/ при необходимости
	@if [ ! -f $(YQ) ]; then \
		echo '📦 Installing yq...'; \
		GOBIN=$(BIN_DIR) go install github.com/mikefarah/yq/v4@$(YQ_VERSION); \
	fi

openapi-gen-dir: create-dirs ogen-install yq-install redocly-cli-bundle  ## Генерация Go-кода из OpenAPI в shared/pkg/openapi
	@echo "🚀 Генерация кода из OpenAPI спецификаций в $(OPENAPI_GEN_DIR)..."
	@if [ -f "$(OPENAPI_BUNDLE_FILE)" ]; then \
		echo "🚀 Generating from bundle: $(OPENAPI_BUNDLE_FILE)"; \
		target_dir="order/v1"; \
		package="orderv1"; \
		echo "📁 Target directory: $$target_dir"; \
		echo "📦 Package: $$package"; \
		mkdir -p "$(OPENAPI_GEN_DIR)/$$target_dir"; \
		$(OGEN) \
			--target "$(OPENAPI_GEN_DIR)/$$target_dir" \
			--package "$$package" \
			--clean \
			"$(OPENAPI_BUNDLE_FILE)" || (echo "❌ Ошибка генерации из $(OPENAPI_BUNDLE_FILE)"; exit 1); \
		echo "✅ Успешно сгенерировано: $(OPENAPI_GEN_DIR)/$$target_dir"; \
	else \
		echo "❌ Бандл OpenAPI не найден: $(OPENAPI_BUNDLE_FILE)"; \
		echo "Создаем бандл..."; \
		make redocly-cli-bundle; \
		if [ -f "$(OPENAPI_BUNDLE_FILE)" ]; then \
			echo "🚀 Generating from bundle: $(OPENAPI_BUNDLE_FILE)"; \
			mkdir -p "$(OPENAPI_GEN_DIR)/$$target_dir"; \
			$(OGEN) \
				--target "$(OPENAPI_GEN_DIR)/$$target_dir" \
				--package "$$package" \
				--clean \
				"$(OPENAPI_BUNDLE_FILE)" || (echo "❌ Ошибка генерации из $(OPENAPI_BUNDLE_FILE)"; exit 1); \
		else \
			echo "❌ Не удалось создать бандл OpenAPI"; \
			exit 1; \
		fi; \
	fi

gen: proto-gen-dir openapi-gen-dir  ## Генерация всех proto и OpenAPI деклараций в соответствующие директории
	@echo "✅ Вся генерация завершена!"

deps-update:  ## Обновление зависимостей в go.mod во всех модулях
	@echo "🔄 Обновление зависимостей в go.mod во всех модулях"
	@for mod in $(MODULES); do \
		if [ -d "$$mod" ]; then \
			echo "🔄 Обновление зависимостей в $$mod"; \
			(cd "$$mod" && go mod tidy -compat=1.24) || exit 1; \
		fi; \
	done
	@echo "🔄 Обновление зависимостей в shared"; \
	(cd shared && go mod tidy -compat=1.24) || exit 1;

grpcurl-install:  ## Устанавливает grpcurl в каталог bin
	@if [ ! -f $(GRPCURL) ]; then \
		echo '📦 Устанавливаем grpcurl $(GRPCURL_VERSION)...'; \
		GOBIN=$(BIN_DIR) go install github.com/fullstorydev/grpcurl/cmd/grpcurl@$(GRPCURL_VERSION); \
	fi

test-api: grpcurl-install  ## 🧪 Запуск тестов для проверки API микросервисов
	@echo "🧪 Тестирование API микросервисов через gRPC и REST"
	@echo "⚠️  Реализация тестов будет добавлена позже"

clean:  ## Очистка сгенерированных файлов
	@echo "🧹 Очистка сгенерированных файлов..."
	@rm -rf $(BIN_DIR)
	@rm -rf shared/api/bundles
	@rm -rf $(PROTO_GEN_DIR)/*
	@rm -rf $(OPENAPI_GEN_DIR)/*
	@for module in $(MODULES); do \
		if [ -d "$$module" ]; then \
			find "$$module" -name "*.gen.go" -type f -delete; \
			find "$$module" -name "*_ogen*" -type f -delete; \
			find "$$module" -name "*pb.go" -type f -delete; \
			find "$$module" -name "*_grpc.pb.go" -type f -delete; \
		fi; \
	done
	@echo "✅ Очистка завершена"

all: format lint gen deps-update  ## Выполнить все основные задачи: форматирование, линтинг, генерацию кода и обновление зависимостей
	@echo "✅ Все задачи выполнены успешно!"

# Отладочные цели для проверки
debug-sources:  ## Показать исходные файлы для генерации
	@echo "=== Proto sources ==="
	@find shared/proto -name '*.proto' 2>/dev/null | while read file; do echo "  $$file"; done || echo "  не найдено"
	@echo ""
	@echo "=== OpenAPI sources ==="
	@find shared/api -name '*.yaml' -o -name '*.yml' 2>/dev/null | while read file; do echo "  $$file"; done || echo "  не найдено"
	@echo ""
	@echo "=== Главный OpenAPI файл ==="
	@if [ -f "$(OPENAPI_MAIN_FILE)" ]; then \
		echo "  ✅ $(OPENAPI_MAIN_FILE) - существует"; \
	else \
		echo "  ❌ $(OPENAPI_MAIN_FILE) - не существует"; \
	fi
	@echo ""
	@echo "=== Generated files ==="
	@echo "PROTO_GEN_DIR ($(PROTO_GEN_DIR)):"
	@ls -la $(PROTO_GEN_DIR) 2>/dev/null || echo "  не существует или пусто"
	@echo ""
	@echo "OPENAPI_GEN_DIR ($(OPENAPI_GEN_DIR)):"
	@ls -la $(OPENAPI_GEN_DIR) 2>/dev/null || echo "  не существует или пусто"

debug-proto-paths:  ## Показать пути для proto генерации
	@echo "📁 ROOT_DIR: $(ROOT_DIR)"
	@echo "📁 PROTO_GEN_DIR: $(PROTO_GEN_DIR)"
	@echo "📁 Относительный путь из shared/proto в pkg/proto: ../../../pkg/proto"
	@echo "📁 Текущая директория proto:"
	@cd shared/proto && pwd