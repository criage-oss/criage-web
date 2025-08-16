/**
 * Централизованная навигация для сайта Criage
 * Поддерживает мультиязычность и автоматическое определение активных ссылок
 */

const NavigationConfig = {
    // Конфигурация языков и навигации
    languages: {
        en: {
            name: "EN",
            links: [
                { id: "home", href: "index.html", text: "Home" },
                { id: "client-docs", href: "docs.html", text: "Client Docs" },
                { id: "repository-server", href: "repository-server.html", text: "Repository Server" },
                { id: "mcp-server", href: "mcp-server.html", text: "MCP Server" },
                { id: "cicd-examples", href: "cicd-examples.html", text: "CI/CD Examples" }
            ]
        },
        ru: {
            name: "RU", 
            links: [
                { id: "home", href: "index_ru.html", text: "Главная" },
                { id: "client-docs", href: "docs_ru.html", text: "Документация клиента" },
                { id: "repository-server", href: "repository-server_ru.html", text: "Сервер репозитория" },
                { id: "mcp-server", href: "mcp-server_ru.html", text: "MCP Server" },
                { id: "cicd-examples", href: "cicd-examples_ru.html", text: "Примеры CI/CD" }
            ]
        }
    },

    // GitHub ссылки для разных типов страниц
    githubLinks: {
        "index": "https://github.com/criage-oss/criage-client",
        "docs": "https://github.com/criage-oss/criage-client", 
        "repository-server": "https://github.com/criage-oss/criage-server",
        "mcp-server": "https://github.com/criage-oss/criage-mcp",
        "cicd-examples": "https://github.com/criage-oss/criage-client"
    },

    // Маппинг файлов на ID страниц для определения активных ссылок
    pageMapping: {
        "index.html": "home",
        "index_ru.html": "home",
        "docs.html": "client-docs",
        "docs_ru.html": "client-docs",
        "repository-server.html": "repository-server",
        "repository-server_ru.html": "repository-server",
        "mcp-server.html": "mcp-server",
        "mcp-server_ru.html": "mcp-server",
        "cicd-examples.html": "cicd-examples",
        "cicd-examples_ru.html": "cicd-examples"
    }
};

class NavigationManager {
    constructor() {
        this.currentPage = this.getCurrentPageName();
        this.currentLanguage = this.detectLanguage();
        this.currentPageId = this.getCurrentPageId();
    }

    /**
     * Получает имя текущей страницы из URL
     */
    getCurrentPageName() {
        const path = window.location.pathname;
        return path.split('/').pop() || 'index.html';
    }

    /**
     * Определяет язык страницы по имени файла
     */
    detectLanguage() {
        return this.currentPage.includes('_ru.') ? 'ru' : 'en';
    }

    /**
     * Определяет ID текущей страницы для выделения активной ссылки
     */
    getCurrentPageId() {
        return NavigationConfig.pageMapping[this.currentPage] || null;
    }

    /**
     * Получает базовое имя страницы без языкового суффикса
     */
    getPageBaseName() {
        return this.currentPage.replace('_ru.html', '.html').replace('.html', '');
    }

    /**
     * Получает подходящую GitHub ссылку для текущей страницы
     */
    getGitHubLink() {
        const baseName = this.getPageBaseName();
        return NavigationConfig.githubLinks[baseName] || NavigationConfig.githubLinks["index"];
    }

    /**
     * Создает навигационные ссылки
     */
    createNavigationLinks() {
        let links = '';
        
        // На главной странице показываем только локальные якорные ссылки
        if (this.currentPageId === 'home') {
            const localLinks = this.createLocalPageLinks();
            if (localLinks) {
                links += localLinks;
            }
        } else {
            // На других страницах показываем глобальные навигационные ссылки
            const config = NavigationConfig.languages[this.currentLanguage];
            const globalLinks = config.links.map(link => {
                const isActive = link.id === this.currentPageId;
                const activeClass = isActive ? ' class="active"' : '';
                return `<li><a href="${link.href}"${activeClass}>${link.text}</a></li>`;
            }).join('\n                    ');
            
            links += globalLinks;
        }

        // Добавляем GitHub ссылку
        const githubLink = this.getGitHubLink();
        links += `\n                    <li><a href="${githubLink}" target="_blank">GitHub</a></li>`;
        
        return links;
    }

    /**
     * Создает локальные якорные ссылки для главной страницы
     */
    createLocalPageLinks() {
        if (this.currentPageId !== 'home') return '';
        
        const isRussian = this.currentLanguage === 'ru';
        const localLinks = isRussian ? [
            { href: "#features", text: "Возможности" },
            { href: "#installation", text: "Установка" },
            { href: "#commands", text: "Команды" }
        ] : [
            { href: "#features", text: "Features" },
            { href: "#installation", text: "Installation" },
            { href: "#commands", text: "Commands" }
        ];

        return localLinks.map(link => 
            `<li><a href="${link.href}">${link.text}</a></li>`
        ).join('\n                    ');
    }

    /**
     * Создает переключатель языков
     */
    createLanguageSwitcher() {
        const otherLang = this.currentLanguage === 'en' ? 'ru' : 'en';
        const currentConfig = NavigationConfig.languages[this.currentLanguage];
        const otherConfig = NavigationConfig.languages[otherLang];
        
        // Находим соответствующую страницу на другом языке
        const currentLink = currentConfig.links.find(link => link.id === this.currentPageId);
        const otherLink = otherConfig.links.find(link => link.id === this.currentPageId);
        
        const currentHref = currentLink ? currentLink.href : this.currentPage;
        const otherHref = otherLink ? otherLink.href : this.currentPage;

        return `
                    <a href="${currentHref}" class="lang-btn active">${currentConfig.name}</a>
                    <a href="${otherHref}" class="lang-btn">${otherConfig.name}</a>
                `;
    }

    /**
     * Рендерит навигацию в указанный селектор
     */
    renderNavigation(selector = '.nav-links') {
        const navContainer = document.querySelector(selector);
        if (navContainer) {
            // Всегда полная замена навигации
            navContainer.innerHTML = this.createNavigationLinks();
        }

        const langSwitcher = document.querySelector('.language-switcher');
        if (langSwitcher) {
            langSwitcher.innerHTML = this.createLanguageSwitcher();
        }
    }

    /**
     * Инициализация навигации после загрузки DOM
     */
    init() {
        // Ждем загрузки DOM
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => {
                this.renderNavigation();
            });
        } else {
            this.renderNavigation();
        }
    }
}

// Глобальная функция для добавления новых языков
window.addNavigationLanguage = function(langCode, config) {
    NavigationConfig.languages[langCode] = config;
};

// Глобальная функция для обновления GitHub ссылок
window.updateGitHubLinks = function(newLinks) {
    Object.assign(NavigationConfig.githubLinks, newLinks);
};

// Автоматическая инициализация
const navigationManager = new NavigationManager();
navigationManager.init();

// Экспорт для использования в других скриптах
window.NavigationManager = NavigationManager;
window.navigationManager = navigationManager;
