# 📚 Link Keeper

## Описание
Link Keeper — это Telegram-бот, который сохраняет ссылки, отправленные пользователями, и по их запросу выдаёт случайную из сохранённых. Это идеальное решение для тех, кто часто откладывает статьи "на потом", но затем забывает их прочитать. Бот помогает организовать временный список материалов: каждая выданная ссылка автоматически удаляется после отправки.

Код проекта спроектирован с акцентом на расширяемость. Например, поддержку другого мессенджера можно добавить, просто реализовав новый клиент — основная логика останется неизменной.

## Основные функции
- **Сохранение ссылок**: Отправка URL привязывает его к пользователю.
- **Выдача случайной ссылки**: Команда `/rnd` отправляет случайную ссылку и удаляет её из списка.
- **Команды**: `/start` (приветствие), `/help` (справка).
- **Уникальность**: Дубликаты ссылок для одного пользователя не сохраняются.

## Технологии
- **Библиотеки**:
  - `github.com/mattn/go-sqlite3` — для SQLite.
- **Хранилище**:
  - **Файловая система**: Сериализация данных в файлы с помощью `gob`.
  - **SQLite**: Альтернативная реализация с SQL-запросами.
- **Архитектура**:
  - Модульная структура: пакеты `clients`, `events`, `storage`, `consumer`.
  - Интерфейсы (`Fetcher`, `Processor`, `Storage`) для абстрагирования от реализаций.
- **Обработка ошибок**: Кастомная утилита `errwrap` для обёртки ошибок с контекстом.

## Структура проекта
- **`clients/telegram`**: HTTP-клиент для работы с Telegram API.
- **`events/telegram`**: Логика обработки событий и команд.
- **`storage`**: Интерфейс хранилища с реализациями для файлов (`files`) и SQLite (`sqlite`).
- **`consumer`**: Цикл обработки событий с настраиваемым батчем.
- **`lib/errwrap`**: Утилита для управления ошибками.

## Особенности
- **Абстрагированность**:
  - Интерфейсы `Fetcher` и `Processor` отделяют логику событий от Telegram, упрощая интеграцию с другими платформами.
  - Интерфейс `Storage` изолирует хранилище, поддерживая разные варианты (файлы, SQLite) и новые реализации.
- **Чистая структура**: Модульность и интерфейсы обеспечивают читаемость и удобство доработки.
