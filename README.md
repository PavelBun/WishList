# Wishlist API

REST API сервис для создания вишлистов (списков подарков) с публичным доступом по уникальной ссылке и возможностью бронирования подарков.

Стек: Go 1.25.1, PostgreSQL, pgx, Chi, JWT, Swagger, Docker Compose.

## Возможности

- Регистрация и аутентификация пользователей (JWT)
- CRUD операции с вишлистами и подарками
- Публичный доступ к вишлисту по уникальному токену (UUID)
- Бронирование подарков гостями без авторизации
- Защита от CSRF (Go 1.25 CrossOriginProtection)
- Настраиваемые CORS
- Graceful shutdown
- Миграции базы данных
- Юнит-тесты сервисов и HTTP‑обработчиков

## Запуск

```bash
git clone https://github.com/yourname/wishlist-api
cd wishlist-api
cp .env.example .env   # отредактируйте при необходимости
docker-compose up --build
```

Если файл .env отсутствует, используются значения по умолчанию
Сервис запустится на `http://localhost:8080`.  
Swagger UI: `http://localhost:8080/swagger/index.html`.

## Переменные окружения (.env)

| Переменная            | Описание                                 | Значение по умолчанию          |
|-----------------------|------------------------------------------|--------------------------------|
| `DB_HOST`             | Хост PostgreSQL                          | `db`                           |
| `DB_PORT`             | Порт PostgreSQL                          | `5432`                         |
| `DB_USER`             | Пользователь БД                          | `postgres`                     |
| `DB_PASSWORD`         | Пароль БД                                | `postgres`                     |
| `DB_NAME`             | Имя базы данных                          | `wishlist`                     |
| `APP_PORT`            | Порт HTTP-сервера                        | `8080`                         |
| `JWT_SECRET`          | Секретный ключ для JWT                   | **обязательно**                |
| `JWT_EXPIRY_HOURS`    | Время жизни токена в часах               | `24`                           |
| `CORS_ALLOWED_ORIGINS`| Разрешённые источники (через запятую)    | `*` (для разработки)           |

В коде `CORS_ALLOWED_ORIGINS` отключил для удобства

## API Эндпоинты

### Публичные (без авторизации)

#### Регистрация
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepass"
}
```
Ответ: `201 Created` с объектом пользователя.

#### Вход
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepass"
}
```
Ответ: `200 OK` с JWT-токеном.

#### Просмотр вишлиста по токену
```http
GET /public/wishlists/{access_token}
```
`access_token` – UUID, возвращается при создании вишлиста.

#### Бронирование подарка
```http
POST /public/wishlists/{access_token}/items/{item_id}/book
```
При успехе возвращает `204 No Content`.  
Повторное бронирование – ошибка `409 Conflict`.

### Приватные (требуют авторизации)

Добавьте заголовок: `Authorization: Bearer <JWT>`.

#### Вишлисты

| Метод   | Путь                     | Описание                        |
|---------|--------------------------|---------------------------------|
| POST    | /wishlists               | Создать вишлист                 |
| GET     | /wishlists               | Получить все свои вишлисты      |
| GET     | /wishlists/{id}          | Получить вишлист по ID          |
| PUT     | /wishlists/{id}          | Обновить вишлист                |
| DELETE  | /wishlists/{id}          | Удалить вишлист                 |

Пример создания:
```json
{
  "title": "День рождения",
  "description": "Подарки на 25 лет",
  "event_date": "2026-05-20"
}
```

#### Позиции вишлиста

| Метод   | Путь                                  | Описание                        |
|---------|---------------------------------------|---------------------------------|
| POST    | /wishlists/{wishlist_id}/items        | Добавить позицию                |
| GET     | /wishlists/{wishlist_id}/items        | Получить все позиции            |
| GET     | /items/{item_id}                      | Получить позицию по ID          |
| PUT     | /items/{item_id}                      | Обновить позицию                |
| DELETE  | /items/{item_id}                      | Удалить позицию                 |

Пример добавления:
```json
{
  "title": "Книга «Чистый код»",
  "description": "Роберт Мартин",
  "product_link": "https://ozon.ru/...",
  "priority": 5
}
```

## Примеры запросов curl

```bash
# 1. Регистрация
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass"}'

# 2. Вход (получить токен)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass"}'

# 3. Создание вишлиста (замените <токен>)
curl -X POST http://localhost:8080/wishlists \
  -H "Authorization: Bearer <токен>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Новый год","description":"Праздник","event_date":"2026-12-31"}'

# 4. Просмотр публичного вишлиста (используйте access_token из ответа выше)
curl http://localhost:8080/public/wishlists/<access_token>

# 5. Бронирование подарка
curl -X POST http://localhost:8080/public/wishlists/<access_token>/items/<item_id>/book
```

## Разработка

### Структура проекта
```
.
├── cmd/server/main.go         # точка входа
├── internal
│   ├── app                    # сборка зависимостей, запуск, graceful shutdown
│   ├── config                 # конфигурация из переменных окружения
│   ├── db                     # подключение к PostgreSQL (pgxpool)
│   ├── dto                    # объекты передачи данных
│   ├── handlers               # HTTP-обработчики
│   ├── middleware             # CORS, request ID, JWT-авторизация
│   ├── models                 # сущности БД
│   ├── repository             # слой доступа к данным (PostgreSQL)
│   ├── service                # бизнес-логика и интерфейсы
│   └── validator              # валидация входных данных
├── migrations                 # SQL-миграции (golang-migrate)
├── tests
│   ├── service                # юнит-тесты сервисов
│   └── k6                     # сценарий нагрузочного тестирования
├── docs                       # сгенерированная Swagger-документация
├── index.html                 # простой фронтенд для демонстрации
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

### Миграции
Применяются автоматически контейнером `migrate` в docker-compose. Для ручного запуска:
```bash
docker-compose run --rm migrate -path /migrations -database "postgres://..." up
```

### Генерация Swagger-документации
```bash
swag init -g cmd/server/main.go
```
Документация доступна по `/swagger/index.html`.

### Запуск тестов
```bash
# Юнит-тесты
make test

# Линтер
make lint

# Покрытие тестами (отчёт coverage.html)
make test-cover
```

### Makefile команды
- `make up` – запустить приложение в Docker
- `make down` – остановить контейнеры
- `make test` – запустить все тесты
- `make lint` – проверить код линтером
- `make test-cover` – сгенерировать отчёт о покрытии

## Безопасность

- **JWT** – подпись HMAC‑SHA256, явная проверка алгоритма подписи
- **Пароли** – хэширование bcrypt
- **CSRF** – middleware `http.CrossOriginProtection` (Go 1.25+)
- **CORS** – настраиваемый список разрешённых источников
- **Graceful shutdown** – корректное завершение HTTP‑сервера и закрытие соединений с БД

## Нагрузочное тестирование

В папке `tests/k6` находится сценарий `reserve_test.js` для проверки поведения API под нагрузкой.

### Сценарий

1. **Setup** – создаётся тестовый пользователь, выпускается JWT, генерируются **три вишлиста** с двумя подарками в каждом.
2. **Main** – каждый виртуальный пользователь:
   - случайно выбирает один из созданных вишлистов,
   - запрашивает публичную страницу вишлиста (замеряется время ответа),
   - пытается забронировать случайный подарок из этого вишлиста,
   - ожидает случайную паузу от 0.5 до 2.5 секунд.
3. **Teardown** – выводит в лог количество созданных вишлистов.

### Профиль нагрузки

- **Этапы** (`stages`):
  - 10 секунд разогрева до 5 VU,
  - 20 секунд основной нагрузки с 30 VU,
  - 10 секунд плавного завершения до 0 VU.

### Метрики и пороги

- `public_view_duration` – время ответа публичного эндпоинта (95-й перцентиль < 500 мс).
- `booking_success_rate` – доля успешных бронирований (должна быть > 1% за весь тест).
- `http_req_failed` – общая доля ошибочных запросов (< 1%).

### Запуск

В другом терминале запустите приложение 

```bash
cd tests/k6
k6 run reserve_test.js
```