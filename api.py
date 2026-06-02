import sqlite3
from http import HTTPStatus
import threading
import time
from typing import Any

from flask import Flask, jsonify, g, request

DB_PATH = "test.db"

app = Flask(__name__)


def get_db() -> sqlite3.Connection:
    """
    Возвращает соединение с БД, привязанное к контексту запроса.
    Соединение автоматически закрывается в конце запроса.
    """
    if "db" not in g:
        conn = sqlite3.connect(DB_PATH)
        conn.row_factory = sqlite3.Row
        g.db = conn
    return g.db


@app.teardown_appcontext
def close_db(_exc: BaseException | None) -> None:  # type: ignore[override]
    db = g.pop("db", None)
    if db is not None:
        db.close()


def init_db() -> None:
    """
    Инициализирует таблицу пользователей при первом запуске.
    """
    conn = sqlite3.connect(DB_PATH)
    try:
        cur = conn.cursor()
        cur.execute(
            "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT NOT NULL)"
        )
        conn.commit()
    finally:
        conn.close()


@app.route("/adduser", methods=["POST"])
def add_user() -> Any:
    """
    Создаёт пользователя.

    Ожидает JSON:
    {
        "name": "<строка>"
    }
    """
    payload = request.get_json(silent=True) or {}
    name = payload.get("name")

    if not isinstance(name, str) or not name.strip():
        return jsonify({"error": "name is required"}), HTTPStatus.BAD_REQUEST

    db = get_db()
    cur = db.cursor()
    cur.execute("INSERT INTO users (name) VALUES (?)", (name.strip(),))
    db.commit()

    user_id = cur.lastrowid
    return (
        jsonify({"status": "ok", "id": int(user_id), "name": name.strip()}),
        HTTPStatus.CREATED,
    )


@app.route("/user/<int:uid>", methods=["GET"])
def get_user(uid: int) -> Any:
    """
    Возвращает пользователя по id.
    """
    db = get_db()
    cur = db.cursor()
    cur.execute("SELECT id, name FROM users WHERE id = ?", (uid,))
    row = cur.fetchone()

    if row is None:
        return jsonify({"error": "not_found"}), HTTPStatus.NOT_FOUND

    return jsonify({"id": row["id"], "name": row["name"]}), HTTPStatus.OK


# Менеджмент активных пользователей (без гонок)
_active_users: list[int] = []
_active_lock = threading.Lock()


def add_active_user(user_id: int, max_size: int = 100) -> None:
    """
    Добавляет пользователя в список активных (ограничение длины).
    Потокобезопасно.
    """
    with _active_lock:
        _active_users.append(user_id)
        if len(_active_users) > max_size:
            overflow = len(_active_users) - max_size
            del _active_users[:overflow]


def get_active_users() -> list[int]:
    """
    Возвращает копию списка активных пользователей.
    """
    with _active_lock:
        return list(_active_users)


@app.route("/activate/<int:uid>", methods=["POST", "GET"])
def activate(uid: int) -> Any:
    """
    Помечает пользователя активным и возвращает актуальный список.
    Делается синхронно, без лишних потоков и гонок.
    """
    # Небольшая имитация задержки, если нужно показать "работу"
    time.sleep(0.1)
    add_active_user(uid)
    return jsonify({"status": "ok", "active": get_active_users()}), HTTPStatus.OK


def _slow_task(iterations: int = 200_000) -> int:
    """
    Тяжёлая CPU-задача, вынесенная в отдельную функцию.
    """
    x = 0
    for i in range(iterations):
        x += i * (i + 1) // 2
    return x


@app.route("/slow", methods=["GET"])
def slow() -> Any:
    """
    Демонстрационный эндпоинт: запускает тяжёлую задачу в отдельном потоке
    и сразу возвращает ответ, чтобы не блокировать сервер.
    """

    def worker() -> None:
        _slow_task()

    threading.Thread(target=worker, daemon=True).start()
    return jsonify({"status": "scheduled"}), HTTPStatus.ACCEPTED


@app.route("/wrong", methods=["GET"])
def wrong() -> Any:
    """
    Пример корректной обработки ошибки.
    """
    try:
        _ = 10 / 0
    except ZeroDivisionError:
        return (
            jsonify({"msg": "error", "detail": "division by zero"}),
            HTTPStatus.INTERNAL_SERVER_ERROR,
        )

    # Этот код фактически недостижим, но оставлен как пример "успеха"
    return jsonify({"msg": "ok"}), HTTPStatus.OK


init_db()

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080, debug=True)
