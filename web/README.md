# Bailanysta Frontend 🌟

Современный React frontend для образовательной платформы Bailanysta.

## 🛠 Технологии

- **React 18** + **TypeScript** - Современный UI фреймворк с строгой типизацией
- **Vite** - Быстрая сборка и hot reload
- **TanStack Query** - Мощное управление серверным состоянием
- **Zustand** - Легковесное управление клиентским состоянием
- **Tailwind CSS** - Utility-first CSS фреймворк
- **Framer Motion** - Анимации и переходы
- **Radix UI** - Доступные UI примитивы
- **Lucide React** - Красивые иконки
- **React Router DOM** - Клиентская маршрутизация

## 🏗 Архитектура

```
src/
├── api/           # API клиент и типы запросов
├── components/    # React компоненты
│   ├── ui/        # Базовые UI компоненты (Button, Input, Card...)
│   ├── layout/    # Компоненты макета (Header, Sidebar, Layout)
│   └── features/  # Функциональные компоненты (PostCard, UserCard...)
├── hooks/         # Кастомные React хуки
├── pages/         # Страницы приложения
├── stores/        # Zustand stores (auth, theme, ui)
├── types/         # TypeScript типы
├── lib/           # Утилиты и хелперы
└── assets/        # Статические ресурсы
```

## 🚀 Запуск

### Разработка

```bash
npm install
npm run dev
```

Приложение будет доступно на http://localhost:3000

### Сборка

```bash
npm run build
npm run preview
```

## 🎨 UI Система

Приложение использует современную систему дизайна с:

- **Темная/Светлая тема** - Автоматическое переключение и сохранение
- **Responsive дизайн** - Адаптивная верстка для всех устройств
- **Анимации** - Плавные переходы с Framer Motion
- **Accessibility** - Доступность через Radix UI компоненты

## 📱 Основные функции

### Аутентификация
- Регистрация и вход с email/паролем
- JWT токены с автоматическим обновлением
- Защищенные маршруты

### Посты и лента
- Создание постов с поддержкой хештегов
- Бесконечная прокрутка ленты
- Лайки и комментарии
- Поиск по контенту

### Социальные функции
- Подписки на пользователей
- Уведомления в реальном времени
- Профили пользователей
- Система комментариев

### AI интеграция
- Генерация постов с помощью AI
- Создание комментариев
- Объяснение концепций
- Генерация конспектов и викторин

## 🔧 Конфигурация

### Переменные окружения

```bash
# API Configuration
VITE_API_URL=http://localhost:8080

# App Configuration
VITE_APP_NAME=Bailanysta
VITE_APP_DESCRIPTION="Образовательная платформа с AI-интеграцией"

# Feature Flags
VITE_ENABLE_AI_FEATURES=true
VITE_ENABLE_NOTIFICATIONS=true
VITE_ENABLE_REAL_TIME=false
```

## 🧪 Тестирование

```bash
# Запуск тестов
npm run test

# Запуск тестов в watch режиме
npm run test:watch

# Линтинг
npm run lint
```

## 📦 State Management

### Auth Store (Zustand)
```typescript
const { user, login, logout, isAuthenticated } = useAuthStore()
```

### Theme Store
```typescript
const { theme, setTheme, toggleTheme } = useThemeStore()
```

### UI Store
```typescript
const { sidebarOpen, toggleSidebar, setComposerOpen } = useUIStore()
```

## 🌐 API Integration

Приложение использует TanStack Query для управления серверным состоянием:

```typescript
// Получение ленты с бесконечной прокруткой
const { data, fetchNextPage, hasNextPage } = useInfiniteQuery({
  queryKey: ['feed'],
  queryFn: ({ pageParam }) => apiClient.getFeed(20, pageParam),
  getNextPageParam: (lastPage) => lastPage.nextOffset
})

// Создание поста с оптимистичным обновлением
const createPostMutation = useMutation({
  mutationFn: apiClient.createPost,
  onSuccess: () => {
    queryClient.invalidateQueries(['feed'])
  }
})
```

## 🎯 Компоненты

### UI Компоненты
- `Button` - Кнопки с вариантами стилей
- `Input/Textarea` - Поля ввода с валидацией
- `Card` - Карточки для контента
- `Avatar` - Аватары пользователей
- `Toast` - Уведомления

### Feature Компоненты
- `PostCard` - Карточка поста с действиями
- `PostComposer` - Форма создания поста
- `UserCard` - Карточка пользователя
- `CommentCard` - Карточка комментария

### Layout Компоненты
- `Header` - Шапка с навигацией и поиском
- `Sidebar` - Боковая панель навигации
- `Layout` - Основной макет приложения

## 🔒 Безопасность

- JWT токены с автоматическим обновлением
- Защищенные маршруты с `ProtectedRoute`
- Валидация входных данных
- XSS защита через sanitization

## 📱 Responsive Design

Приложение полностью адаптивно:
- **Mobile** (< 768px) - Мобильное меню, компактная навигация
- **Tablet** (768px - 1024px) - Адаптивная сетка
- **Desktop** (> 1024px) - Полный интерфейс с боковой панелью

## 🚀 Производительность

- **Code Splitting** - Автоматическое разделение кода
- **Lazy Loading** - Ленивая загрузка компонентов
- **Image Optimization** - Оптимизация изображений
- **Bundle Analysis** - Анализ размера бандла

## 🤝 Разработка

### Создание нового компонента

```typescript
// components/features/NewComponent.tsx
import { motion } from 'framer-motion'
import { Button } from '@/components/ui/button'

interface NewComponentProps {
  title: string
  onAction: () => void
}

export default function NewComponent({ title, onAction }: NewComponentProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="p-4 bg-card rounded-lg"
    >
      <h2 className="text-xl font-bold mb-4">{title}</h2>
      <Button onClick={onAction}>Действие</Button>
    </motion.div>
  )
}
```

### Создание новой страницы

```typescript
// pages/NewPage.tsx
export default function NewPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Новая страница</h1>
      {/* Контент страницы */}
    </div>
  )
}
```

## 📋 TODO

- [ ] Добавить E2E тесты с Playwright
- [ ] Реализовать PWA функциональность
- [ ] Добавить поддержку файлов/изображений
- [ ] Интеграция с WebSocket для real-time обновлений
- [ ] Оптимизация производительности

---

Создано с ❤️ для современного образования
