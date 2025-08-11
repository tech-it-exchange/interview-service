generate-swagger: ## Swag init
	swag init --parseDependency --parseInternal

.PHONY: help build run test clean generate-swagger

help: ## Показать справку по командам
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Собрать проект
	export PATH=/opt/homebrew/opt/go/libexec/bin:$$PATH && go build -o bin/interview-service .

run: ## Запустить сервис
	export PATH=/opt/homebrew/opt/go/libexec/bin:$$PATH && go run .

run-background: ## Запустить сервис в фоне
	export PATH=/opt/homebrew/opt/go/libexec/bin:$$PATH && go run . &

stop: ## Остановить сервис
	pkill -f "go run ." || true

test: ## Запустить тесты
	export PATH=/opt/homebrew/opt/go/libexec/bin:$$PATH && go test ./...

test-coverage: ## Запустить тесты с покрытием
	export PATH=/opt/homebrew/opt/go/libexec/bin:$$PATH && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

clean: ## Очистить артефакты сборки
	rm -rf bin/ coverage.out

generate-swagger: ## Генерировать Swagger документацию
	export PATH=/opt/homebrew/opt/go/libexec/bin:$$PATH && swag init --parseDependency --parseInternal

lint: ## Проверить код линтером
	export PATH=/opt/homebrew/opt/go/libexec/bin:$$PATH && golangci-lint run

fmt: ## Форматировать код
	export PATH=/opt/homebrew/opt/go/libexec/bin:$$PATH && go fmt ./...

mod-tidy: ## Привести в порядок зависимости
	export PATH=/opt/homebrew/opt/go/libexec/bin:$$PATH && go mod tidy

check-api: ## Проверить API эндпоинты
	@echo "Проверка API эндпоинтов..."
	@echo "1. Swagger UI: http://localhost:9099/docs/index.html"
	@echo "2. Get Loss: http://localhost:8085/api/v2/client/user/get-loss"
	@echo "3. Kafka Topics: http://localhost:8085/api/v2/client/kafka/topics"
	@echo "4. Create Order: POST http://localhost:8085/api/v2/client/order/create"
	@echo "5. Create Instrument: POST http://localhost:8085/api/v2/client/instrument/create"

install-deps: ## Установить зависимости разработки
	export PATH=/opt/homebrew/opt/go/libexec/bin:$$PATH && go install github.com/swaggo/swag/cmd/swag@latest
	export PATH=/opt/homebrew/opt/go/libexec/bin:$$PATH && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

dev: ## Режим разработки (сборка + запуск)
	make build && make run

all: ## Полная сборка и проверка
	make clean && make mod-tidy && make generate-swagger && make build && make test