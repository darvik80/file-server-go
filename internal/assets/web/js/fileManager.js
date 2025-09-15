// Основной класс для управления файлами с поддержкой авторизации

class FileManager {
    constructor() {
        this.files = [];
        this.filteredFiles = [];
        this.currentSort = 'name';
        this.currentView = 'list';
        this.searchTerm = '';
        this.token = localStorage.getItem('token');

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

    // Выход из системы
    logout() {
        localStorage.removeItem('token');
        this.token = null;
        fileAPI.setToken(null);
        this.showLoginSection();
        showNotification(translator.get('logoutSuccess'), 'info');
    }

    // Загрузка файлов
    async refreshFiles() {
        if (!this.token) {
            this.showLoginSection();
            return;
        }

        try {
            showLoading(document.getElementById('filesList'));
            this.files = await fileAPI.getFiles();
            this.filterAndRenderFiles();
        } catch (error) {
            if (error.message === 'unauthorized') {
                this.logout();
                return;
            }
            handleError(error, translator.get('refreshError'));
            document.getElementById('filesList').innerHTML = `<div class="text-center py-4">${translator.get('error')}</div>`;
        }
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
        const icon = PathUtils.getFileIcon(file.name, false);
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
                                    ${canPreviewFile ? `<button class="btn btn-secondary btn-sm" onclick="fileManager.previewFile('${file.name}')">${translator.get('preview')}</button>` : ''}
                                    <button class="btn btn-primary btn-sm" onclick="fileManager.downloadFile('${file.name}')">${translator.get('download')}</button>
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
                        ${canPreviewFile ? `<button class="btn btn-secondary" onclick="fileManager.previewFile('${file.name}')">${translator.get('preview')}</button>` : ''}
                        <button class="btn btn-primary" onclick="fileManager.downloadFile('${file.name}')">${translator.get('download')}</button>
                    </div>
                </div>
            `;
        }
    }

    // Загрузка файлов
    async handleFileUpload(files) {
        if (files.length === 0) return;

        if (!this.token) {
            this.showLoginSection();
            return;
        }

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

                await fileAPI.uploadFile(file, (progress) => {
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
            const modal = document.getElementById('previewModal');
            const title = document.getElementById('previewTitle');
            const content = document.getElementById('previewContent');

            title.textContent = translator.get('previewTitle', {filename: filename});
            content.innerHTML = `<div class="text-center py-4">${translator.get('loading')}</div>`;

            // Показываем модальное окно с анимацией Bootstrap
            const bootstrapModal = new bootstrap.Modal(modal);
            bootstrapModal.show();

            const ext = PathUtils.getFileExtension(filename);

            if (['jpg', 'jpeg', 'png', 'gif'].includes(ext)) {
                // Предварительный просмотр изображений
                try {
                    const blob = await fileAPI.downloadFile(filename);
                    const imageUrl = URL.createObjectURL(blob);
                    content.innerHTML = `<img src="${imageUrl}" class="img-fluid" alt="${filename}" style="max-width: 100%; height: auto;">`;

                    // Освобождаем память при закрытии модального окна
                    modal.addEventListener('hidden.bs.modal', () => {
                        URL.revokeObjectURL(imageUrl);
                    }, { once: true });
                } catch (error) {
                    if (error.message === 'unauthorized') {
                        this.logout();
                        return;
                    }
                    throw error;
                }
            } else if (ext === 'md') {
                // Предварительный просмотр Markdown файлов
                const blob = await fileAPI.downloadFile(filename);
                const text = await blob.text();
                const html = marked.parse(text);
                content.innerHTML = `<div class="markdown-body">${html}</div>`;
            } else if (ext === 'html') {
                // Предварительный просмотр HTML файлов
                const blob = await fileAPI.downloadFile(filename);
                const text = await blob.text();
                content.innerHTML = text;
            } else if (['txt', 'json', 'css', 'js'].includes(ext)) {
                // Предварительный просмотр текстовых файлов
                const blob = await fileAPI.downloadFile(filename);
                const text = await blob.text();
                content.innerHTML = `<pre class="bg-light p-3 rounded">${this.escapeHtml(text)}</pre>`;
            } else {
                content.innerHTML = `<p class="text-muted">${translator.get('previewNotAvailable')}</p>`;
            }
        } catch (error) {
            if (error.message === 'unauthorized') {
                this.logout();
                return;
            }
            handleError(error, translator.get('previewError'));
            this.closePreview();
        }
    }

    // Закрытие предварительного просмотра
    closePreview() {
        const modal = document.getElementById('previewModal');
        const bootstrapModal = bootstrap.Modal.getInstance(modal);
        if (bootstrapModal) {
            bootstrapModal.hide();
        } else {
            modal.style.display = 'none';
        }
    }

    // Переключение вида
    toggleView() {
        this.currentView = this.currentView === 'list' ? 'grid' : 'list';
        showNotification(translator.get('viewToggle', {view: this.currentView === 'list' ? translator.get('listView') : translator.get('gridView')}), 'info');

        // Обновляем отображение файлов
        this.renderFiles();
    }

    // Экранирование HTML
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Создаем глобальный экземпляр менеджера файлов
window.fileManager = null;