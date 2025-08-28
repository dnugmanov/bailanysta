
# 1. Функциональные требования (FRD)

## 1.1. Роли и аутентификация

- **Регистрация/логин:** email + пароль
- **Профиль пользователя:** аватар (опционально), имя, био, навыки/интересы
- **JWT:** access + refresh, хранение refresh в HttpOnly cookie

## 1.2. Посты

- **CRUD постов:** создать, читать, редактировать, удалять (автором)
- **Содержимое:** текст (обязательно), теги/хештеги (извлекаются автоматически), привязка к курсу/модулю (опционально)
- **Лайки, комментарии**
- **Поиск:** по ключевым словам и #хештегам

## 1.3. Лента

- **Главная лента:** посты авторов, на которых подписан, + персональный микс (рекомендуемые/популярные)
- **Лента курса/модуля:** фильтр по `course_id`/`module_id`
- **Пагинация:** limit/offset или cursor-pagination

## 1.4. Подписки и уведомления

- **Подписка на авторов**
- **Уведомления:** новый комментарий/лайк/подписчик (MVP — в приложении, push позже)
- **Канал обновлений:** SSE или WebSocket (MVP — polling)

## 1.5. Темы и UI/UX

- **Светлая/тёмная тема** + сохранение выбора (localStorage)
- **Лоадеры/скелетоны**
- **Роутинг:** `/feed`, `/profile/:username`, `/post/:id`, `/compose`, `/search`

## 1.6. AI‑интеграция (OpenAI‑совместимый gpt-oss-120b)

- **Серверный эндпоинт:** `/ai/generate` (на вход — промпт/параметры, на выход — сгенерированный текст)
- **Возможность "Создать пост по промпту":** UI-форма → вызывает бэкенд → черновик

## 1.7. Нефункциональные требования

- **Производительность:** P95 < 200–300 мс для основных GET (на локальном/дев окружении — реалистично)
- **Безопасность:** bcrypt/argon2, валидация входных данных, CORS, rate-limit
- **Логи и трассировка:** structured logs
- **Миграции БД, сиды**
- **Докеризация:** docker-compose, README, публичный GitHub

---

# 2. Архитектура и стек

## Backend (Go)

- **HTTP:** chi или fiber (выбрано chi за простоту и совместимость со средой тестов)
- **Слои:** handlers (transport) → services (usecases) → repositories (db/sqlc)
- **Валидация:** go-playground/validator
- **БД:** PostgreSQL + sqlc (типобезопасные запросы) или GORM (выбрано sqlc)
- **Миграции:** golang-migrate
- **Авторизация:** JWT (golang-jwt/jwt/v5), refresh‑cookies HttpOnly
- **AI:** клиент OpenAI‑совместимого API (например, [github.com/sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) с BaseURL)
- **Поиск:** pg_trgm + GIN индекс, простая обработка хэштегов

## Frontend (React + TypeScript)

- **Сборка:** Vite
- **Роутинг:** React Router
- **Данные:** TanStack Query (react-query)
- **Состояние UI:** Zustand (тема, временные флаги)
- **UI:** Tailwind CSS (быстрая тёмная тема), Headless UI/Radix (по желанию)
- **Тесты:** Vitest + Testing Library

## Deploy

- **docker-compose:** api, web, db, migrate, (опц.) pgadmin
- **Nginx:** для статики web и реверс‑прокси на api
- **ENV:** через `.env` и `.env.example`

---

# 3. Модель данных (ERD) и миграции

## Сущности

- **users:** id, username, email, password_hash, bio, avatar_url, created_at
- **posts:** id, author_id, text, course_id, module_id, created_at, updated_at
- **likes:** user_id, post_id, created_at
- **comments:** id, post_id, author_id, text, created_at
- **follows:** follower_id, followee_id, created_at
- **courses:** id, title, description
- **modules:** id, course_id, title, order
- **hashtags:** id, tag
- **post_hashtags:** post_id, hashtag_id
- **notifications:** id, user_id, type, entity_id, payload_json, read_at, created_at

## Быстрые миграции (черновик)
Быстрые миграции (черновик):

```
-- 0001_init.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  bio TEXT DEFAULT '',
  avatar_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE courses (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  title TEXT NOT NULL,
  description TEXT DEFAULT ''
);

CREATE TABLE modules (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  "order" INT NOT NULL DEFAULT 0
);

CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  author_id UUID REFERENCES users(id) ON DELETE CASCADE,
  text TEXT NOT NULL,
  course_id UUID REFERENCES courses(id),
  module_id UUID REFERENCES modules(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE comments (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
  author_id UUID REFERENCES users(id) ON DELETE CASCADE,
  text TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE likes (
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, post_id)
);

CREATE TABLE follows (
  follower_id UUID REFERENCES users(id) ON DELETE CASCADE,
  followee_id UUID REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (follower_id, followee_id)
);

CREATE TABLE hashtags (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  tag TEXT UNIQUE NOT NULL
);

CREATE TABLE post_hashtags (
  post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
  hashtag_id UUID REFERENCES hashtags(id) ON DELETE CASCADE,
  PRIMARY KEY (post_id, hashtag_id)
);

CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  type TEXT NOT NULL,              -- like|comment|follow
  entity_id UUID,                  -- пост/коммент/пользователь
  payload_json JSONB NOT NULL,
  read_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

-- Поиск
CREATE INDEX posts_text_trgm_idx ON posts USING GIN (text gin_trgm_ops);
# 📋 Техническая спецификация: Социальная платформа с AI

## 🔍 **База данных - Поиск**

```sql
-- Полнотекстовый поиск с триграммами
CREATE INDEX posts_text_trgm_idx ON posts USING GIN (text gin_trgm_ops);
```

---

## 🚀 **API спецификация v1**

### Общие положения
- **Версия:** `/api/v1`
- **Формат:** JSON, `Content-Type: application/json`
- **Авторизация:** `Authorization: Bearer <access_token>`
- **Пагинация:** `?limit=20&offset=0`
- **Ошибки:** 
```json
{ 
  "error": { 
    "code": "BAD_REQUEST", 
    "message": "..." 
  } 
}
```

### 🔐 **Auth**
- `POST /auth/register` — `{email, username, password}` → `201 {user, tokens}`
- `POST /auth/login` — `{email, password}` → `200 {user, tokens}`
- `POST /auth/refresh` — refresh cookie → `200 {access_token}`
- `POST /auth/logout` — инвалидация refresh

### 👤 **Users**
- `GET /me` → профиль текущего пользователя
- `GET /users/:id` → профиль пользователя
- `PATCH /me` — обновление био/аватара
- `POST /users/:id/follow` / `DELETE /users/:id/follow`

### 📝 **Posts**
- `POST /posts` — создать пост `{text, course_id?, module_id?}`
- `GET /posts/:id` — пост + агрегаты `{like_count, comment_count}`
- `PATCH /posts/:id` — изменить свой пост
- `DELETE /posts/:id`
- `POST /posts/:id/like` / `DELETE /posts/:id/like`
- `GET /posts/:id/comments` — список комментариев
- `POST /posts/:id/comments` — добавить комментарий

### 📰 **Feed/Search**
- `GET /feed` — лента подписок (плюс популярные)
- `GET /search?query=...` — посты и/или пользователи; поддержка `#tag`

### 📚 **Courses/Modules**
- `GET /courses` | `GET /courses/:id/modules`
- *(опц.)* `POST /courses` / `POST /modules` для админки

### 🤖 **AI**
- `POST /ai/generate` — `{prompt, max_tokens?, temperature?}` → `{text}`
- Вызов OpenAI‑совместимого API (gpt‑oss‑120b) с серверной стороны

### 🔔 **Notifications**
- `GET /notifications?unread_only=true`
- *(опц.)* `GET /notifications/stream` — SSE

---

## 📋 **Примеры контрактов**

### POST /posts (создание поста)
```json
// Запрос
{ 
  "text": "Мой конспект по #transformers", 
  "course_id": "uuid?", 
  "module_id": "uuid?" 
}

// Ответ
{
  "id": "uuid",
  "author_id": "uuid",
  "text": "Мой конспект по #transformers",
  "created_at": "2025-08-28T12:00:00Z",
  "like_count": 0,
  "comment_count": 0,
  "hashtags": ["transformers"]
}
```

### POST /ai/generate
```json
// Запрос
{ 
  "prompt": "Сделай конспект главы про Self-Attention", 
  "max_tokens": 300 
}

// Ответ
{ 
  "text": "Self-Attention — это механизм, при котором..." 
}
```

---

## 🎨 **UI спецификация**

### Роуты
- `/feed` — главная лента
- `/compose` — создание поста (есть вкладка «Сгенерировать с AI»)
- `/profile/:username` — профиль, список постов, кнопка Follow/Unfollow
- `/post/:id` — страница поста с комментариями
- `/search?q=...` — результаты поиска
- `/login`, `/register`, `/settings`

### Главные экраны и состояния
- **Лента:** skeleton при загрузке, EmptyState без подписок (предложить подписаться)
- **Профиль:** аватар/имя/био, табы: Posts / About
- **Создание поста:** textarea + подсветка #тегов; таб «AI»: поле prompt → прогресс → предпросмотр → «Вставить в черновик»
- **Тема:** тёмная/светлая тема в шапке, состояние в localStorage и Zustand

### Ключевые компоненты
- `AppShell` (Header, Sidebar, ThemeToggle)
- `PostCard`, `PostComposer`, `CommentList`, `LikeButton`, `FollowButton`
- `FeedList` (TanStack Query + infinite scroll)
- `SearchBar` (+поиск по Enter)
- `Skeleton.*` (PostSkeleton, ProfileSkeleton)
- `ThemeProvider` (чтение/запись localStorage, класс dark на html)

### Состояние и запросы
- **TanStack Query:** кэш по ключам `['feed', params]`, `['post', id]`, `['profile', username]`, `['search', q]`
- **Мутации:** createPost, likePost, addComment, followUser
- **Авторизация:** accessToken в памяти (react-query fetcher вставляет заголовок); refresh-cookie прозрачен

---

# Примерный план по файлам

```bash
bailanysta/
  api/
    cmd/api/main.go
    internal/
      config/
      http/
        middleware/
        handlers/
      domain/
      services/
      repositories/
      db/
        queries/        # sqlc .sql файлы
        migrations/     # golang-migrate
      pkg/
        auth/
        ai/
        logger/
        validation/
      openapi/
        openapi.yaml
  web/
    src/
      app/
      components/
      pages/
      routes/
      hooks/
      store/
      lib/
    public/
    index.html
    vite.config.ts
    tailwind.config.js
  infra/
    docker/
      nginx.conf
    docker-compose.yml
    .env.example
  Makefile
  README.md
```