# Руководство по навигации Criage

Система навигации была полностью переработана для лучшего UX и современного дизайна.

## 🎯 Что изменилось

### ✅ Исправлено

- **Убраны дублирующиеся ссылки** в шапке главной страницы
- **Локальные якорные ссылки** больше не показываются в навигации
- **Четкое разделение** между страницами и секциями

### 🎨 Улучшен дизайн

- **Кнопочный дизайн** навигационных ссылок
- **Подсветка активной страницы** с подчеркиванием
- **Специальное выделение** для GitHub ссылки
- **Hover эффекты** и анимации
- **Мобильное меню** с гамбургер-кнопкой

## 📁 Структура файлов

```
docs/
├── navigation.js           # Основная логика навигации
├── navigation-styles.css   # Стили для всех страниц
├── mobile-menu.js         # Мобильное меню (универсальное)
├── os-switcher.js         # Переключатель ОС (опционально)
└── os-switcher.css        # Стили переключателя ОС
```

## 🚀 Быстрая интеграция

### Для новых страниц

Добавьте в `<head>`:

```html
<link rel="stylesheet" href="navigation-styles.css">
```

Добавьте перед закрывающим `</body>`:

```html
<script src="navigation.js"></script>
<script src="mobile-menu.js"></script>
```

### Обновление существующих страниц

#### 1. HTML структура

Замените навигацию на:

```html
<header>
    <nav class="container">
        <a href="index.html" class="logo">
            <img src="logo.png" alt="Criage Logo">
            Criage
        </a>
        <div class="nav-content">
            <ul class="nav-links">
                <!-- Заполняется navigation.js -->
            </ul>
            <div class="language-switcher">
                <!-- Заполняется navigation.js -->
            </div>
            <button class="mobile-menu-toggle" onclick="toggleMobileMenu()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                    <line x1="3" y1="6" x2="21" y2="6"></line>
                    <line x1="3" y1="12" x2="21" y2="12"></line>
                    <line x1="3" y1="18" x2="21" y2="18"></line>
                </svg>
            </button>
        </div>
        <div class="mobile-nav" id="mobileNav">
            <ul class="nav-links">
                <!-- Заполняется navigation.js -->
            </ul>
        </div>
    </nav>
</header>
```

#### 2. Подключение стилей

Добавьте в CSS или подключите файл:

```html
<link rel="stylesheet" href="navigation-styles.css">
```

#### 3. Подключение скриптов

```html
<script src="navigation.js"></script>
<script src="mobile-menu.js"></script>
```

## 🎨 Кастомизация

### CSS переменные

```css
:root {
    --primary-color: #2563eb;
    --primary-dark: #1d4ed8;
    --text-secondary: #64748b;
    --text-primary: #1e293b;
    --border: #e2e8f0;
    --surface: #f8fafc;
    --background: #ffffff;
    --shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}
```

### Переопределение стилей

```css
/* Кастомный цвет для вашей страницы */
.nav-link.active {
    color: #your-color !important;
}

.nav-link.active::after {
    background: #your-color !important;
}
```

## 📱 Мобильное меню

### Автоматическая инициализация

Мобильное меню инициализируется автоматически при загрузке страницы.

### Программное управление

```javascript
// Открыть меню
window.mobileMenu.open();

// Закрыть меню
window.mobileMenu.close();

// Переключить состояние
window.mobileMenu.toggle();
```

### События

```javascript
// Слушать обновления навигации
window.addEventListener('navigationUpdated', () => {
    console.log('Navigation updated');
});
```

## 🌐 Мультиязычность

Навигация автоматически определяет язык по URL:

- `page.html` → английский
- `page_ru.html` → русский

### Добавление нового языка

В `navigation.js`:

```javascript
NavigationConfig.languages.de = {
    name: "DE",
    links: [
        { id: "home", href: "index_de.html", text: "Startseite" },
        // ... другие ссылки
    ]
};
```

## 🔧 Конфигурация

### Добавление новой страницы

В `navigation.js`:

```javascript
NavigationConfig.languages.en.links.push({
    id: "new-page",
    href: "new-page.html", 
    text: "New Page"
});

NavigationConfig.languages.ru.links.push({
    id: "new-page",
    href: "new-page_ru.html",
    text: "Новая страница"
});

NavigationConfig.pageMapping["new-page.html"] = "new-page";
NavigationConfig.pageMapping["new-page_ru.html"] = "new-page";
```

### Кастомные GitHub ссылки

```javascript
NavigationConfig.githubLinks["new-page"] = "https://github.com/criage-oss/new-repo";
```

## 🎯 Доступность

Навигация поддерживает:

- ✅ **Клавиатурную навигацию** (Tab, Enter, Escape)
- ✅ **Screen readers** (aria-labels, roles)
- ✅ **Focus management** в мобильном меню
- ✅ **Высокий контраст** для текста и кнопок

## 🚨 Troubleshooting

### Навигация не загружается

1. Проверьте подключение `navigation.js`
2. Откройте консоль браузера
3. Убедитесь что элементы `.nav-links` существуют

### Мобильное меню не работает

1. Проверьте подключение `mobile-menu.js`
2. Убедитесь что есть кнопка `.mobile-menu-toggle`
3. Проверьте наличие `.mobile-nav` контейнера

### Стили не применяются

1. Проверьте подключение `navigation-styles.css`
2. Убедитесь в правильности CSS переменных
3. Проверьте специфичность селекторов

### Конфликты стилей

Используйте `!important` для критичных стилей:

```css
.nav-link {
    color: var(--primary-color) !important;
}
```

## 📊 Производительность

- **CSS файл**: ~8KB (сжатый)
- **JS файлы**: ~15KB (сжатые)
- **Загрузка**: Асинхронная, не блокирует рендеринг
- **Совместимость**: IE11+, все современные браузеры

## 🔄 Обновления

При обновлении навигации:

1. Обновите `navigation.js` для новых ссылок
2. Добавьте переводы для новых языков
3. Обновите `pageMapping` для новых страниц
4. Протестируйте мобильное меню

## 💡 Примеры использования

### Страница документации

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <link rel="stylesheet" href="navigation-styles.css">
</head>
<body>
    <header><!-- навигация --></header>
    <main><!-- контент --></main>
    <script src="navigation.js"></script>
    <script src="mobile-menu.js"></script>
</body>
</html>
```

### Страница с переключателем ОС

```html
<head>
    <link rel="stylesheet" href="navigation-styles.css">
    <link rel="stylesheet" href="os-switcher.css">
</head>
<body>
    <div class="auto-os-switch" data-template="install" data-package="criage-client"></div>
    <script src="navigation.js"></script>
    <script src="mobile-menu.js"></script>
    <script src="os-switcher.js"></script>
</body>
```

---

**Поддержка**: [GitHub Issues](https://github.com/criage-oss/criage.ru/issues)  
**Документация**: [criage.ru](https://criage.ru/)
