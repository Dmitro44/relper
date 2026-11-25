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

places=(
  "$MINSK|Кавярня на Немиге|food|ул. Немига 5|53.9035|27.5540|Кофе и завтраки рядом с Немигой.|daily_life|Загляните утром|До полудня меньше людей и свежая выпечка."
  "$MINSK|Парк Горького|hobby|ул. Фрунзе|53.9270|27.5680|Большой парк с набережной и прокатом.|daily_life|Возьмите воду и удобную обувь|Дорожки длинные; в жару лучше идти утром или вечером."
  "$MINSK|Национальная библиотека|hobby|пр. Независимости 116|53.9315|27.6460|Библиотека с обзорной площадкой и тихими залами.|integration|Запишитесь заранее|Для залов и обзорной площадки нужен читательский билет."
  "$MINSK|Комаровский рынок|shopping|ул. Веры Хоружей 8|53.9170|27.5830|Рынок с продуктами, специями и ремёслами.|daily_life|Приходите до полудня|Утром больше выбора и короче очереди."
  "$MINSK|Площадь Ленина|transport|пл. Ленина|53.8930|27.5480|Узел метро и наземного транспорта в центре.|daily_life|Сверьтесь с указателями выходов|Выходы метро ведут в разные стороны площади."
  "$MINSK|Галерея Минск|shopping|пр. Победителей 9|53.9080|27.5490|Торговый центр рядом с центром города.|daily_life|Отдохните на верхнем этаже|Там меньше людей и есть зоны отдыха."
  "$MINSK|Минск-Пассажирский|transport|пл. Привокзальная|53.8907|27.5506|Железнодорожный вокзал и выход к метро.|documents|Возьмите паспорт|Для междугородних билетов нужен документ."
  "$MINSK|Парк Челюскинцев|hobby|ул. Челюскинцев|53.9250|27.6150|Парк с аттракционами, катком и дорожками.|daily_life|Уточните работу аттракционов|График зависит от сезона и погоды."
  "$MINSK|Кофейня Zerno|food|ул. Октябрьская 16|53.9100|27.5600|Кофе и десерты в центре Минска.|daily_life|Спросите зерно дня|Бариста подскажет сорт под ваш вкус."
  "$MINSK|Поликлиника №1|healthcare|ул. Кирова 2|53.8950|27.5550|Городская поликлиника рядом с центром.|documents|Возьмите полис и паспорт|На первом приёме нужны документы и направление."
)

place_ids=()
for row in "${places[@]}"; do
  IFS='|' read -r location name category address lat lng description tip_category tip_title tip_content <<< "$row"
  place_id="$(create_place "$location" "$name" "$category" "$address" "$lat" "$lng" "$description")"
  place_ids+=("$place_id")
  tip "$place_id" "$TOKEN_DEMO" "$tip_category" "$tip_title" "$tip_content"
done

comments=(
  "Отличное место, вернусь ещё."
  "Удобно добираться, рядом всё нужное."
  "Спокойно и без лишнего шума."
  "Хороший вариант для первого знакомства с районом."
  "Персонал помог разобраться с деталями."
)

for i in "${!place_ids[@]}"; do
  place_id="${place_ids[$i]}"
  rate "$place_id" "$TOKEN_ANNA" "$((5 - (i % 2)))"
  rate "$place_id" "$TOKEN_MAKS" "4"
  rate "$place_id" "$TOKEN_LENA" "$((3 + (i % 3)))"
  comment "$place_id" "$TOKEN_ANNA" "${comments[$((i % ${#comments[@]}))]}"
done


echo "Demo admin: $DEMO_EMAIL / $DEMO_PASSWORD"
echo "Reviewers: anna@relper.local, maks@relper.local, lena@relper.local / $DEMO_PASSWORD"
echo "Seeded 1 location, ${#place_ids[@]} places, ratings, comments and tips."
