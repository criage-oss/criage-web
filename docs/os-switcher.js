/**
 * Универсальный переключатель ОС для инструкций установки
 * Автоматически инициализируется на всех страницах
 */

class OSSwitcher {
    constructor() {
        this.currentOS = this.detectOS();
        this.translations = {
            en: {
                linux: '🐧 Linux',
                macos: '🍎 macOS', 
                windows: '🪟 Windows',
                copy: 'Copy',
                copied: 'Copied!'
            },
            ru: {
                linux: '🐧 Linux',
                macos: '🍎 macOS',
                windows: '🪟 Windows', 
                copy: 'Копировать',
                copied: 'Скопировано!'
            }
        };
        this.currentLang = this.detectLanguage();
    }

    /**
     * Определяет ОС пользователя
     */
    detectOS() {
        const userAgent = navigator.userAgent.toLowerCase();
        if (userAgent.includes('win')) return 'windows';
        if (userAgent.includes('mac')) return 'macos';
        return 'linux';
    }

    /**
     * Определяет язык страницы
     */
    detectLanguage() {
        const path = window.location.pathname;
        return path.includes('_ru.') ? 'ru' : 'en';
    }

    /**
     * Создает HTML для переключателя ОС
     */
    createOSSwitch(config = {}) {
        const lang = this.translations[this.currentLang];
        const defaultConfig = {
            showLinux: true,
            showMacOS: true,
            showWindows: true,
            defaultOS: this.currentOS
        };
        const options = { ...defaultConfig, ...config };

        let tabs = '';
        if (options.showLinux) tabs += `<button class="os-tab ${options.defaultOS === 'linux' ? 'active' : ''}" data-os="linux">${lang.linux}</button>`;
        if (options.showMacOS) tabs += `<button class="os-tab ${options.defaultOS === 'macos' ? 'active' : ''}" data-os="macos">${lang.macos}</button>`;
        if (options.showWindows) tabs += `<button class="os-tab ${options.defaultOS === 'windows' ? 'active' : ''}" data-os="windows">${lang.windows}</button>`;

        return `
            <div class="os-switcher">
                <div class="os-tabs">
                    ${tabs}
                </div>
            </div>
        `;
    }

    /**
     * Создает контейнер для кода под конкретную ОС
     */
    createCodeBlock(os, code, isActive = false) {
        const lang = this.translations[this.currentLang];
        return `
            <div class="code-snippet os-content ${isActive ? 'active' : ''}" data-os="${os}">
                <button class="copy-code" onclick="copyToClipboard(this)">${lang.copy}</button>
                <pre>${code}</pre>
            </div>
        `;
    }

    /**
     * Инициализирует переключатели на странице
     */
    init() {
        // Автоматически находим и инициализируем все переключатели
        document.querySelectorAll('.os-switcher').forEach(switcher => {
            this.initSingleSwitcher(switcher);
        });

        // Автоматическая замена простых блоков установки на переключатели
        this.autoConvertInstallBlocks();
    }

    /**
     * Инициализирует один переключатель
     */
    initSingleSwitcher(switcher) {
        const tabs = switcher.querySelectorAll('.os-tab');
        const contents = switcher.querySelectorAll('.os-content');

        tabs.forEach(tab => {
            tab.addEventListener('click', () => {
                const targetOS = tab.dataset.os;

                // Обновляем активную вкладку
                tabs.forEach(t => t.classList.remove('active'));
                tab.classList.add('active');

                // Обновляем активный контент
                contents.forEach(content => {
                    content.classList.remove('active');
                    if (content.dataset.os === targetOS) {
                        content.classList.add('active');
                    }
                });

                // Сохраняем выбор пользователя
                localStorage.setItem('criage-preferred-os', targetOS);
            });
        });

        // Восстанавливаем сохраненный выбор пользователя
        const savedOS = localStorage.getItem('criage-preferred-os');
        if (savedOS) {
            const savedTab = switcher.querySelector(`[data-os="${savedOS}"]`);
            if (savedTab) {
                savedTab.click();
            }
        }
    }

    /**
     * Автоматическое преобразование простых блоков установки
     */
    autoConvertInstallBlocks() {
        // Ищем блоки с классом 'auto-os-switch'
        document.querySelectorAll('.auto-os-switch').forEach(block => {
            const template = block.dataset.template;
            if (template) {
                this.convertTemplateBlock(block, template);
            }
        });
    }

    /**
     * Преобразует блок с шаблоном
     */
    convertTemplateBlock(block, template) {
        const packageName = block.dataset.package || 'criage-client';
        const lang = this.currentLang === 'ru' ? 'ru' : 'en';
        
        const templates = {
            'install': {
                linux: this.generateLinuxInstall(packageName),
                macos: this.generateMacOSInstall(packageName),
                windows: this.generateWindowsInstall(packageName)
            },
            'docker': {
                linux: this.generateDockerCommands(packageName, 'linux'),
                macos: this.generateDockerCommands(packageName, 'macos'),
                windows: this.generateDockerCommands(packageName, 'windows')
            }
        };

        if (templates[template]) {
            const switcherHTML = this.createOSSwitch();
            const codeBlocks = Object.entries(templates[template])
                .map(([os, code]) => this.createCodeBlock(os, code, os === this.currentOS))
                .join('');

            block.innerHTML = switcherHTML + codeBlocks;
            this.initSingleSwitcher(block.querySelector('.os-switcher'));
        }
    }

    /**
     * Генерирует команды установки для Linux
     */
    generateLinuxInstall(packageName) {
        const repoName = packageName === 'criage-client' ? 'criage-client' : 
                        packageName === 'criage-server' ? 'criage-server' : 'criage-mcp';
        const binaryName = packageName === 'criage-client' ? 'criage' :
                          packageName === 'criage-server' ? 'criage-server' : 'criage-mcp';
        
        return `# Install ${packageName} on Linux
wget https://github.com/criage-oss/${repoName}/releases/latest/download/${binaryName}-linux-amd64.tar.gz
tar -xzf ${binaryName}-linux-amd64.tar.gz
sudo mv ${binaryName}-linux-amd64 /usr/local/bin/${binaryName}

# Verify installation
${binaryName} --version`;
    }

    /**
     * Генерирует команды установки для macOS
     */
    generateMacOSInstall(packageName) {
        const repoName = packageName === 'criage-client' ? 'criage-client' : 
                        packageName === 'criage-server' ? 'criage-server' : 'criage-mcp';
        const binaryName = packageName === 'criage-client' ? 'criage' :
                          packageName === 'criage-server' ? 'criage-server' : 'criage-mcp';
        
        return `# Install ${packageName} on macOS
curl -L https://github.com/criage-oss/${repoName}/releases/latest/download/${binaryName}-darwin-amd64.tar.gz -o ${binaryName}-darwin-amd64.tar.gz
tar -xzf ${binaryName}-darwin-amd64.tar.gz
sudo mv ${binaryName}-darwin-amd64 /usr/local/bin/${binaryName}

# Verify installation
${binaryName} --version`;
    }

    /**
     * Генерирует команды установки для Windows
     */
    generateWindowsInstall(packageName) {
        const repoName = packageName === 'criage-client' ? 'criage-client' : 
                        packageName === 'criage-server' ? 'criage-server' : 'criage-mcp';
        const binaryName = packageName === 'criage-client' ? 'criage' :
                          packageName === 'criage-server' ? 'criage-server' : 'criage-mcp';
        
        return `# Install ${packageName} on Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/criage-oss/${repoName}/releases/latest/download/${binaryName}-windows-amd64.zip" -OutFile "${binaryName}-windows-amd64.zip"
Expand-Archive -Path "${binaryName}-windows-amd64.zip" -DestinationPath "."
Move-Item "${binaryName}-windows-amd64.exe" "$env:USERPROFILE\\AppData\\Local\\Microsoft\\WindowsApps\\${binaryName}.exe"

# Verify installation
${binaryName} --version`;
    }

    /**
     * Генерирует Docker команды
     */
    generateDockerCommands(packageName, os) {
        const containerName = packageName === 'criage-client' ? 'huginnn/client' :
                             packageName === 'criage-server' ? 'huginnn/server' : 'huginnn/mcp';
        
        const osSpecific = os === 'windows' ? 
            '# Windows (PowerShell)\ndocker run --rm ' : 
            `# ${os.charAt(0).toUpperCase() + os.slice(1)}\ndocker run --rm `;
            
        return `${osSpecific}${containerName}:latest --help

# Run with volume mount
docker run --rm -v "\${PWD}:/workspace" ${containerName}:latest`;
    }
}

// Глобальная функция копирования (если не определена)
if (typeof copyToClipboard === 'undefined') {
    window.copyToClipboard = function(btn) {
        const parent = btn.parentElement;
        let text = '';

        if (parent.classList.contains('code-snippet')) {
            text = parent.querySelector('pre').textContent.trim();
        }

        const lang = window.location.pathname.includes('_ru.') ? 'ru' : 'en';
        const copyText = lang === 'ru' ? 'Копировать' : 'Copy';
        const copiedText = lang === 'ru' ? 'Скопировано!' : 'Copied!';

        navigator.clipboard.writeText(text).then(() => {
            const originalText = btn.textContent;
            btn.textContent = copiedText;
            setTimeout(() => {
                btn.textContent = copyText;
            }, 2000);
        });
    };
}

// Автоматическая инициализация
document.addEventListener('DOMContentLoaded', () => {
    window.osSwitcher = new OSSwitcher();
    window.osSwitcher.init();
});

// Экспорт для использования в других скриптах
window.OSSwitcher = OSSwitcher;
