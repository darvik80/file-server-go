// Основной класс для управления файлами с поддержкой авторизации

class FileManager {
    constructor() {
        this.files = [];
        this.filteredFiles = [];
        this.currentSort = 'name';
        this.currentView = 'list';
        this.searchTerm = '';
        this.token = localStorage.getItem('token');
        this.userRole = null;
        this.currentPath = '.'; // Текущая директория

        this.init();
    }

    init() {
        this.bindEvents();
        this.checkAuth();
    }

    bindEvents() {
        // Авторизация
        const loginForm = document.getElementById('loginForm');
        if (loginForm) {
            loginForm.addEventListener('submit', (e) => this.handleLogin(e));
        }

        const logoutBtn = document.getElementById('logoutBtn');
        if (logoutBtn) {
            logoutBtn.addEventListener('click', () => this.logout());
        }

        // Загрузка файлов
        const uploadArea = document.getElementById('uploadArea');
        const fileInput = document.getElementById('fileInput');

        if (uploadArea && fileInput) {
            // Drag & Drop
            uploadArea.addEventListener('dragover', (e) => {
                e.preventDefault();
                uploadArea.classList.add('dragover');
            });

            uploadArea.addEventListener('dragleave', (e) => {
                e.preventDefault();
                uploadArea.classList.remove('dragover');
            });

            uploadArea.addEventListener('drop', (e) => {
                e.preventDefault();
                uploadArea.classList.remove('dragover');

                const files = Array.from(e.dataTransfer.files);
                this.handleFileUpload(files);
            });

            // Клик по области загрузки
            uploadArea.addEventListener('click', () => {
                fileInput.click();
            });

            // Выбор файлов через input
            fileInput.addEventListener('change', (e) => {
                const files = Array.from(e.target.files);
                this.handleFileUpload(files);
            });
        }

        // Поиск
        const searchInput = document.getElementById('searchInput');
        if (searchInput) {
            searchInput.addEventListener('input', debounce((e) => {
                this.searchTerm = e.target.value.toLowerCase();
                this.filterAndRenderFiles();
            }, 300));
        }

        // Сортировка
        const sortSelect = document.getElementById('sortSelect');
        if (sortSelect) {
            sortSelect.addEventListener('change', (e) => {
                this.currentSort = e.target.value;
                this.filterAndRenderFiles();
            });
        }

        // Обработчик формы создания директории
        const createDirForm = document.getElementById('createDirForm');
        if (createDirForm) {
            createDirForm.addEventListener('submit', (e) => this.handleCreateDirectory(e));
        }
    }

    // Проверка авторизации
    checkAuth() {
        if (this.token) {
            fileAPI.setToken(this.token);
            this.showMainSection();
            this.refreshFiles();
        } else {
            this.showLoginModal();
        }
    }

    // Показать модальное окно логина
    showLoginModal() {
        const loginModal = new bootstrap.Modal(document.getElementById('loginModal'));
        loginModal.show();
    }

    // Показать секцию логина (устаревший метод, оставлен для обратной совместимости)
    showLoginSection() {
        this.showLoginModal();
    }

    // Показать основную секцию
    showMainSection() {
        const mainSection = document.getElementById('mainSection');
        if (mainSection) mainSection.style.display = 'block';
    }

    // Обработка логина
    async handleLogin(e) {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        if (!username || !password) {
            showNotification(translator.get('error'), translator.get('unauthorized'));
            return;
        }

        try {
            const data = await fileAPI.login(username, password);
            if (data.token) {
                this.token = data.token;
                localStorage.setItem('token', this.token);
                fileAPI.setToken(this.token);

                // Получаем информацию о пользователе
                await this.loadUserInfo();

                // Закрываем модальное окно логина
                const loginModal = bootstrap.Modal.getInstance(document.getElementById('loginModal'));
                if (loginModal) {
                    loginModal.hide();
                }

                showNotification(translator.get('loginSuccess'), 'success');
                this.showMainSection();
                this.refreshFiles();
            }
        } catch (error) {
            handleError(error, translator.get('error'));
        }
    }

    // Загрузка информации о пользователе
    async loadUserInfo() {
        try {
            const userInfo = await fileAPI.getUserInfo();
            this.userRole = userInfo.role;

            // Скрываем блок загрузки файлов для пользователей с ролью reader
            this.updateUploadSectionVisibility();

            // Скрываем вкладки Users и Applications для пользователей без прав администратора
            this.updateTabsVisibility();
        } catch (error) {
            console.error('Error loading user info:', error);
            this.userRole = null;
        }
    }

    // Обновление видимости секции загрузки файлов
    updateUploadSectionVisibility() {
        const uploadSection = document.querySelector('.upload-section');
        if (uploadSection) {
            // Скрываем секцию загрузки для пользователей с ролью reader
            if (this.userRole === 'reader') {
                uploadSection.style.display = 'none';
            } else {
                uploadSection.style.display = 'block';
            }
        }
    }

    // Обновление видимости вкладок
    updateTabsVisibility() {
        const usersTab = document.getElementById('users-tab');
        const appsTab = document.getElementById('apps-tab');
        const usersPane = document.getElementById('users');
        const appsPane = document.getElementById('applications');

        // Скрываем вкладки Users и Applications для пользователей без прав администратора
        if (this.userRole !== 'admin') {
            if (usersTab) usersTab.style.display = 'none';
            if (appsTab) appsTab.style.display = 'none';
            // Также скрываем соответствующие панели
            if (usersPane) usersPane.style.display = 'none';
            if (appsPane) appsPane.style.display = 'none';

            // Если текущая активная вкладка была скрыта, переключаемся на вкладку файлов
            const activeTab = document.querySelector('.nav-link.active');
            if (activeTab && (activeTab.id === 'users-tab' || activeTab.id === 'apps-tab')) {
                const filesTab = document.getElementById('files-tab');
                const filesPane = document.getElementById('files');
                if (filesTab) filesTab.classList.add('active');
                if (filesPane) filesPane.classList.add('show', 'active');
            }
        } else {
            // Показываем вкладки для администраторов
            if (usersTab) usersTab.style.display = 'block';
            if (appsTab) appsTab.style.display = 'block';
        }
    }

    // Выход из системы
    logout() {
        localStorage.removeItem('token');
        this.token = null;
        this.userRole = null;
        fileAPI.setToken(null);
        this.showLoginSection();
        showNotification(translator.get('logoutSuccess'), 'info');
    }

    // Обновление списка файлов
    async refreshFiles() {
        if (!this.token) {
            this.showLoginSection();
            return;
        }

        try {
            showLoading(document.getElementById('filesList'));
            this.files = await fileAPI.getFiles(this.currentPath);
            this.filterAndRenderFiles();
            this.updateBreadcrumb(); // Обновляем навигационную цепочку
        } catch (error) {
            if (error.message === 'unauthorized') {
                this.logout();
                return;
            }
            handleError(error, translator.get('refreshError'));
            document.getElementById('filesList').innerHTML = `<div class="text-center py-4">${translator.get('error')}</div>`;
        }
    }

    // Обновление навигационной цепочки (breadcrumb)
    updateBreadcrumb() {
        const breadcrumb = document.getElementById('breadcrumb');
        if (!breadcrumb) return;

        // Очищаем breadcrumb
        breadcrumb.innerHTML = '';

        // Если мы в корневой директории, показываем только "Файлы"
        if (this.currentPath === '.' || this.currentPath === '/') {
            breadcrumb.innerHTML = `<li class="breadcrumb-item active" aria-current="page">${translator.get('files')}</li>`;
            return;
        }

        // Разбиваем путь на части
        const pathParts = this.currentPath.split('/');

        // Добавляем ссылку на корневую директорию
        const rootItem = document.createElement('li');
        rootItem.className = 'breadcrumb-item';
        rootItem.innerHTML = `<a href="#" onclick="fileManager.navigateToFolder('.')">${translator.get('files')}</a>`;
        breadcrumb.appendChild(rootItem);

        // Добавляем промежуточные директории
        let pathSoFar = '';
        for (let i = 0; i < pathParts.length; i++) {
            if (pathParts[i] === '') continue;

            if (pathSoFar === '') {
                pathSoFar = pathParts[i];
            } else {
                pathSoFar += '/' + pathParts[i];
            }

            const item = document.createElement('li');
            item.className = 'breadcrumb-item';

            // Для последнего элемента делаем его активным (без ссылки)
            if (i === pathParts.length - 1) {
                item.className += ' active';
                item.setAttribute('aria-current', 'page');
                item.textContent = pathParts[i];
            } else {
                item.innerHTML = `<a href="#" onclick="fileManager.navigateToFolder('${pathSoFar}')">${pathParts[i]}</a>`;
            }

            breadcrumb.appendChild(item);
        }
    }

    // Переход в папку
    async navigateToFolder(path) {
        this.currentPath = path || '.';
        await this.refreshFiles();
    }

    // Загрузка файлов
    async handleFileUpload(files) {
        // Проверяем роль пользователя перед загрузкой
        if (this.userRole === 'reader') {
            showNotification(translator.get('uploadNotAllowed'), 'error');
            return;
        }

        if (files.length === 0) return;

        if (!this.token) {
            this.showLoginSection();
            return;
        }

        // Отладочный вывод для проверки значения currentPath
        console.log('Current path before upload:', this.currentPath);

        const progressContainer = document.getElementById('uploadProgress');
        const progressFill = document.getElementById('progressFill');
        const progressText = document.getElementById('progressText');

        progressContainer.style.display = 'block';

        for (let i = 0; i < files.length; i++) {
            const file = files[i];

            try {
                // Проверка размера файла
                if (file.size > CONSTANTS.MAX_FILE_SIZE) {
                    throw new Error(translator.get('fileTooLarge', {filename: file.name, maxSize: formatFileSize(CONSTANTS.MAX_FILE_SIZE)}));
                }

                progressText.textContent = translator.get('uploadProgress', {filename: file.name, current: i + 1, total: files.length});

                // Передаем текущий путь в функцию загрузки
                await fileAPI.uploadFile(file, this.currentPath, (progress) => {
                    progressFill.style.width = `${progress}%`;
                });

                showNotification(translator.get('uploadSuccess', {filename: file.name}), 'success');
            } catch (error) {
                if (error.message === 'unauthorized') {
                    this.logout();
                    return;
                }
                handleError(error, translator.get('uploadError', {filename: file.name}));
            }
        }

        // Скрываем прогресс и обновляем список
        progressContainer.style.display = 'none';
        progressFill.style.width = '0%';
        progressText.textContent = '0%';

        // Очищаем input
        const fileInput = document.getElementById('fileInput');
        if (fileInput) fileInput.value = '';

        // Обновляем список файлов
        await this.refreshFiles();
    }

    // Фильтрация и рендеринг
    filterAndRenderFiles() {
        // Фильтрация по поиску
        this.filteredFiles = this.files.filter(file => {
            if (!this.searchTerm) return true;
            return file.name.toLowerCase().includes(this.searchTerm);
        });

        // Сортировка
        this.filteredFiles.sort((a, b) => {
            switch (this.currentSort) {
                case 'name':
                    return a.name.localeCompare(b.name);
                case 'size':
                    return (b.size || 0) - (a.size || 0);
                case 'date':
                    return new Date(b.modified || 0) - new Date(a.modified || 0);
                default:
                    return 0;
            }
        });

        this.renderFiles();
    }

    // Рендеринг списка файлов
    renderFiles() {
        const container = document.getElementById('filesList');

        if (this.filteredFiles.length === 0) {
            container.innerHTML = this.searchTerm ?
                `<div class="text-center py-4">${translator.get('noFilesFound')}</div>` :
                `<div class="text-center py-4">${translator.get('noFiles')}</div>`;
            return;
        }

        // Устанавливаем класс контейнера в зависимости от режима отображения
        if (this.currentView === 'grid') {
            container.className = 'row g-3';
        } else {
            container.className = 'files-list';
        }

        const filesHTML = this.filteredFiles.map(file => this.createFileHTML(file)).join('');
        container.innerHTML = filesHTML;

        // Применяем переводы к новым элементам
        translator.applyTranslations();
    }

    // Создание HTML для файла
    createFileHTML(file) {
        const icon = PathUtils.getFileIcon(file.name, file.isDir);
        const size = formatFileSize(file.size || 0);
        const canPreviewFile = canPreview(file.name);

        if (this.currentView === 'grid') {
            // Режим сетки
            return `
                <div class="file-item col-md-4 mb-3">
                    <div class="card h-100">
                        <div class="card-body d-flex flex-column">
                            <div class="text-center mb-3">
                                <div class="file-icon" style="font-size: 2rem;">${icon}</div>
                            </div>
                            <h5 class="card-title text-truncate">${file.name}</h5>
                            <p class="card-text text-muted">${size}</p>
                            <div class="mt-auto">
                                <div class="file-actions d-flex justify-content-center gap-2">
                                    ${file.name === '..' ? 
                                        `<button class="btn btn-secondary btn-sm" onclick="fileManager.navigateToFolder('${file.path}')">${translator.get('back')}</button>` :
                                        file.isDir ? 
                                            `<button class="btn btn-primary btn-sm" onclick="fileManager.navigateToFolder('${file.path}')">${translator.get('open')}</button>` : 
                                            `${canPreviewFile ? `<button class="btn btn-secondary btn-sm" onclick="fileManager.previewFile('${file.name}')">${translator.get('preview')}</button>` : ''}
                                            <button class="btn btn-primary btn-sm" onclick="fileManager.downloadFile('${file.name}')">${translator.get('download')}</button>`}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            `;
        } else {
            // Режим списка (по умолчанию)
            return `
                <div class="file-item">
                    <div class="file-icon">${icon}</div>
                    <div class="file-info">
                        <div class="file-name">${file.name}</div>
                        <div class="file-meta">${size}</div>
                    </div>
                    <div class="file-actions">
                        ${file.name === '..' ? 
                            `<button class="btn btn-secondary" onclick="fileManager.navigateToFolder('${file.path}')">${translator.get('back')}</button>` :
                            file.isDir ? 
                                `<button class="btn btn-primary" onclick="fileManager.navigateToFolder('${file.path}')">${translator.get('open')}</button>` : 
                                `${canPreviewFile ? `<button class="btn btn-secondary" onclick="fileManager.previewFile('${file.name}')">${translator.get('preview')}</button>` : ''}
                                <button class="btn btn-primary" onclick="fileManager.downloadFile('${file.name}')">${translator.get('download')}</button>`}
                    </div>
                </div>
            `;
        }
    }

    // Скачивание файла
    async downloadFile(filename) {
        if (!this.token) {
            this.showLoginSection();
            return;
        }

        try {
            showNotification(translator.get('downloadSuccess', {filename: filename}), 'info');

            const blob = await fileAPI.downloadFile(filename);

            // Создаем ссылку для скачивания
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.download = filename;

            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);

            window.URL.revokeObjectURL(url);

            showNotification(translator.get('downloadSuccess', {filename: filename}), 'success');
        } catch (error) {
            if (error.message === 'unauthorized') {
                this.logout();
                return;
            }
            handleError(error, translator.get('downloadError', {filename: filename}));
        }
    }

    // Предварительный просмотр файла
    async previewFile(filename) {
        if (!this.token) {
            this.showLoginSection();
            return;
        }

        try {
            // Проверяем, можно ли предварительно просматривать файл
            if (!canPreview(filename)) {
                showNotification(translator.get('previewNotAvailable'), 'error');
                return;
            }

            // Получаем расширение файла
            const ext = PathUtils.getFileExtension(filename).toLowerCase();

            // Устанавливаем заголовок модального окна
            const previewTitle = document.getElementById('previewTitle');
            if (previewTitle) {
                previewTitle.textContent = translator.get('previewTitle', {filename: filename});
            }

            // Получаем содержимое файла
            const blob = await fileAPI.downloadFile(filename);

            // Определяем тип контента
            const contentType = blob.type || 'application/octet-stream';

            // Получаем элемент контента модального окна
            const previewContent = document.getElementById('previewContent');
            if (!previewContent) return;

            // Очищаем содержимое
            previewContent.innerHTML = '';
            previewContent.className = ''; // Убираем предыдущие классы

            // Обрабатываем разные типы файлов
            if (contentType.startsWith('image/')) {
                // Для изображений
                const img = document.createElement('img');
                img.src = URL.createObjectURL(blob);
                img.className = 'img-fluid';
                img.style.maxWidth = '100%';
                img.style.height = 'auto';
                img.style.maxHeight = '70vh';
                img.style.objectFit = 'contain';
                img.onload = () => URL.revokeObjectURL(img.src);
                previewContent.appendChild(img);
            } else if (contentType === 'text/html' || ext === 'html') {
                // Для HTML файлов
                const text = await blob.text();
                // Создаем iframe для безопасного отображения HTML
                const iframe = document.createElement('iframe');
                iframe.style.width = '100%';
                iframe.style.height = '70vh';
                iframe.style.border = '1px solid #ddd';
                iframe.style.borderRadius = '5px';
                iframe.sandbox = 'allow-scripts allow-same-origin'; // Ограничиваем возможности iframe для безопасности
                previewContent.appendChild(iframe);

                // После добавления iframe в DOM, записываем в него содержимое
                setTimeout(() => {
                    try {
                        const doc = iframe.contentDocument || iframe.contentWindow.document;
                        doc.open();
                        doc.write(`
                            <!DOCTYPE html>
                            <html>
                            <head>
                                <meta charset="UTF-8">
                                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                                <style>
                                    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif; margin: 0; padding: 20px; }
                                </style>
                            </head>
                            <body>${text}</body>
                            </html>
                        `);
                        doc.close();
                    } catch (e) {
                        // Если не удалось записать в iframe, отображаем как текст
                        previewContent.innerHTML = `<pre class="preview-text" style="white-space: pre-wrap; word-break: break-word; max-height: 70vh; overflow: auto;">${escapeHtml(text)}</pre>`;
                    }
                }, 100);
            } else if (contentType === 'text/markdown' || ext === 'md') {
                // Для Markdown файлов
                const text = await blob.text();
                previewContent.innerHTML = marked.parse(text);
                previewContent.className = 'markdown-body';
            } else if (contentType.startsWith('text/') ||
                ['json', 'js', 'css', 'html', 'txt'].includes(ext)) {
                // Для текстовых файлов
                const text = await blob.text();
                const pre = document.createElement('pre');
                pre.className = 'preview-text';
                pre.textContent = text;
                pre.style.whiteSpace = 'pre-wrap';
                pre.style.wordBreak = 'break-word';
                pre.style.maxHeight = '70vh';
                pre.style.overflow = 'auto';
                pre.style.padding = '15px';
                pre.style.backgroundColor = '#f8f9fa';
                pre.style.borderRadius = '5px';
                pre.style.border = '1px solid #ddd';
                previewContent.appendChild(pre);
            } else {
                // Для других типов файлов
                previewContent.innerHTML = `<p>${translator.get('previewNotAvailable')}</p>`;
            }

            // Показываем модальное окно
            const modal = new bootstrap.Modal(document.getElementById('previewModal'));
            modal.show();
        } catch (error) {
            if (error.message === 'unauthorized') {
                this.logout();
                return;
            }
            handleError(error, translator.get('previewError'));
        }
    }

    // Закрытие предварительного просмотра
    closePreview() {
        const modal = bootstrap.Modal.getInstance(document.getElementById('previewModal'));
        if (modal) {
            modal.hide();
        }

        // Очищаем содержимое при закрытии
        const previewContent = document.getElementById('previewContent');
        if (previewContent) {
            previewContent.innerHTML = '';
        }
    }

    // Показать модальное окно создания директории
    showCreateDirectoryModal() {
        const modal = new bootstrap.Modal(document.getElementById('createDirModal'));
        modal.show();
    }

    // Обработка создания директории
    async handleCreateDirectory(e) {
        e.preventDefault();

        const dirname = document.getElementById('newDirName').value;
        const path = document.getElementById('newDirPath').value;

        if (!dirname) {
            showNotification(translator.get('fillAllFields'), 'error');
            return;
        }

        try {
            const requestData = {
                dirname: dirname
            };

            // Добавляем путь, если он задан
            if (path) {
                requestData.path = path;
            }

            const response = await fetch('/create-dir', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${this.token}`
                },
                body: JSON.stringify(requestData)
            });

            if (!response.ok) {
                if (response.status === 401) {
                    this.logout();
                    return;
                }
                throw new Error('Ошибка при создании директории');
            }

            const data = await response.json();
            showNotification(data.message || translator.get('dirCreated'), 'success');

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('createDirModal'));
            if (modal) {
                modal.hide();
            }

            // Очищаем форму
            document.getElementById('createDirForm').reset();

            // Обновляем список файлов
            this.refreshFiles();
        } catch (error) {
            handleError(error, translator.get('createDirError'));
        }
    }

    // Создание директории
    async createDirectory(path) {
        // Показываем модальное окно для ввода имени новой директории
        this.showCreateDirectoryModal();

        // Устанавливаем путь в скрытое поле формы
        const pathInput = document.getElementById('newDirPath');
        if (pathInput) {
            pathInput.value = path;
        }
    }
}
// Вспомогательная функция для экранирования HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}