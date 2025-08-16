# OS Switcher для документации Criage

Универсальная система переключения между инструкциями установки для разных операционных систем.

## Возможности

- 🖱️ **Автоматическое определение ОС пользователя**
- 🌐 **Поддержка мультиязычности** (EN/RU)
- 💾 **Сохранение выбора пользователя** в localStorage
- 📱 **Адаптивный дизайн** для мобильных устройств
- 🎨 **Настраиваемые темы** (светлая/темная)
- 📋 **Копирование кода одним кликом**
- ⚡ **Автоматическое преобразование** простых блоков

## Быстрый старт

### 1. Подключение файлов

```html
<link rel="stylesheet" href="os-switcher.css">
<script src="os-switcher.js"></script>
```

### 2. Простое использование

Добавьте класс `auto-os-switch` к блоку с инструкциями:

```html
<div class="auto-os-switch" data-template="install" data-package="criage-client">
    <!-- Автоматически сгенерируется переключатель ОС -->
</div>
```

### 3. Ручное создание

```html
<div class="os-switcher">
    <div class="os-tabs">
        <button class="os-tab active" data-os="linux">🐧 Linux</button>
        <button class="os-tab" data-os="macos">🍎 macOS</button>
        <button class="os-tab" data-os="windows">🪟 Windows</button>
    </div>
    
    <div class="code-snippet os-content active" data-os="linux">
        <button class="copy-code" onclick="copyToClipboard(this)">Copy</button>
        <pre># Linux команды
sudo apt install criage</pre>
    </div>
    
    <div class="code-snippet os-content" data-os="macos">
        <button class="copy-code" onclick="copyToClipboard(this)">Copy</button>
        <pre># macOS команды
brew install criage</pre>
    </div>
    
    <div class="code-snippet os-content" data-os="windows">
        <button class="copy-code" onclick="copyToClipboard(this)">Copy</button>
        <pre># Windows команды
choco install criage</pre>
    </div>
</div>
```

## Доступные шаблоны

### install
Генерирует команды установки из GitHub releases:

```html
<div class="auto-os-switch" data-template="install" data-package="criage-client"></div>
<div class="auto-os-switch" data-template="install" data-package="criage-server"></div>
<div class="auto-os-switch" data-template="install" data-package="criage-mcp"></div>
```

### docker
Генерирует Docker команды:

```html
<div class="auto-os-switch" data-template="docker" data-package="criage-client"></div>
```

## API

### Создание программного переключателя

```javascript
const switcher = new OSSwitcher();

// Создать HTML переключателя
const html = switcher.createOSSwitch({
    showLinux: true,
    showMacOS: true,
    showWindows: true,
    defaultOS: 'linux'
});

// Создать блок кода
const codeBlock = switcher.createCodeBlock('linux', '# Linux команды\nls -la', true);

// Инициализировать все переключатели на странице
switcher.init();
```

### Определение ОС и языка

```javascript
const switcher = new OSSwitcher();
console.log(switcher.currentOS); // 'linux', 'macos', 'windows'
console.log(switcher.currentLang); // 'en', 'ru'
```

## Настройка стилей

Компонент использует CSS переменные для легкой настройки:

```css
:root {
    --primary-color: #2563eb;
    --primary-dark: #1d4ed8;
    --background: #ffffff;
    --surface: #f8fafc;
    --text-secondary: #64748b;
    --border: #e2e8f0;
}
```

## Расширенные возможности

### Кастомные генераторы команд

```javascript
class CustomOSSwitcher extends OSSwitcher {
    generateCustomInstall(packageName) {
        return `# Custom installation for ${packageName}
custom-install ${packageName}`;
    }
}
```

### Добавление новых ОС

```javascript
// В переопределенном классе
createOSSwitch(config = {}) {
    // Добавить поддержку FreeBSD, Android и т.д.
    return super.createOSSwitch({
        ...config,
        showFreeBSD: true
    });
}
```

## Мультиязычность

Компонент автоматически определяет язык по URL:
- `index.html` → английский
- `index_ru.html` → русский

Переводы находятся в объекте `translations`:

```javascript
translations: {
    en: {
        linux: '🐧 Linux',
        copy: 'Copy',
        copied: 'Copied!'
    },
    ru: {
        linux: '🐧 Linux', 
        copy: 'Копировать',
        copied: 'Скопировано!'
    }
}
```

## Интеграция с существующими страницами

1. **Добавьте CSS и JS файлы** в `<head>` секцию
2. **Замените простые блоки кода** на переключатели
3. **Используйте автоматические шаблоны** для стандартных случаев
4. **Настройте стили** под дизайн вашей страницы

## Примеры интеграции

### Документация API
```html
<div class="auto-os-switch" data-template="install" data-package="criage-client">
    <h3>Установка Criage Client</h3>
</div>
```

### Страница сервера
```html
<div class="auto-os-switch" data-template="docker" data-package="criage-server">
    <h3>Запуск сервера в Docker</h3>
</div>
```

### MCP Server документация  
```html
<div class="auto-os-switch" data-template="install" data-package="criage-mcp">
    <h3>Установка MCP Server</h3>
</div>
```

## Отладка

Включите логирование в консоли:

```javascript
window.osSwitcher.debug = true;
```

Проверьте корректность данных:

```javascript
// Проверить определение ОС
console.log('Detected OS:', window.osSwitcher.currentOS);

// Проверить сохраненные настройки
console.log('Saved OS:', localStorage.getItem('criage-preferred-os'));
```

## Требования

- Современный браузер с поддержкой ES6+
- CSS Grid и Flexbox
- LocalStorage для сохранения настроек
- Clipboard API для копирования (опционально)

## Лицензия

MIT License - используйте свободно в проектах Criage.
