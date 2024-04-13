## Обзор Проекта

Доброго времени суток! В связи с болезнью, начало работы над проектом было перенесено на среду. Ниже представлены ключевые детали проекта, включая инструкции по запуску, использованные технологии и структуру проекта.
## Запуск проекта

Для запуска проекта используйте команду:

```bash
make dc
```

Команда запустит необходимые сервисы через docker-compose, описанные в файле docker-compose.yml.
## Тестирование

Тестирование производительности сервиса проведено с использованием инструмента [k6](https://k6.io/). Настройки тестирования включали 95% запросов на чтение и 5% на запись. Сценарий тестирования расположен в файле:

```plaintext
test/k6/read_and_random_write.js
```
Ниже представлены результаты тестирования:

![k6](/tests/k6/read_and_random_write_result.png)
## Использованные инструменты
### Нейросети

В проекте использовался GitHub Copilot для генерации бойлерплейт кода.
### Технологический стек

- Язык программирования: Go
- База данных: PostgreSQL
- Основные библиотеки и фреймворки:
  - Роутер: Fiber
  - Кэширование: GoCache
  - Логгирование: Slog
  - Драйвер PostgreSQL: Pgx

### Пропущенные задания

В рамках текущего задания не были реализованы дополнительные задачи. Отсутствие попыток реализации дополнительных задач обусловлено ограниченным временем, вызванным задержкой начала работы над проектом.