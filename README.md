# url-shortener

Простое и производительное REST API на Go для сокращения URL-адресов.

Поддерживает создание коротких ссылок, получение оригинальных URL, удаление, метрики Prometheus и автогенерацию Swagger-документации.

---

## Возможности

- Создание короткой ссылки `POST /urls`
- Получение оригинального URL по короткому коду `GET /urls/{short}`
- Удаление короткой ссылки `DELETE /urls/{short}`
- Проверка статуса сервиса `GET /health`
- Метрики Prometheus `GET /metrics`
- Swagger-документация `GET /swagger/index.html`
- Легковесная SQLite-база
- Middleware для логирования, метрик и обработки ошибок

---

## Установка

1. Клонируйте репозиторий:

```bash
git clone https://github.com/zen-flo/url-shortener.git
cd url-shortener
```

2. Соберите сервер:

```bash
go build -o url-shortener ./cmd
```

3. Запустите:

```bash
./url-shortener
```

По умолчанию сервер поднимается на http://localhost:8080

---

## Примеры использования

### Создать короткий URL

```bash
curl -X POST http://localhost:8080/urls \
-H "Content-Type: application/json" \
-d '{"url":"https://google.com"}'
```

### Получить оригинальный URL

```bash
curl http://localhost:8080/urls/abc123
```

### Удалить короткий URL

```bash
curl -X DELETE http://localhost:8080/urls/abc123
```

### Проверить статус сервиса

```bash
curl http://localhost:8080/health
```

### Посмотреть метрики Prometheus

```bash
curl http://localhost:8080/metrics
```

### Открыть Swagger UI

```bash
http://localhost:8080/swagger/index.html
```

---

## Структура проекта

```bash
.
├── cmd/
├── internal/
│   ├── db/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   └── service/
├── docs/ (Swagger)
└── main.go
```

---

## Тестирование

Проект покрыт юнит и интеграционными тестами, включая тесты HTTP-роутера и in-memory SQLite.

### Запуск всех тестов с покрытием:

```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### Текущее покрытие:

* cmd: ~40%

* db: ~75%

* handler: ~73%

* middleware: 100%

* service: ~77%

**Общее покрытие:** ~69%

---
