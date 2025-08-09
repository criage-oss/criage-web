# Процесс релиза Criage

Этот документ описывает процедуру создания релизов для всех компонентов экосистемы Criage.

## 🏗️ Архитектура CI/CD

Проект Criage состоит из трех основных компонентов:

1. **criage-client** - основной пакетный менеджер
2. **criage-server** - HTTP сервер репозитория пакетов  
3. **criage-mcp** - MCP сервер для интеграции с AI

Каждый компонент имеет собственный GitHub Actions workflow для сборки и публикации.

## 📋 Предварительные требования

### GitHub Secrets

Для работы CI/CD необходимы следующие секреты:

#### В каждом репозитории компонентов

- `DOCKER_USERNAME` - имя пользователя Docker Hub
- `DOCKER_PASSWORD` - пароль Docker Hub

#### В главном репозитории (criage.ru)

- `PAT_TOKEN` - Personal Access Token с правами на создание тегов во всех репозиториях

### Настройка PAT Token

1. Перейдите в GitHub Settings → Developer settings → Personal access tokens
2. Создайте новый token с правами:
   - `repo` (полный доступ к репозиториям)
   - `workflow` (управление GitHub Actions)
3. Добавьте токен как секрет `PAT_TOKEN` в репозиторий criage.ru

## 🚀 Процесс релиза

### Автоматический релиз (рекомендуется)

1. Перейдите в репозиторий [criage.ru](https://github.com/criage-oss/criage.ru)
2. Откройте Actions → "Release All Criage Components"
3. Нажмите "Run workflow"
4. Укажите версию (например, `1.0.0`)
5. Выберите, является ли релиз pre-release
6. Нажмите "Run workflow"

Система автоматически:

- Создаст теги во всех репозиториях
- Запустит сборки для всех платформ
- Создаст релизы с артефактами
- Опубликует Docker образы
- Обновит информацию на сайте

### Ручной релиз

Если нужно выпустить только один компонент:

#### 1. Создание тега

```bash
git tag v1.0.0
git push origin v1.0.0
```

#### 2. Автоматическая сборка

GitHub Actions автоматически начнет сборку при создании тега.

## 📦 Артефакты сборки

### Каждый компонент создает

**Бинарные файлы:**

- Windows (amd64, arm64) - `.exe` в `.zip` архивах
- Linux (amd64, arm64) - бинарии в `.tar.gz` архивах  
- macOS (amd64, arm64) - бинарии в `.tar.gz` архивах

**Docker образы:**

- `criage/client:версия`
- `criage/server:версия`
- `criage/mcp-server:версия`

**Конфигурационные файлы:**

- Примеры конфигурации
- Документация
- README файлы

## 🏷️ Схема версионирования

Проект использует [Semantic Versioning](https://semver.org/):

- `MAJOR.MINOR.PATCH` (например, `1.0.0`)
- `MAJOR` - breaking changes
- `MINOR` - новая функциональность (обратно совместимая)
- `PATCH` - исправления ошибок

### Примеры

- `1.0.0` - первый стабильный релиз
- `1.1.0` - добавлена новая функция
- `1.1.1` - исправлена ошибка
- `2.0.0` - breaking changes

## 🔄 Workflow этапы

### 1. Test Stage

- Запуск unit тестов
- Проверка форматирования кода
- Статический анализ (go vet)

### 2. Build Stage  

- Кроссплатформенная сборка
- Создание архивов
- Загрузка артефактов

### 3. Release Stage (только для тегов)

- Создание GitHub релиза
- Прикрепление артефактов
- Генерация changelog

### 4. Docker Stage

- Мультиплатформенная сборка образов
- Публикация в Docker Hub
- Тегирование latest для main ветки

## 🐛 Устранение неполадок

### Сборка не запускается

- Проверьте права PAT токена
- Убедитесь, что тег создан корректно
- Проверьте статус GitHub Actions

### Docker публикация не работает

- Проверьте DOCKER_USERNAME и DOCKER_PASSWORD
- Убедитесь в правильности имен образов

### Зависимости между компонентами

- criage-server зависит от criage-client через replace directive
- При изменении criage-client может потребоваться обновление server

## 📊 Мониторинг релизов

### Проверка статуса сборок

- [Client Actions](https://github.com/criage-oss/criage-client/actions)
- [Server Actions](https://github.com/criage-oss/criage-server/actions)  
- [MCP Actions](https://github.com/criage-oss/criage-mcp/actions)

### Docker Hub

- [criage/client](https://hub.docker.com/r/criage/client)
- [criage/server](https://hub.docker.com/r/criage/server)
- [criage/mcp-server](https://hub.docker.com/r/criage/mcp-server)

## 📝 Чеклист релиза

- [ ] Все тесты проходят
- [ ] Документация обновлена
- [ ] Версии синхронизированы во всех компонентах
- [ ] Changelog подготовлен
- [ ] PAT токен актуален
- [ ] Docker credentials настроены
- [ ] Релиз протестирован на основных платформах

## 🔧 Настройка среды разработки

```bash
# Клонирование всех репозиториев
git clone https://github.com/criage-oss/criage-client.git
git clone https://github.com/criage-oss/criage-server.git  
git clone https://github.com/criage-oss/criage-mcp.git
git clone https://github.com/criage-oss/criage.ru.git

# Настройка зависимостей
cd criage-client && go mod tidy
cd ../criage-server && go mod tidy
cd ../criage-mcp && go mod tidy
```

## 🆘 Контакты

При проблемах с релизами обращайтесь:

- GitHub Issues в соответствующем репозитории
- Команда разработки Criage

---

**Последнее обновление:** 09.01.2025  
**Версия документа:** 1.0.0
