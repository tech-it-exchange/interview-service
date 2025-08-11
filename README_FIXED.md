# Interview Service - Исправленная версия

## 🎯 Статус исправления багов

| Баг | Статус | Описание |
|-----|--------|----------|
| 1 | ✅ **ИСПРАВЛЕН** | Приложение запускается без ошибок |
| 2 | ✅ **ИСПРАВЛЕН** | Swagger UI доступен по http://localhost:9099/docs/index.html |
| 3 | ✅ **ИСПРАВЛЕН** | Эндпоинт get-loss работает корректно |
| 4 | ✅ **ИСПРАВЛЕН** | JSON возвращает `{"loss":-1000}` |
| 5 | ✅ **ИСПРАВЛЕН** | Валидация side="ask" работает |
| 6 | ✅ **ИСПРАВЛЕН** | Kafka topics API возвращает массив топиков |
| 7 | ✅ **ИСПРАВЛЕН** | SQL Injection уязвимость устранена |
| 8 | ✅ **ИСПРАВЛЕН** | Логика создания инструмента оптимизирована |

## 🚀 Быстрый запуск

### Предварительные требования
- Go 1.22+ (рекомендуется Go 1.24)
- Homebrew (для macOS)

### Установка Go
```bash
# Установка актуальной версии Go
brew install go

# Проверка версии
go version
```

### Запуск проекта
```bash
# Клонирование и переход в директорию
cd interview-service

# Установка зависимостей
go mod tidy

# Запуск сервиса
make run
```

## 📡 Доступные эндпоинты

### HTTP API (порт 8085)
- **GET** `/api/v2/client/user/get-loss` - Получить убыток
- **POST** `/api/v2/client/order/create` - Создать ордер
- **POST** `/api/v2/client/instrument/create` - Создать инструмент
- **GET** `/api/v2/client/user/by/{name}` - Получить пользователя по имени
- **GET** `/api/v2/client/instrument/balance/{instrumentId}` - Получить баланс инструмента
- **GET** `/api/v2/client/kafka/topics` - Получить активные топики Kafka

### Swagger UI (порт 9099)
- **URL:** http://localhost:9099/docs/index.html

## 🛠️ Команды Makefile

```bash
make help              # Показать справку по командам
make build             # Собрать проект
make run               # Запустить сервис
make run-background    # Запустить в фоне
make stop              # Остановить сервис
make test              # Запустить тесты
make clean             # Очистить артефакты
make check-api         # Проверить API эндпоинты
make dev               # Режим разработки
```

## 🧪 Тестирование API

### 1. Проверка Swagger
```bash
curl http://localhost:9099/docs/index.html
```

### 2. Тест get-loss
```bash
curl http://localhost:8085/api/v2/client/user/get-loss
# Ожидаемый результат: {"loss":-1000}
```

### 3. Тест создания ордера
```bash
curl -X POST http://localhost:8085/api/v2/client/order/create \
  -H "Content-Type: application/json" \
  -d '{"instrumentId":"0fae8035-7e4b-4f60-b1da-e3752d53f2f4","side":"ask","qty":"1","price":"100"}'
```

### 4. Тест создания инструмента
```bash
curl -X POST http://localhost:8085/api/v2/client/instrument/create \
  -H "Content-Type: application/json" \
  -d '{"instrumentId":"0fae8035-7e4b-4f60-b1da-e3752d53f2f4","isActive":true,"isListed":true}'
```

### 5. Тест Kafka topics
```bash
curl http://localhost:8085/api/v2/client/kafka/topics
# Ожидаемый результат: ["confirm.topic","command.topic"]
```

## 🔒 Исправленные уязвимости

### SQL Injection
- **Проблема:** Использование `fmt.Sprintf` в SQL запросах
- **Решение:** Заменено на параметризованные запросы
- **Файл:** `adapter/outbound/postgres/repositories/user_repository.go`

### Производительность
- **Проблема:** Создание 100M горутин для начисления баланса
- **Решение:** Оптимизировано до одного арифметического действия
- **Файл:** `application/worker/events/balance_worker.go`

## 📊 Архитектура проекта

```
interview-service/
├── adapter/           # Адаптеры (HTTP, Kafka, Storage)
├── application/       # Бизнес-логика (Use Cases, Services)
├── domain/           # Доменные модели
├── infrastructure/   # Конфигурация и утилиты
├── port/            # Интерфейсы (Ports)
└── main.go          # Точка входа
```

## 🚨 Известные ограничения

1. **Mock реализации:** Некоторые компоненты (Kafka, PostgreSQL) используют mock-реализации
2. **Конфигурация:** Настройки захардкожены в коде
3. **Логирование:** Используется базовый `log.Logger` вместо структурированного
4. **Тесты:** Отсутствуют unit и integration тесты

## 🔧 Рекомендации по доработке

1. **Добавить тесты** для всех use case'ов
2. **Реализовать реальные** Kafka и PostgreSQL адаптеры
3. **Добавить мониторинг** (Prometheus, health checks)
4. **Улучшить обработку ошибок** и graceful shutdown
5. **Добавить валидацию** входных данных
6. **Реализовать аутентификацию** и авторизацию

## 📝 Логи

При запуске сервиса вы увидите:
```
[app] 2025/08/11 18:27:29 Запуск http сервера на порту :8085
[app] 2025/08/11 18:27:29 Swagger UI available on port 9099 path /docs/
```

## 🆘 Устранение неполадок

### Проблема: "command not found: go"
```bash
export PATH=/opt/homebrew/opt/go/libexec/bin:$PATH
```

### Проблема: "port already in use"
```bash
make stop
# или
lsof -ti:8085 | xargs kill -9
lsof -ti:9099 | xargs kill -9
```

### Проблема: "module requires Go 1.22"
```bash
go version  # Проверьте версию
brew install go  # Обновите Go
```

## 📞 Поддержка

При возникновении проблем:
1. Проверьте версию Go: `go version`
2. Убедитесь, что порты 8085 и 9099 свободны
3. Проверьте логи запуска
4. Используйте `make help` для списка доступных команд
