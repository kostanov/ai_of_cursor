#!/usr/bin/env bash
# Проверка всех эндпоинтов go-server (curl).
# Использование:
#   ./scripts/test-endpoints.sh
#   BASE_URL=http://my-host:8080 ./scripts/test-endpoints.sh

set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"
BASE_URL="${BASE_URL%/}"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

pass=0
fail=0

request() {
  local method="$1"
  local path="$2"
  local body="${3:-}"
  local expected="${4}"

  local url="${BASE_URL}${path}"
  local tmp
  tmp="$(mktemp)"

  local code
  if [[ -n "$body" ]]; then
    code="$(curl -sS -o "$tmp" -w '%{http_code}' -X "$method" \
      -H 'Content-Type: application/json' \
      -d "$body" "$url")"
  else
    code="$(curl -sS -o "$tmp" -w '%{http_code}' -X "$method" "$url")"
  fi

  local response
  response="$(cat "$tmp")"
  rm -f "$tmp"

  if [[ "$code" == "$expected" ]]; then
    echo -e "${GREEN}OK${NC}  $method $path -> $code"
    ((pass++)) || true
  else
    echo -e "${RED}FAIL${NC} $method $path -> $code (ожидался $expected)"
    echo "    body: $response"
    ((fail++)) || true
  fi

  printf '%s' "$response"
}

echo "=== smoke-test go-server ==="
echo "BASE_URL=$BASE_URL"
echo

# 1) POST /adduser
resp="$(request POST /adduser '{"name":"SmokeTest"}' 201)"
if ! echo "$resp" | grep -q '"status":"ok"'; then
  echo -e "${RED}FAIL${NC} POST /adduser: нет status ok в теле"
  ((fail++)) || true
else
  ((pass++)) || true
fi

user_id="$(echo "$resp" | sed -n 's/.*"id":\([0-9]*\).*/\1/p' | head -1)"
if [[ -z "$user_id" ]]; then
  echo -e "${RED}FAIL${NC} не удалось извлечь id из ответа adduser"
  exit 1
fi
echo "    создан user_id=$user_id"
echo

# 2) GET /user/{id}
resp="$(request GET "/user/${user_id}" '' 200)"
if ! echo "$resp" | grep -q '"name":"SmokeTest"'; then
  echo -e "${RED}FAIL${NC} GET /user: имя не совпало"
  ((fail++)) || true
else
  ((pass++)) || true
fi
echo

# 3) GET /user/999999 — not found
request GET /user/999999 '' 404 >/dev/null
echo

# 4) POST /adduser — пустое имя
request POST /adduser '{"name":"   "}' 400 >/dev/null
echo

# 5) GET /activate/{id}
resp="$(request GET "/activate/${user_id}" '' 200)"
if ! echo "$resp" | grep -q '"active"'; then
  echo -e "${RED}FAIL${NC} GET /activate: нет поля active"
  ((fail++)) || true
else
  ((pass++)) || true
fi
echo

# 6) POST /activate/{id}
request POST "/activate/${user_id}" '' 200 >/dev/null
echo

# 7) GET /slow
request GET /slow '' 202 >/dev/null
echo

# 8) GET /wrong
resp="$(request GET /wrong '' 500)"
if ! echo "$resp" | grep -q 'division by zero'; then
  echo -e "${RED}FAIL${NC} GET /wrong: нет detail division by zero"
  ((fail++)) || true
else
  ((pass++)) || true
fi
echo

# 9) 405 — неверный метод
request GET /adduser '' 405 >/dev/null
echo

echo "=== итог: ${pass} проверок OK, ${fail} FAIL ==="
if [[ "$fail" -gt 0 ]]; then
  exit 1
fi
