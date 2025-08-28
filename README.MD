# Bailanysta 🎓

Современная образовательная платформа с AI-интеграцией для обучения и общения.

## 🚀 Ключевые возможности

- **📚 Образовательный контент**: Курсы, модули и учебные материалы
- **🤖 AI-ассистент**: Генерация контента, викторины, объяснения концепций
- **👥 Социальная сеть**: Посты, комментарии, подписки, уведомления
- **🔍 Поиск**: Полнотекстовый поиск с хештегами
- **📱 Современный UI**: React + TypeScript + Tailwind CSS
- **🐳 Docker**: Легкое развертывание в любом окружении

## 🏗️ Архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React Web     │    │   Go Backend    │    │   PostgreSQL    │
│   Frontend      │◄──►│   API Server    │◄──►│   Database      │
│                 │    │                 │    │                 │
│ • TypeScript    │    │ • REST API      │    │ • Users         │
│ • Tailwind CSS  │    │ • JWT Auth      │    │ • Posts         │
│ • Zustand       │    │ • AI Integration│    │ • Comments      │
│ • React Router  │    │ • Rate Limiting │    │ • Notifications │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              ▲
                              │
                       ┌─────────────────┐
                       │     Redis       │
                       │   (Optional)    │
                       │                 │
                       │ • Caching       │
                       │ • Sessions      │
                       └─────────────────┘
```

## 🛠️ Технологии

### Backend
- **Go 1.23** - Высокопроизводительный бэкенд
- **PostgreSQL 15** - Надежная база данных
- **Redis 7** - Кеширование и сессии
- **JWT** - Аутентификация
- **Chi Router** - HTTP роутинг
- **sqlc** - Type-safe SQL генерация

### Frontend
- **React 18** - Современный UI фреймворк
- **TypeScript** - Строгая типизация
- **Tailwind CSS** - Utility-first CSS
- **Vite** - Быстрая сборка
- **Zustand** - State management
- **TanStack Query** - API state management

### DevOps
- **Docker & Docker Compose** - Контейнеризация
- **Nginx** - Reverse proxy и static files
- **Air** - Hot reload для Go
- **Vitest** - Тестирование

## 🚀 Быстрый старт

### Предварительные требования
- Docker 20.10+
- Docker Compose 2.0+
- Git

### Развертывание

1. **Клонируем репозиторий:**
   ```bash
   git clone <repository-url>
   cd bailanysta
   ```

2. **Запускаем development среду:**
   ```bash
   make dev
   # или
   ./infra/docker/deploy.sh dev
   ```

3. **Открываем приложение:**
   - 🌐 Frontend: http://localhost:3000
   - 🔌 API: http://localhost:8080
   - 🗄️ Database: localhost:5432

### Production развертывание

1. **Настраиваем переменные окружения:**
   ```bash
   cp infra/docker/env.example .env
   # Редактируем .env с production настройками
   ```

2. **Запускаем production:**
   ```bash
   make prod
   # или
   ./infra/docker/deploy.sh prod
   ```

## 📁 Структура проекта

```
bailanysta/
├── api/                          # Go backend
│   ├── cmd/api/                  # Main application
│   ├── internal/
│   │   ├── config/               # Configuration
│   │   ├── db/                   # Database layer
│   │   │   ├── migrations/       # SQL migrations
│   │   │   └── queries/          # SQL queries for sqlc
│   │   ├── http/                 # HTTP handlers
│   │   ├── pkg/                  # Shared packages
│   │   ├── services/             # Business logic
│   │   └── logger/               # Logging
│   └── pkg/                      # Public packages
├── web/                          # React frontend
│   ├── src/
│   │   ├── components/           # React components
│   │   ├── pages/                # Page components
│   │   ├── hooks/                # Custom hooks
│   │   ├── stores/               # Zustand stores
│   │   ├── lib/                  # Utilities
│   │   ├── types/                # TypeScript types
│   │   └── api/                  # API client
│   └── nginx.conf               # Nginx config for web
├── infra/                        # Infrastructure
│   ├── docker/                   # Docker configs
│   ├── nginx/                    # Nginx configs
│   └── ssl/                      # SSL certificates
├── docker-compose.yml           # Main compose file
├── docker-compose.prod.yml      # Production overrides
├── docker-compose.override.yml  # Development overrides
├── Dockerfile                   # Go production build
├── Dockerfile.dev              # Go development build
└── Makefile                     # Development helpers
```

## 🧪 Тестирование

### Запуск всех тестов
```bash
make test
```

### Тестирование только бэкенда
```bash
make test-go
```

### Тестирование только фронтенда
```bash
make test-frontend
```

### Тестирование с watch режимом
```bash
make test-frontend-watch
```

## 📚 API Документация

### Аутентификация
- `POST /api/v1/auth/register` - Регистрация
- `POST /api/v1/auth/login` - Вход
- `POST /api/v1/auth/refresh` - Обновление токена
- `GET /api/v1/auth/me` - Текущий пользователь

### Посты
- `GET /api/v1/feed` - Лента постов
- `POST /api/v1/posts` - Создать пост
- `GET /api/v1/posts/{id}` - Получить пост
- `PUT /api/v1/posts/{id}` - Обновить пост
- `DELETE /api/v1/posts/{id}` - Удалить пост
- `POST /api/v1/posts/{id}/like` - Лайкнуть пост

### Социальные функции
- `POST /api/v1/follow/{user_id}` - Подписаться
- `DELETE /api/v1/follow/{user_id}` - Отписаться
- `GET /api/v1/users/{username}/followers` - Подписчики
- `GET /api/v1/users/{username}/following` - Подписки

### Поиск
- `GET /api/v1/search?q={query}` - Поиск постов и пользователей

### Уведомления
- `GET /api/v1/notifications` - Список уведомлений
- `POST /api/v1/notifications/mark-read` - Отметить все как прочитанные
- `GET /api/v1/notifications/unread-count` - Количество непрочитанных

### AI Функции
- `POST /api/v1/ai/generate` - Генерация текста
- `POST /api/v1/ai/generate-post` - Генерация поста
- `POST /api/v1/ai/generate-comment` - Генерация комментария
- `POST /api/v1/ai/generate-study-notes` - Генерация конспекта
- `POST /api/v1/ai/generate-quiz` - Генерация викторины
- `POST /api/v1/ai/explain-concept` - Объяснение концепции

## 🔧 Разработка

### Установка зависимостей
```bash
make install
```

### Запуск в режиме разработки
```bash
make dev
```

### Линтинг и форматирование
```bash
make lint
make format
```

### Работа с базой данных
```bash
make migrate    # Запуск миграций
make backup     # Создание бэкапа
```

## 🚀 Деплой

### Development
```bash
make dev
```

### Production
```bash
make deploy-check  # Предварительные проверки
make prod         # Запуск production среды
```

### Управление сервисами
```bash
make logs         # Все логи
make logs-api     # Логи API
make stop         # Остановить все
make cleanup      # Очистка ресурсов
```

## 🤝 Вклад в проект

1. Fork проект
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📝 Лицензия

Этот проект распространяется под лицензией MIT. См. файл `LICENSE` для подробностей.

## 📞 Контакты

- **Email**: your-email@example.com
- **GitHub**: https://github.com/yourusername/bailanysta
- **Документация**: [Ссылка на документацию]

---

## Компромиссы
- SSE/Realtime отложены (MVP — polling).
- Полнотекстовый поиск упрощён (trgm, без TS‑ранжирования).

## Известные проблемы
- Нет email‑верификации (можно добавить позже).
- Отсутствует масштабирование (горизонтально — в будущем).


**Bailanysta** - ваша платформа для современного образования! 🎓✨