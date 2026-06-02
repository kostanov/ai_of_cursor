# ai-of-cursor

Учебный Python-проект: утилиты для работы с пользователями в SQLite (`utils.py`) и HTTP API на Flask (`api.py`).

## Требования

- Python >= 3.14
- [uv](https://docs.astral.sh/uv/) — управление зависимостями и виртуальным окружением
- [just](https://github.com/casey/just) — команды разработки (опционально, но рекомендуется)

## Быстрый старт

```bash
uv init   # только при первом клонировании, если pyproject.toml ещё нет
uv sync
```

Активация окружения (опционально):

```bash
source .venv/bin/activate
```

## Зависимости

Основной источник зависимостей — `pyproject.toml` и lock-файл `uv.lock`.

Для установки без uv можно использовать pip:

```bash
pip install -r requirements.txt
pip install -r requirements-dev.txt   # линтер (ruff)
```

| Файл | Назначение |
|------|------------|
| `requirements.txt` | runtime: Flask и транзитивные пакеты |
| `requirements-dev.txt` | dev: ruff |
| `pyproject.toml` | декларация проекта для uv |

Обновить lock и экспорт requirements:

```bash
uv sync
uv export --format requirements-txt --no-hashes --no-dev -o requirements.txt
uv export --format requirements-txt --no-hashes --only-dev -o requirements-dev.txt
```

## Команды (just)

| Команда | Описание |
|---------|----------|
| `just` / `just default` | список рецептов |
| `just run` | запуск `main.py` |
| `just api` | запуск Flask API (`api.py`, порт **8080**) |
| `just lint` | `uv run ruff check .` |
| `just fix` | `uv run ruff check --fix .` |
| `just format` | `uv run ruff format .` |

Без just:

```bash
uv run python main.py
uv run python api.py
```

## Структура проекта

```
.
├── api.py              # Flask HTTP API (test.db)
├── utils.py            # библиотека пользователей (users.db)
├── main.py             # точка входа-заглушка
├── pyproject.toml      # метаданные и зависимости (uv)
├── requirements.txt    # runtime-зависимости (pip)
├── requirements-dev.txt
├── justfile
├── CHANGE_LOG.md       # история рефакторингов
└── README.md
```

## Модуль `utils.py`

Локальная работа с пользователями без HTTP.

- БД: `users.db`, таблица `users` (id, name, tags в JSON).
- Пароли: PBKDF2-HMAC-SHA256 в `passwords.txt` (соль + хеш).
- Активные пользователи: потокобезопасный список (до 5 в `set_active`, до 100 в API).

Основные функции:

| Функция | Описание |
|---------|----------|
| `init_db()` | создаёт таблицу при необходимости |
| `add_user(name, tags=None)` | добавляет пользователя, возвращает `int` id |
| `get_user_by_name(name)` | `dict` или `None` |
| `store_password(user_id, password)` | сохраняет хеш пароля |
| `set_active(user_id)` | помечает пользователя активным |
| `get_active_users_snapshot()` | копия списка активных id |
| `is_admin(user)` | проверка роли `admin` |
| `self_test()` | smoke-тест модуля |

Запуск self-test:

```bash
uv run python utils.py
```

## HTTP API (`api.py`)

Сервер: `http://0.0.0.0:8080` (режим `debug=True` при прямом запуске).

БД: `test.db`, таблица `users` (id, name). Соединение — на контекст запроса (`flask.g`), закрывается в `teardown_appcontext`.

### Эндпоинты

#### `POST /adduser`

Создание пользователя.

**Тело (JSON):**

```json
{ "name": "Alice" }
```

**Ответ `201`:**

```json
{ "status": "ok", "id": 1, "name": "Alice" }
```

**Ошибка `400`:** пустое или нестроковое `name`.

**Пример:**

```bash
curl -s -X POST http://127.0.0.1:8080/adduser \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice"}'
```

#### `GET /user/<uid>`

Пользователь по id.

**Ответ `200`:** `{ "id": 1, "name": "Alice" }`  
**Ответ `404`:** `{ "error": "not_found" }`

```bash
curl -s http://127.0.0.1:8080/user/1
```

#### `POST|GET /activate/<uid>`

Добавляет id в список активных (потокобезопасно, лимит 100).

**Ответ `200`:**

```json
{ "status": "ok", "active": [1, 2] }
```

#### `GET /slow`

Запускает тяжёлую задачу в фоновом daemon-потоке.

**Ответ `202`:** `{ "status": "scheduled" }`

#### `GET /wrong`

Демонстрация обработки ошибки (деление на ноль → `500` с JSON).

## Линтинг и форматирование

```bash
just lint
just fix
just format
```

Или:

```bash
uv run ruff check .
uv run ruff check --fix .
uv run ruff format .
```

## Файлы данных (gitignore)

В `.gitignore` исключены `*.db` и `.venv`. При работе локально появляются:

- `users.db`, `passwords.txt` — из `utils.py`
- `test.db` — из `api.py`

## История изменений

См. [CHANGE_LOG.md](CHANGE_LOG.md).
