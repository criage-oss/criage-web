/**
 * Упрощенная навигация для GitHub Pages
 * Работает без сложных зависимостей
 */

(function() {
    'use strict';
    
    // Конфигурация навигации
    const NAV_CONFIG = {
        en: {
            links: [
                { href: "index.html", text: "Home", active: ["index.html", ""] },
                { href: "docs.html", text: "Client Docs", active: ["docs.html"] },
                { href: "repository-server.html", text: "Repository Server", active: ["repository-server.html"] },
                { href: "mcp-server.html", text: "MCP Server", active: ["mcp-server.html"] },
                { href: "cicd-examples.html", text: "CI/CD Examples", active: ["cicd-examples.html"] },
                { href: "https://github.com/criage-oss/criage-client", text: "GitHub", external: true }
            ],
            langSwitcher: { current: "EN", other: "RU", otherHref: "_ru.html" }
        },
        ru: {
            links: [
                { href: "index_ru.html", text: "Главная", active: ["index_ru.html"] },
                { href: "docs_ru.html", text: "Документация клиента", active: ["docs_ru.html"] },
                { href: "repository-server_ru.html", text: "Сервер репозитория", active: ["repository-server_ru.html"] },
                { href: "mcp-server_ru.html", text: "MCP Server", active: ["mcp-server_ru.html"] },
                { href: "cicd-examples_ru.html", text: "Примеры CI/CD", active: ["cicd-examples_ru.html"] },
                { href: "https://github.com/criage-oss/criage-client", text: "GitHub", external: true }
            ],
            langSwitcher: { current: "RU", other: "EN", otherHref: ".html" }
        }
    };

    // Определяем текущий язык и страницу
    function getCurrentLanguage() {
        const path = window.location.pathname;
        return path.includes('_ru') ? 'ru' : 'en';
    }

    function getCurrentPage() {
        const path = window.location.pathname;
        return path.split('/').pop() || 'index.html';
    }

    // Создаем HTML для навигации
    function createNavigationHTML(config, currentPage) {
        const links = config.links.map(link => {
            const isActive = link.active && link.active.some(activePage => 
                currentPage === activePage || (activePage === "" && currentPage === "index.html")
            );
            
            const classes = ['nav-link'];
            if (isActive) classes.push('active');
            if (link.external) classes.push('github-link');

            const target = link.external ? ' target="_blank"' : '';
            const githubIcon = link.external ? 
                '<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>' 
                : '';

            return `<li><a href="${link.href}" class="${classes.join(' ')}"${target}>${githubIcon}${link.text}</a></li>`;
        }).join('');

        return links;
    }

    // Создаем переключатель языков
    function createLanguageSwitcherHTML(config, currentPage) {
        const ls = config.langSwitcher;
        const otherHref = currentPage.replace(/(_ru)?\.html$/, ls.otherHref);
        
        return `
            <a href="${currentPage}" class="lang-btn active">${ls.current}</a>
            <a href="${otherHref}" class="lang-btn">${ls.other}</a>
        `;
    }

    // Инициализация навигации
    function initNavigation() {
        console.log('🚀 Simple Navigation: Initializing...');
        
        try {
            const lang = getCurrentLanguage();
            const currentPage = getCurrentPage();
            const config = NAV_CONFIG[lang];
            
            console.log(`📍 Language: ${lang}, Page: ${currentPage}`);

            // Основная навигация
            const navLinks = document.querySelector('.nav-links');
            if (navLinks) {
                navLinks.innerHTML = createNavigationHTML(config, currentPage);
                console.log('✅ Main navigation created');
            }

            // Мобильная навигация
            const mobileNavLinks = document.querySelector('.mobile-nav .nav-links');
            if (mobileNavLinks) {
                mobileNavLinks.innerHTML = createNavigationHTML(config, currentPage);
                console.log('✅ Mobile navigation created');
            }

            // Переключатель языков
            const langSwitcher = document.querySelector('.language-switcher');
            if (langSwitcher) {
                langSwitcher.innerHTML = createLanguageSwitcherHTML(config, currentPage);
                console.log('✅ Language switcher created');
            }

            console.log('🎉 Simple Navigation: Initialization complete!');
            
        } catch (error) {
            console.error('❌ Simple Navigation: Error during initialization:', error);
        }
    }

    // Функция для мобильного меню
    function setupMobileMenu() {
        const toggle = document.querySelector('.mobile-menu-toggle');
        const nav = document.querySelector('.mobile-nav');
        
        if (!toggle || !nav) return;

        toggle.addEventListener('click', function(e) {
            e.stopPropagation();
            nav.classList.toggle('open');
            
            // Анимация иконки
            const svg = toggle.querySelector('svg');
            if (svg) {
                if (nav.classList.contains('open')) {
                    svg.innerHTML = '<line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line>';
                } else {
                    svg.innerHTML = '<line x1="3" y1="6" x2="21" y2="6"></line><line x1="3" y1="12" x2="21" y2="12"></line><line x1="3" y1="18" x2="21" y2="18"></line>';
                }
            }
        });

        // Закрытие при клике вне меню
        document.addEventListener('click', function(e) {
            if (!nav.contains(e.target) && !toggle.contains(e.target)) {
                nav.classList.remove('open');
                const svg = toggle.querySelector('svg');
                if (svg) {
                    svg.innerHTML = '<line x1="3" y1="6" x2="21" y2="6"></line><line x1="3" y1="12" x2="21" y2="12"></line><line x1="3" y1="18" x2="21" y2="18"></line>';
                }
            }
        });

        console.log('📱 Mobile menu setup complete');
    }

    // Глобальная функция для обратной совместимости
    window.toggleMobileMenu = function() {
        const nav = document.querySelector('.mobile-nav');
        if (nav) {
            nav.classList.toggle('open');
        }
    };

    // Инициализация при загрузке DOM
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function() {
            initNavigation();
            setupMobileMenu();
        });
    } else {
        // DOM уже загружен
        initNavigation();
        setupMobileMenu();
    }

    // Дополнительная инициализация после полной загрузки
    window.addEventListener('load', function() {
        // Двойная проверка через небольшую задержку
        setTimeout(function() {
            if (document.querySelector('.nav-links') && 
                document.querySelector('.nav-links').children.length === 0) {
                console.log('🔄 Simple Navigation: Re-initializing...');
                initNavigation();
            }
        }, 100);
    });

    console.log('📦 Simple Navigation script loaded');
})();
