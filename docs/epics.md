# Эпик 1: БД и миграции

- Настроить `golang-migrate` и миграции (см. SQL выше).
- Включить расширения `uuid-ossp`, `pg_trgm`.
- Сид-user для локалки.

# Эпик 2: Бэкенд каркас

- Конфиг (`env`): `PORT`, `DATABASE_URL`, `JWT_SECRET`, `OPENAI_BASE_URL`, `OPENAI_API_KEY`.
- Роутер `chi`, middleware (logger, recover, CORS, request-id, rate-limit).
- Подключить БД, генерация sqlc.

# Эпик 3: Аутентификация

- `POST /auth/register`, `POST /auth/login`, `POST /auth/refresh`, `POST /auth/logout`.
- Хеширование паролей (`bcrypt`).
- JWT + refresh cookie HttpOnly + ротация.

# Эпик 4: Посты/Комментарии/Лайки

- CRUD постов.
- Комментарии.
- Лайки.

# Эпик 5: Подписки и лента

- Follow/Unfollow.
- `GET /feed`.

# Эпик 6: Поиск

- `GET /search` по тексту/хештегам.
- Выделение хэштегов при создании поста (бэкенд: парсер).

# Эпик 7: AI-интеграция

- `POST /ai/generate` (сервер → gpt-oss-120b).
- UI вкладка «Сгенерировать».

# Эпик 8: Frontend

- Каркас Vite + TS + React Router + TanStack Query + Tailwind + Zustand.
- Страницы `/feed`, `/profile/:username`, `/post/:id`, `/compose`, `/search`, `/login|/register`.
- Тёмная/светлая тема.
- Лоадеры/скелетоны.

# Эпик 9: Notifications (MVP)

- Создание нотификаций при лайке/комменте/фоллоу.
- `GET /notifications?unread_only=true` (polling).

# Эпик 10: Докер/деплой

- `docker-compose.yml` (db, migrate, api, web, nginx).
- Nginx конфиг (SPA + прокси на API).
- README с инструкциями.

# Эпик 11: Тесты

- Go: unit/handler tests (`httptest`), репозитории (`sqlmock` или `testcontainers`).
- FE: Vitest + RTL (реквесты мокируются MSW).