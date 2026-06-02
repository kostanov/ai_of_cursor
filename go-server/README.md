# Go Server

Этот каталог содержит порт `api.py` на Go с эквивалентными маршрутами и поведением.

## Что реализовано

- SQLite база `test.db` и автоинициализация таблицы `users`.
- Эндпоинты:
  - `POST /adduser`
  - `GET /user/{uid}`
  - `GET|POST /activate/{uid}`
  - `GET /slow`
  - `GET /wrong`
- Потокобезопасный список активных пользователей (лимит 100).
- Сервер слушает `0.0.0.0:8080`.

## Требования

- Go 1.22+

## Структура проекта

```
go-server/
├── cmd/server/main.go       # точка входа
├── internal/
│   ├── config/              # константы (порт, путь к БД)
│   ├── models/              # DTO и модели ответов
│   ├── database/            # подключение и схема SQLite
│   ├── repository/          # работа с таблицей users
│   ├── service/             # бизнес-логика (активные пользователи, фоновые задачи)
│   ├── handler/             # HTTP-обработчики
│   ├── router/              # регистрация маршрутов
│   └── util/                # JSON-ответы, разбор пути
├── go.mod
└── README.md
```

## Команды (just)

| Команда | Описание |
|---------|----------|
| `just` | список рецептов |
| `just run` | запуск сервера |
| `just build` | сборка бинарника `server` |
| `just tidy` | `go mod tidy` |
| `just fmt` | `gofmt` для `cmd` и `internal` |
| `just vet` | `go vet ./...` |
| `just test` | `go test ./...` |

## Установка и запуск

```bash
cd go-server
go mod tidy
just run
```

Сервер поднимется на `http://127.0.0.1:8080`.

## Сборка

```bash
cd go-server
go build -o server ./cmd/server
./server
```

## API: примеры

### 1) Создать пользователя

```bash
curl -s -X POST http://127.0.0.1:8080/adduser \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice"}'
```

Ответ:

```json
{"status":"ok","id":1,"name":"Alice"}
```

### 2) Получить пользователя

```bash
curl -s http://127.0.0.1:8080/user/1
```

Ответ `200`:

```json
{"id":1,"name":"Alice"}
```

Если не найден:

```json
{"error":"not_found"}
```

### 3) Активировать пользователя

```bash
curl -s -X POST http://127.0.0.1:8080/activate/1
```

Ответ:

```json
{"status":"ok","active":[1]}
```

### 4) Фоновая задача

```bash
curl -s http://127.0.0.1:8080/slow
```

Ответ `202`:

```json
{"status":"scheduled"}
```

### 5) Демонстрация ошибки

```bash
curl -s http://127.0.0.1:8080/wrong
```

Ответ `500`:

```json
{"msg":"error","detail":"division by zero"}
```

## Примечания

- База `test.db` создается в каталоге `go-server`.
- Для удаления локальных данных можно удалить файл `test.db`.
