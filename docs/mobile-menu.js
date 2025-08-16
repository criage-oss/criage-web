/**
 * Универсальная функциональность мобильного меню
 * Автоматически инициализируется на всех страницах
 */

class MobileMenu {
    constructor() {
        this.isOpen = false;
        this.init();
    }

    /**
     * Инициализация мобильного меню
     */
    init() {
        // Проверяем наличие элементов мобильного меню
        this.toggle = document.querySelector('.mobile-menu-toggle');
        this.nav = document.querySelector('.mobile-nav');
        
        if (!this.toggle || !this.nav) {
            // Создаем элементы если их нет
            this.createMobileMenuElements();
        }

        this.bindEvents();
    }

    /**
     * Создает элементы мобильного меню если их нет
     */
    createMobileMenuElements() {
        const navContent = document.querySelector('.nav-content');
        if (!navContent) return;

        // Создаем кнопку переключения если её нет
        if (!this.toggle) {
            this.toggle = document.createElement('button');
            this.toggle.className = 'mobile-menu-toggle';
            this.toggle.innerHTML = `
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                    <line x1="3" y1="6" x2="21" y2="6"></line>
                    <line x1="3" y1="12" x2="21" y2="12"></line>
                    <line x1="3" y1="18" x2="21" y2="18"></line>
                </svg>
            `;
            navContent.appendChild(this.toggle);
        }

        // Создаем мобильную навигацию если её нет
        if (!this.nav) {
            const nav = document.querySelector('nav');
            if (nav) {
                this.nav = document.createElement('div');
                this.nav.className = 'mobile-nav';
                this.nav.id = 'mobileNav';
                
                const navLinks = document.querySelector('.nav-links');
                if (navLinks) {
                    this.nav.innerHTML = `<ul class="nav-links">${navLinks.innerHTML}</ul>`;
                }
                
                nav.appendChild(this.nav);
            }
        }
    }

    /**
     * Привязка событий
     */
    bindEvents() {
        // Клик по кнопке переключения
        if (this.toggle) {
            this.toggle.addEventListener('click', (e) => {
                e.stopPropagation();
                this.toggle();
            });
        }

        // Клик вне меню для закрытия
        document.addEventListener('click', (e) => {
            if (this.isOpen && 
                this.nav && 
                !this.nav.contains(e.target) && 
                this.toggle && 
                !this.toggle.contains(e.target)) {
                this.close();
            }
        });

        // Клик по ссылкам в мобильном меню
        if (this.nav) {
            const links = this.nav.querySelectorAll('a');
            links.forEach(link => {
                link.addEventListener('click', () => {
                    this.close();
                });
            });
        }

        // Закрытие при изменении размера экрана
        window.addEventListener('resize', () => {
            if (window.innerWidth > 768 && this.isOpen) {
                this.close();
            }
        });

        // Закрытие по Escape
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && this.isOpen) {
                this.close();
            }
        });
    }

    /**
     * Переключение состояния меню
     */
    toggle() {
        if (this.isOpen) {
            this.close();
        } else {
            this.open();
        }
    }

    /**
     * Открытие меню
     */
    open() {
        if (!this.nav) return;
        
        this.nav.classList.add('open');
        this.isOpen = true;
        
        // Блокируем прокрутку страницы
        document.body.style.overflow = 'hidden';
        
        // Устанавливаем фокус на первую ссылку
        const firstLink = this.nav.querySelector('a');
        if (firstLink) {
            firstLink.focus();
        }

        // Анимация кнопки
        if (this.toggle) {
            this.toggle.setAttribute('aria-expanded', 'true');
            this.animateToggleButton(true);
        }
    }

    /**
     * Закрытие меню
     */
    close() {
        if (!this.nav) return;
        
        this.nav.classList.remove('open');
        this.isOpen = false;
        
        // Восстанавливаем прокрутку страницы
        document.body.style.overflow = '';
        
        // Возвращаем фокус на кнопку
        if (this.toggle) {
            this.toggle.focus();
            this.toggle.setAttribute('aria-expanded', 'false');
            this.animateToggleButton(false);
        }
    }

    /**
     * Анимация кнопки переключения
     */
    animateToggleButton(isOpen) {
        if (!this.toggle) return;
        
        const svg = this.toggle.querySelector('svg');
        if (svg) {
            if (isOpen) {
                // Анимация в иконку X
                svg.innerHTML = `
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                `;
            } else {
                // Анимация в гамбургер
                svg.innerHTML = `
                    <line x1="3" y1="6" x2="21" y2="6"></line>
                    <line x1="3" y1="12" x2="21" y2="12"></line>
                    <line x1="3" y1="18" x2="21" y2="18"></line>
                `;
            }
        }
    }

    /**
     * Обновление содержимого мобильного меню
     */
    updateContent() {
        if (!this.nav) return;
        
        const mainNavLinks = document.querySelector('.nav-content .nav-links');
        const mobileNavLinks = this.nav.querySelector('.nav-links');
        
        if (mainNavLinks && mobileNavLinks) {
            mobileNavLinks.innerHTML = mainNavLinks.innerHTML;
            this.bindMobileLinkEvents();
        }
    }

    /**
     * Привязка событий к ссылкам мобильного меню
     */
    bindMobileLinkEvents() {
        if (!this.nav) return;
        
        const links = this.nav.querySelectorAll('a');
        links.forEach(link => {
            link.addEventListener('click', () => {
                this.close();
            });
        });
    }
}

// Глобальные функции для обратной совместимости
window.toggleMobileMenu = function() {
    if (window.mobileMenu) {
        window.mobileMenu.toggle();
    }
};

window.closeMobileMenu = function() {
    if (window.mobileMenu) {
        window.mobileMenu.close();
    }
};

// Автоматическая инициализация
document.addEventListener('DOMContentLoaded', () => {
    window.mobileMenu = new MobileMenu();
});

// Обновление при изменении навигации (для интеграции с navigation.js)
window.addEventListener('navigationUpdated', () => {
    if (window.mobileMenu) {
        window.mobileMenu.updateContent();
    }
});

// Экспорт для использования в других скриптах
window.MobileMenu = MobileMenu;
