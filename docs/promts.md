7.1. api/cmd/api/main.go
Create file api/cmd/api/main.go with a Go HTTP server using chi. Read config from env (PORT, DATABASE_URL, JWT_SECRET, OPENAI_BASE_URL, OPENAI_API_KEY, CORS_ORIGIN). Initialize logger, DB (pgxpool), run migrations if MIGRATE_ON_START=true. Mount routes from internal/http/handlers. Health endpoint GET /health returns {"ok":true}. Use structured logging.

7.2. Конфиг и роутер
Create files:
- api/internal/config/config.go: struct Config { Port, DatabaseURL, JwtSecret, OpenAIBaseURL, OpenAIApiKey, CORSOrigin string; MigrateOnStart bool }. Load from env with defaults and validate.
- api/internal/http/router.go: chi.NewRouter with middlewares (requestID, realIP, recoverer, logger, CORS for CORSOrigin, rate-limiter token bucket 100 rpm).
Export func NewRouter(deps *Deps) *chi.Mux where Deps holds services.

7.3. Домены/репозитории/sqlc
Use sqlc. Create SQL query files under api/internal/db/queries:
- users.sql, posts.sql, comments.sql, likes.sql, follows.sql, hashtags.sql, notifications.sql
Implement typical CRUD and aggregation queries (get post with counts, feed by follows, search by trigram).
Generate sqlc.yaml and Go package db with typed methods.

7.4. Auth сервис и хендлеры
Create:
- api/internal/pkg/auth/jwt.go: CreateAccessToken(userID), ParseAccessToken, Refresh flow with HttpOnly cookies.
- api/internal/services/auth.go with Register, Login using bcrypt and repositories.
- api/internal/http/handlers/auth.go with endpoints /auth/register, /auth/login, /auth/refresh, /auth/logout. Validate inputs via go-playground/validator. Return JSON contracts as in spec.

7.5. Posts/Comments/Likes
Create:
- api/internal/services/posts.go: CreatePost (parse hashtags), GetPost (with counts), UpdatePost (author-only), DeletePost.
- api/internal/http/handlers/posts.go: routes /posts (POST, GET by id, PATCH, DELETE), /posts/{id}/like (POST/DELETE), /posts/{id}/comments (GET/POST).
Implement permission checks, errors -> proper HTTP codes.

7.6. Follow/Feed/Search
Create:
- api/internal/services/social.go: FollowUser, UnfollowUser, GetFeed(userID, limit, offset)
- api/internal/http/handlers/social.go: POST/DELETE /users/{id}/follow, GET /feed
- api/internal/http/handlers/search.go: GET /search?query=...
Search: if query starts with '#', search by hashtag join; else trigram on posts.text.

7.7. AI-клиент и эндпоинт
Create:
- api/internal/pkg/ai/client.go using github.com/sashabaranov/go-openai with custom BaseURL from config. Provide method Generate(ctx, prompt string, maxTokens int, temperature float32) (string, error) using ChatCompletions with model "gpt-oss-120b".
- api/internal/http/handlers/ai.go: POST /ai/generate reading {prompt, max_tokens?, temperature?} (validate), call aiClient.Generate, return {text}.
Ensure NO client-side API keys; only server calls external AI.

7.8. OpenAPI (минимум)
Create api/openapi/openapi.yaml describing endpoints from the spec: /auth/*, /posts/*, /feed, /search, /ai/generate, schemas for Post, User, Error, Pagination.

7.9. Frontend scaffold
Create Vite React TS app in /web with:
- React Router routes: /feed, /compose, /profile/:username, /post/:id, /search
- TanStack Query client with fetcher adding Authorization header from state
- Zustand store for theme + auth token
- Tailwind configured (darkMode: 'class'), ThemeToggle component setting 'dark' class on <html> and persisting in localStorage
- Components: PostCard, PostComposer (tabs: Write | AI), CommentList, LikeButton, FollowButton, Skeletons
- Pages implementing data fetching and mutations to backend
- Loading skeletons and error states
- Simple login/register pages

7.10. Докер и Nginx
Create infra/docker/docker-compose.yml with services:
- db: postgres:16, volumes, POSTGRES_PASSWORD from .env
- migrate: migrate/migrate: run migrations against db
- api: build ./api, depends_on db,migrate, expose 8080
- web: build ./web (npm ci && npm run build), then served by nginx:alpine with /usr/share/nginx/html
- nginx: reverse proxy /api -> api:8080, serve web static; ports: 80:80
Create infra/docker/nginx.conf with SPA fallback and proxy_pass for /api/.
Provide .env.example and Makefile targets: run, stop, logs.

7.11. Тесты
Backend:
- Add api/internal/http/handlers/auth_test.go using httptest for register/login
- Add posts_test.go for create/read/like/comment with in-memory or testcontainers Postgres (ory/dockertest)
- Add services tests with table-driven style

Frontend:
- Add web/src/pages/__tests__/Feed.test.tsx using Vitest + Testing Library + MSW to mock /feed
- Add PostComposer AI tab test with MSW mock for /api/v1/ai/generate
