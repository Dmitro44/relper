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
  local category="$1"
  local index="$2"
  local variant=$((index % 4))
  case "$category" in
    food)
      case "$variant" in
        0) printf '%s' "Брала кофе и сырники, всё принесли быстро. Утром здесь спокойно, а к обеду уже шумно." ;;
        1) printf '%s' "Взяла суп и салат, вышло сытно. В обед свободных столов почти нет, лучше приходить чуть раньше." ;;
        2) printf '%s' "Кофе вкусный, но зал небольшой. Если хотите сесть у окна, приходите до полудня." ;;
        3) printf '%s' "Вкусно, но в обед шумно." ;;
      esac ;;
    hobby)
      case "$variant" in
        0) printf '%s' "Провела здесь пару часов и не заметила. В выходные много людей, в будни спокойнее." ;;
        1) printf '%s' "Место интересное, но лучше без спешки. Утром народу меньше, можно всё посмотреть." ;;
        2) printf '%s' "Пришла ближе к вечеру, когда основная толпа уже разошлась. Можно спокойно посидеть и передохнуть." ;;
        3) printf '%s' "Хорошо, что можно просто посидеть." ;;
      esac ;;
    shopping)
      case "$variant" in
        0) printf '%s' "Выбор большой, но цены не всегда ниже, чем в других местах. В будни здесь свободнее и с парковкой проще." ;;
        1) printf '%s' "Нашла всё, что искала, за один заход. В выходные у касс очередь, лучше приехать вечером в будни." ;;
        2) printf '%s' "Магазинов много, можно провести здесь полдня. Если идёте за чем-то конкретным, проверьте этаж заранее." ;;
        3) printf '%s' "Выбор нормальный, но в выходные тесно." ;;
      esac ;;
    transport)
      case "$variant" in
        0) printf '%s' "В час пик здесь тесно, зато пересадки удобные. С чемоданом лучше не ехать в утренний пик." ;;
        1) printf '%s' "Станция нормальная, но выходы легко перепутать. Посмотрите нужный номер заранее." ;;
        2) printf '%s' "Добралась быстро, поезда ходят часто. В вечерний пик вагоны набиваются, но ждать долго не приходится." ;;
        3) printf '%s' "В час пик здесь тесно, зато удобно." ;;
      esac ;;
    healthcare)
      case "$variant" in
        0) printf '%s' "Приняли по записи, ждали недолго. В регистратуре попросили паспорт и полис, лучше взять оба." ;;
        1) printf '%s' "Врач всё объяснил, но к кабинету была очередь. Приходите минут за десять до приёма." ;;
        2) printf '%s' "Записалась заранее, поэтому всё прошло без суеты. Если без записи, можно просидеть в очереди полдня." ;;
        3) printf '%s' "Приняли по записи, всё прошло спокойно." ;;
      esac ;;
    services)
      case "$variant" in
        0) printf '%s' "Вопрос решили за один визит, но пришлось подождать. В обед очередь короче, чем утром." ;;
        1) printf '%s' "Сотрудник подробно объяснил, какие документы нужны. Без паспорта не принимают, проверьте заранее." ;;
        2) printf '%s' "Народу было немного, обслужили быстро. Если идёте с документами, лучше взять копии." ;;
        3) printf '%s' "Вопрос решили, но пришлось подождать." ;;
      esac ;;
    *) printf '%s' "Место нормальное, но лучше приходить в будни. В выходные здесь больше людей." ;;
  esac
}

comment_two_for() {
  local category="$1"
  local index="$2"
  local variant=$((index % 4))
  case "$category" in
    food)
      case "$variant" in
        0) printf '%s' "Зашла по пути, взяла кофе и десерт. Внутри уютно, но столики у окна быстро занимают." ;;
        1) printf '%s' "Была здесь в обед, свободный столик нашёлся не сразу. Кормят хорошо, но в час пик лучше не приходить." ;;
        2) printf '%s' "Взяла только кофе, но задержалась из-за очереди. Бариста вежливые, так что ждать было не обидно." ;;
        3) printf '%s' "Уютно, но столиков мало." ;;
      esac ;;
    hobby)
      case "$variant" in
        0) printf '%s' "Была в будний день, народу почти нет. Можно спокойно всё посмотреть и не торопиться." ;;
        1) printf '%s' "Пришла сюда на пару часов, ушла довольная. В выходные лучше приходить к открытию." ;;
        2) printf '%s' "Внутри интересно, но информации много. Если хотите всё прочитать, закладывайте больше времени." ;;
        3) printf '%s' "Спокойно и без толпы." ;;
      esac ;;
    shopping)
      case "$variant" in
        0) printf '%s' "Зашла за подарком и нашла нужное. В будни продавцы свободнее и помогают быстрее." ;;
        1) printf '%s' "Магазинов хватает, но легко потеряться. Смотрите указатели у входа." ;;
        2) printf '%s' "Была здесь вечером, народу немного. Если едете за конкретным, лучше проверить наличие заранее." ;;
        3) printf '%s' "Много магазинов, есть где походить." ;;
      esac ;;
    transport)
      case "$variant" in
        0) printf '%s' "Ехала с сумкой, пересадка заняла пару минут. В час пик лучше выбирать другой выход." ;;
        1) printf '%s' "Станция чистая, но указатели не всегда понятные. Спросите у сотрудника, если сомневаетесь." ;;
        2) printf '%s' "Добралась без проблем, поезд пришёл быстро. Вечером здесь людно, но терпимо." ;;
        3) printf '%s' "Удобно, если знаешь нужный выход." ;;
      esac ;;
    healthcare)
      case "$variant" in
        0) printf '%s' "Была на приёме утром, очереди почти не было. Всё объяснили и записали на следующий раз." ;;
        1) printf '%s' "Записалась через регистратуру, ждала недолго. В коридоре есть где сесть." ;;
        2) printf '%s' "Врач нормальный, но пришлось подождать. Лучше приходить заранее." ;;
        3) printf '%s' "Всё прошло спокойно." ;;
      esac ;;
    services)
      case "$variant" in
        0) printf '%s' "Пришла с документами, вопрос решили за полчаса. Очередь была, но двигалась быстро." ;;
        1) printf '%s' "Сотрудник помог разобраться, что нужно. Если идёте впервые, возьмите паспорт." ;;
        2) printf '%s' "Была здесь в обед, народу немного. Всё сделали без лишних вопросов." ;;
        3) printf '%s' "Быстро и без нервов." ;;
      esac ;;
    *) printf '%s' "Пришла в будни, народу немного. В выходные здесь обычно больше людей." ;;
  esac
}

second_tip_category() {
  case "$1" in
    healthcare|services) printf '%s' "documents" ;;
    *) printf '%s' "daily_life" ;;
  esac
}

second_tip_title() {
  case "$1" in
    food) printf '%s' "Выберите время без очереди" ;;
    hobby) printf '%s' "Заложите больше времени" ;;
    shopping) printf '%s' "Посмотрите план этажей" ;;
    transport) printf '%s' "Проверьте маршрут заранее" ;;
    healthcare) printf '%s' "Возьмите документы" ;;
    services) printf '%s' "Уточните документы" ;;
    *) printf '%s' "Выберите будний день" ;;
  esac
}

second_tip_content() {
  local category="$1"
  local index="$2"
  local variant=$((index % 4))
  case "$category" in
    food)
      case "$variant" in
        0) printf '%s' "Берите обед до полудня. После часа свободных столов почти не остаётся." ;;
        1) printf '%s' "Спросите, что сегодня свежее. Иногда это проще, чем выбирать по меню." ;;
        2) printf '%s' "Если хотите тихо посидеть, приходите утром. К обеду здесь становится шумно." ;;
        3) printf '%s' "Утром здесь спокойнее." ;;
      esac ;;
    hobby)
      case "$variant" in
        0) printf '%s' "Планируйте пару часов, иначе не успеете всё посмотреть. В будни народу меньше." ;;
        1) printf '%s' "Возьмите с собой воду, особенно летом. Внутри бывает душно." ;;
        2) printf '%s' "Проверьте расписание заранее. В выходные может быть очередь." ;;
        3) printf '%s' "Лучше приходить утром." ;;
      esac ;;
    shopping)
      case "$variant" in
        0) printf '%s' "Смотрите план этажей у входа. Так проще найти нужный магазин." ;;
        1) printf '%s' "В будни здесь свободнее. В выходные у касс очередь." ;;
        2) printf '%s' "Если едете на машине, приезжайте до обеда. Парковка быстро заполняется." ;;
        3) printf '%s' "В будни проще с парковкой." ;;
      esac ;;
    transport)
      case "$variant" in
        0) printf '%s' "В час пик лучше выехать чуть раньше. На платформе тесно, но поезда ходят часто." ;;
        1) printf '%s' "В час пик берите запас времени. Народу много, но поезда ходят часто." ;;
        2) printf '%s' "Если с чемоданом, выбирайте выход с лифтом. С лестницей будет неудобно." ;;
        3) printf '%s' "Проверьте нужный выход." ;;
      esac ;;
    healthcare)
      case "$variant" in
        0) printf '%s' "Возьмите паспорт и полис. Без них могут отправить обратно." ;;
        1) printf '%s' "Записывайтесь заранее. Утром очередь обычно короче." ;;
        2) printf '%s' "Приходите минут за десять до приёма. Так меньше суеты." ;;
        3) printf '%s' "Возьмите документы." ;;
      esac ;;
    services)
      case "$variant" in
        0) printf '%s' "Уточните список документов заранее. Так меньше шансов прийти дважды." ;;
        1) printf '%s' "В обед очередь короче. Утром обычно больше людей." ;;
        2) printf '%s' "Возьмите копии документов. Иногда их просят отдельно." ;;
        3) printf '%s' "Проверьте документы заранее." ;;
      esac ;;
    *) printf '%s' "Приходите в будни. В выходные здесь больше людей." ;;
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
  idx=${#place_ids[@]}
  place_id="$(create_place "$MINSK" "$name" "$category" "$address" "$lat" "$lng" "$description")"
  place_ids+=("$place_id")
  tip "$place_id" "$TOKEN_DEMO" "$tip_category" "$tip_title" "$tip_content"
  if [ "$name" = "БГУИР, 4-й корпус" ]; then
    tip "$place_id" "$TOKEN_MAKS" "daily_life" "Уточните кабинет заранее" "Заранее уточните, в каком кабинете проходят занятия. В корпусе несколько этажей, и указатели не всегда заметны."
    comment "$place_id" "$TOKEN_ANNA" "Удобно, что корпус рядом с метро."
    comment "$place_id" "$TOKEN_LENA" "Внутри обычные аудитории, на этажах есть указатели. На входе могут попросить пропуск, лучше уточнить заранее."
  else
    tip "$place_id" "$TOKEN_MAKS" "$(second_tip_category "$category")" "$(second_tip_title "$category")" "$(second_tip_content "$category" "$idx")"
    comment "$place_id" "$TOKEN_ANNA" "$(comment_for "$category" "$idx")"
    comment "$place_id" "$TOKEN_LENA" "$(comment_two_for "$category" "$idx")"
  fi
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
