import hashlib
import json
import os
import sqlite3
import threading
from typing import Any

DB_PATH = "users.db"
PASSWORDS_PATH = "passwords.txt"
PBKDF2_ITERATIONS = 100_000

active_users: list[int] = []
_active_users_lock = threading.Lock()


def _ensure_users_table(conn: sqlite3.Connection) -> None:
    conn.execute(
        """
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            tags TEXT NOT NULL
        )
        """
    )


def init_db() -> None:
    """Инициализирует структуру базы данных."""
    with sqlite3.connect(DB_PATH) as conn:
        _ensure_users_table(conn)


def add_user(name: str, tags: list[str] | None = None) -> int:
    """Добавляет пользователя в локальную базу и возвращает id."""
    if not name or not name.strip():
        raise ValueError("name must be a non-empty string")

    prepared_tags = list(tags or [])
    if "new" not in prepared_tags:
        prepared_tags.append("new")

    tags_json = json.dumps(prepared_tags, ensure_ascii=False)

    with sqlite3.connect(DB_PATH) as conn:
        _ensure_users_table(conn)
        cursor = conn.execute(
            "INSERT INTO users (name, tags) VALUES (?, ?)",
            (name.strip(), tags_json),
        )
        return int(cursor.lastrowid)


def store_password(user_id: int, password: str) -> None:
    """Сохраняет PBKDF2-хеш пароля в файл."""
    if not password:
        raise ValueError("password must be non-empty")

    salt = os.urandom(16)
    hash_bytes = hashlib.pbkdf2_hmac(
        "sha256",
        password.encode("utf-8"),
        salt,
        PBKDF2_ITERATIONS,
    )
    record = f"{user_id}:{salt.hex()}:{hash_bytes.hex()}\n"

    with open(PASSWORDS_PATH, "a", encoding="utf-8") as passwords_file:
        passwords_file.write(record)


def set_active(user_id: int) -> None:
    """Потокобезопасно добавляет пользователя в активные, храним последние 5."""
    with _active_users_lock:
        active_users.append(user_id)
        if len(active_users) > 5:
            active_users.pop(0)


def get_active_users_snapshot() -> list[int]:
    """Возвращает потокобезопасный снимок активных пользователей."""
    with _active_users_lock:
        return list(active_users)


def get_user_by_name(name: str) -> dict[str, object] | None:
    """Возвращает пользователя по имени или None, если не найден."""
    with sqlite3.connect(DB_PATH) as conn:
        _ensure_users_table(conn)
        row = conn.execute(
            "SELECT id, name, tags FROM users WHERE name = ?",
            (name,),
        ).fetchone()

    if row is None:
        return None

    uid, uname, tags_json = row
    try:
        parsed_tags = json.loads(tags_json)
        if not isinstance(parsed_tags, list):
            parsed_tags = []
    except json.JSONDecodeError:
        parsed_tags = []

    return {"id": uid, "name": uname, "tags": parsed_tags}


def is_admin(user: dict[str, Any] | None) -> bool:
    """Проверяет, является ли пользователь администратором."""
    if not isinstance(user, dict):
        return False
    role = user.get("role")
    return role == "admin"


def long_running_task(n: int) -> int:
    """Вычисляет сумму квадратов от 0 до n-1."""
    result = 0
    for i in range(n):
        result += i * i
    return result


def run_long_running_task_in_thread(n: int) -> threading.Thread:
    """Запускает длительную задачу в daemon-потоке."""

    def worker() -> None:
        long_running_task(n)

    t = threading.Thread(target=worker, daemon=True)
    t.start()
    return t


def self_test() -> None:
    """Простой smoke-тест основных функций модуля."""
    init_db()
    userid = add_user("testuser", tags=["starter"])
    assert isinstance(userid, int)
    user = get_user_by_name("testuser")
    assert user is not None
    assert user["name"] == "testuser"
    assert "new" in user["tags"]
    assert is_admin({"role": "admin"}) is True
    assert is_admin({"role": "user"}) is False
    assert is_admin(None) is False
    set_active(userid)
    snapshot = get_active_users_snapshot()
    assert userid in snapshot
    store_password(userid, "secretpassword")
    assert os.path.exists(PASSWORDS_PATH)
    print("self_test passed")


if __name__ == "__main__":
    self_test()
    threads: list[threading.Thread] = []
    for i in range(5):
        t = threading.Thread(target=set_active, args=(i,))
        t.start()
        threads.append(t)
    for t in threads:
        t.join()
    print("Active users snapshot:", get_active_users_snapshot())
