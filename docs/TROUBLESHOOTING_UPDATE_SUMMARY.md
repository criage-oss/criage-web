# 📋 Troubleshooting Updates Summary

## 🎯 **Задача выполнена!**

Проведен полный анализ кода всех компонентов Criage и обновлены разделы troubleshooting в документации.

## 🔍 **Анализ компонентов**

### **1. Criage Client (`criage-client`)**
**Проанализированные файлы:**
- `main.go` - инициализация и обработка команд
- `commands.go` - логика выполнения команд
- `pkg/config.go` - управление конфигурацией
- `pkg/package_manager.go` - основная логика менеджера пакетов

**Выявленные потенциальные ошибки:**
- Ошибки инициализации конфигурации (`failed to create config manager`)
- Проблемы с домашней директорией (`failed to get home directory`)
- Ошибки создания архивного менеджера (`failed to create archive manager`)
- HTTP таймауты и проблемы сети (`context deadline exceeded`)
- Rate limiting ошибки (`rate limit exceeded`)
- Проблемы с форматами архивов (`unsupported archive format`)
- Ошибки файловой системы (`no space left on device`, `permission denied`)
- Проблемы с локализацией (`translation key not found`)
- JSON parsing ошибки
- Проблемы с зависимостями (`dependency conflict`)

### **2. Criage Server (`criage-server`)**
**Проанализированные файлы:**
- `main.go` - запуск сервера и конфигурация
- `api.go` - API эндпоинты и HTTP обработчики
- `index.go` - управление индексом пакетов

**Выявленные потенциальные ошибки:**
- Ошибки привязки порта (`bind: address already in use`)
- Проблемы с JSON конфигурацией (`invalid character`, `failed to unmarshal`)
- Ошибки создания директорий (`failed to create directories`)
- Проблемы загрузки файлов (`file too large`, `unsupported format`)
- Ошибки аутентификации (`authentication failed`)
- CORS проблемы (`blocked by CORS policy`)
- Проблемы с индексом (`failed to update index`, `index corrupted`)
- API ошибки (`404 Not Found`, `500 Internal Server Error`)
- Проблемы производительности (высокое потребление памяти)
- Дисковые проблемы (`no space left on device`)

### **3. Criage MCP (`criage-mcp`)**
**Проанализированные файлы:**
- `main.go` - MCP протокол и JSON-RPC
- `package_manager.go` - интеграция с пакетным менеджером

**Выявленные потенциальные ошибки:**
- Ошибки инициализации пакетного менеджера (`Не удалось создать пакетный менеджер`)
- JSON-RPC коммуникации (`Ошибка декодирования сообщения`)
- Проблемы с протоколом MCP (`Неизвестный метод`)
- Rate limiter проблемы (`panic: send on closed channel`) - **ИСПРАВЛЕНО**
- Claude Desktop интеграция (tools не отображаются)
- Проблемы с конфигурацией MCP сервера
- Ошибки валидации схем инструментов

## 📝 **Обновленные разделы документации**

### **✅ Criage Client** 
- **Файлы:** `docs.html`, `docs_ru.html`
- **Добавлено:** 15+ новых типов ошибок с решениями
- **Новые разделы:**
  - Installation & Setup Issues
  - Configuration Issues 
  - Package Management Errors
  - Network & Repository Issues
  - Archive & Format Issues
  - File System Issues
  - Debug Mode & Logging
  - Recovery Procedures

### **✅ Criage Server**
- **Файлы:** `repository-server.html`, `repository-server_ru.html` 
- **Добавлено:** 12+ новых типов ошибок с решениями
- **Новые разделы:**
  - Server Startup Issues
  - Upload & File Management Issues
  - Network & API Issues  
  - Performance & Resource Issues
  - Debug Mode & Logging
  - Recovery Procedures

### **✅ Criage MCP Server**
- **Файлы:** `mcp-server.html`, `mcp-server_ru.html`
- **Добавлено:** 10+ новых типов ошибок с решениями
- **Новые разделы:**
  - MCP Protocol Issues
  - Claude Desktop Integration
  - Package Management Errors
  - File System & Permission Issues
  - Debug & Development

## 🛠️ **Типы решений, добавленных в документацию:**

### **🔧 Диагностические команды:**
```bash
# Проверка статуса
criage --version
go version
which criage

# Проверка конфигурации
criage config list
cat ~/.config/criage/config.yaml

# Проверка сети
curl -I https://repo.criage.ru/api/v1/
ping repo.criage.ru
```

### **🔍 Отладка:**
```bash
# Debug режим
export CRIAGE_DEBUG=1
criage install package-name --verbose

# Логи
tail -f /var/log/criage.log
journalctl -u criage-server -f
```

### **🚨 Восстановление:**
```bash
# Сброс конфигурации
rm ~/.config/criage/config.yaml
criage config list

# Очистка кеша
rm -rf ~/.cache/criage/*

# Исправление прав
sudo chown -R $USER:$USER ~/.config/criage
chmod -R 755 ~/.config/criage
```

### **🌐 Проверка сервисов:**
```bash
# Проверка портов
netstat -tlnp | grep :8080
lsof -i :8080

# API тестирование
curl -v http://localhost:8080/api/v1/
curl -X POST http://localhost:8080/api/v1/upload

# MCP протокол
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./criage-mcp
```

## 📊 **Статистика обновлений:**

| Компонент | Старые ошибки | Новые ошибки | Добавлено решений |
|-----------|---------------|--------------|-------------------|
| **Client** | 3 | 15+ | 50+ команд |
| **Server** | 3 | 12+ | 40+ команд |
| **MCP** | 3 | 10+ | 30+ команд |
| **Всего** | **9** | **37+** | **120+ команд** |

## 🎨 **Улучшения форматирования:**

### **До:**
```html
<h4>Network Errors</h4>
<ul>
  <li>Check repository availability</li>
  <li>Increase timeout</li>
</ul>
```

### **После:**
```html
<h4>HTTP Request Timeout</h4>
<div class="warning">
    <strong>Error:</strong> <code>connection timeout</code>
</div>
<p><strong>Solutions:</strong></p>
<div class="code-block">
    <button class="copy-btn">Copy</button>
    <pre># Increase timeout
criage config set timeout 300

# Test connectivity
ping repo.criage.ru</pre>
</div>
```

## 🎯 **Основные улучшения:**

1. **📋 Структурированные ошибки** - каждая ошибка имеет конкретное сообщение и решения
2. **💻 Готовые команды** - все решения содержат copy-paste команды
3. **🎨 Улучшенное форматирование** - использование warning блоков и code-block
4. **🔍 Диагностические процедуры** - пошаговые инструкции для отладки
5. **🚨 Процедуры восстановления** - полное восстановление при критических ошибках
6. **🌍 Мультиплатформенность** - решения для Linux, macOS, Windows
7. **🔧 Конфигурационные исправления** - готовые примеры config.json
8. **📊 Мониторинг** - команды для отслеживания производительности

## 📁 **Обновленные файлы:**

### **Client документация:**
- ✅ `criage.ru/docs/docs.html` 
- ✅ `criage.ru/docs/docs_ru.html`

### **Server документация:**  
- ✅ `criage.ru/docs/repository-server.html`
- ✅ `criage.ru/docs/repository-server_ru.html`

### **MCP документация:**
- ✅ `criage.ru/docs/mcp-server.html` 
- ✅ `criage.ru/docs/mcp-server_ru.html`

---

**🎉 Задача полностью выполнена! Все разделы troubleshooting значительно расширены и содержат практические решения для реальных проблем, выявленных в коде.**
