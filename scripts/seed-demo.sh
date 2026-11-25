#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
BASE="${BASE:-http://localhost:3000}"
DEMO_EMAIL="demo@relper.local"
DEMO_PASSWORD="demo12345"

existing="$(curl -sS "$BASE/locations")"
if [ "$(printf '%s' "$existing" | jq length)" -gt 0 ]; then
  echo "Demo data already exists; nothing to do."
  exit 0
fi

login_token() {
  local email="$1"
  local payload response
  payload="$(jq -n --arg email "$email" --arg password "$DEMO_PASSWORD" '{email:$email,password:$password}')"
  response="$(curl -sS -X POST "$BASE/auth/login" -H 'Content-Type: application/json' -d "$payload")"
  printf '%s' "$response" | jq -r '.access_token // empty'
}

ensure_token() {
  local email="$1"
  local nickname="$2"
  local payload response token
  payload="$(jq -n --arg email "$email" --arg password "$DEMO_PASSWORD" --arg nickname "$nickname" '{email:$email,password:$password,nickname:$nickname}')"
  response="$(curl -sS -X POST "$BASE/auth/register" -H 'Content-Type: application/json' -d "$payload")"
  token="$(printf '%s' "$response" | jq -r '.access_token // empty')"
  if [ -z "$token" ]; then
    token="$(login_token "$email")"
  fi
  if [ -z "$token" ]; then
    echo "Cannot authenticate $email: $response" >&2
    exit 1
  fi
  printf '%s' "$token"
}

create_location() {
  local country="$1"
  local city="$2"
  local lat="$3"
  local lng="$4"
  local payload
  payload="$(jq -n --arg country "$country" --arg city "$city" --argjson lat "$lat" --argjson lng "$lng" '{country:$country,city:$city,coordinates:{lat:$lat,lng:$lng}}')"
  curl -sS -X POST "$BASE/locations" \
    -H "Authorization: Bearer $TOKEN_DEMO" \
    -H 'Content-Type: application/json' \
    -d "$payload" | jq -r '.id'
}

create_place() {
  local location="$1"
  local name="$2"
  local category="$3"
  local address="$4"
  local lat="$5"
  local lng="$6"
  local description="$7"
  local payload
  payload="$(jq -n \
    --arg location "$location" \
    --arg name "$name" \
    --arg category "$category" \
    --arg address "$address" \
    --arg description "$description" \
    --argjson lat "$lat" \
    --argjson lng "$lng" \
    '{location_id:$location,name:$name,category:$category,address:$address,coordinates:{lat:$lat,lng:$lng},description:$description}')"
  curl -sS -X POST "$BASE/places" \
    -H "Authorization: Bearer $TOKEN_DEMO" \
    -H 'Content-Type: application/json' \
    -d "$payload" | jq -r '.id'
}

rate() {
  local place="$1"
  local token="$2"
  local stars="$3"
  local payload
  payload="$(jq -n --argjson stars "$stars" '{stars:$stars}')"
  curl -sS -X PUT "$BASE/places/$place/rating" \
    -H "Authorization: Bearer $token" \
    -H 'Content-Type: application/json' \
    -d "$payload" >/dev/null
}

comment() {
  local place="$1"
  local token="$2"
  local content="$3"
  local payload
  payload="$(jq -n --arg place "$place" --arg content "$content" '{place_id:$place,content:$content}')"
  curl -sS -X POST "$BASE/comments" \
    -H "Authorization: Bearer $token" \
    -H 'Content-Type: application/json' \
    -d "$payload" >/dev/null
}

tip() {
  local place="$1"
  local token="$2"
  local category="$3"
  local title="$4"
  local content="$5"
  local payload
  payload="$(jq -n \
    --arg place "$place" \
    --arg category "$category" \
    --arg title "$title" \
    --arg content "$content" \
    '{place_id:$place,category:$category,title:$title,content:$content}')"
  curl -sS -X POST "$BASE/tips" \
    -H "Authorization: Bearer $token" \
    -H 'Content-Type: application/json' \
    -d "$payload" >/dev/null
}

comment_for() {
  case "$1" in
    food) printf '%s' "Вкусно, порции нормальные." ;;
    hobby) printf '%s' "Интересное место, стоит зайти." ;;
    shopping) printf '%s' "Выбор большой, можно найти нужное." ;;
    transport) printf '%s' "В час пик шумно, зато перебираться удобно." ;;
    healthcare) printf '%s' "Приняли по записи, всё организовано." ;;
    services) printf '%s' "Быстро помогли с вопросом." ;;
    *) printf '%s' "Хорошее место." ;;
  esac
}

TOKEN_DEMO="$(ensure_token "$DEMO_EMAIL" "Демо")"
TOKEN_ANNA="$(ensure_token "anna@relper.local" "Анна")"
TOKEN_MAKS="$(ensure_token "maks@relper.local" "Макс")"
TOKEN_LENA="$(ensure_token "lena@relper.local" "Лена")"

docker compose exec -T postgres psql -U relper -d relper -q -c \
  "insert into user_roles (user_id, role_id) select u.id, r.id from users u, roles r where u.email = '$DEMO_EMAIL' and r.code = 'admin' on conflict do nothing;" >/dev/null

TOKEN_DEMO="$(login_token "$DEMO_EMAIL")"
if [ -z "$TOKEN_DEMO" ]; then
  echo "Cannot log in demo admin" >&2
  exit 1
fi

MINSK="$(create_location "Belarus" "Minsk" "53.9006" "27.5590")"

PLACES_FILE="$ROOT/scripts/demo-places.tsv"
place_ids=()
while IFS='|' read -r name category address lat lng description tip_category tip_title tip_content; do
  case "$name" in
    ''|'#'*) continue ;;
  esac
  place_id="$(create_place "$MINSK" "$name" "$category" "$address" "$lat" "$lng" "$description")"
  place_ids+=("$place_id")
  tip "$place_id" "$TOKEN_DEMO" "$tip_category" "$tip_title" "$tip_content"
  comment "$place_id" "$TOKEN_ANNA" "$(comment_for "$category")"
done < "$PLACES_FILE"

for i in "${!place_ids[@]}"; do
  place_id="${place_ids[$i]}"
  rate "$place_id" "$TOKEN_ANNA" "$((5 - (i % 2)))"
  rate "$place_id" "$TOKEN_MAKS" "4"
  rate "$place_id" "$TOKEN_LENA" "$((3 + (i % 3)))"
done


echo "Demo admin: $DEMO_EMAIL / $DEMO_PASSWORD"
echo "Reviewers: anna@relper.local, maks@relper.local, lena@relper.local / $DEMO_PASSWORD"
echo "Seeded 1 location, ${#place_ids[@]} places, ratings, comments and tips."
