default:
    @just --list

run:
    uv run python main.py

api:
    uv run python api.py

lint:
    uv run ruff check .

fix:
    uv run ruff check --fix .

format:
    uv run ruff format .
