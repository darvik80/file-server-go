import { SELECTORS, domCache, PathUtils } from './constants.js';
import { showError, showToast, formatFileSize } from './utils.js';
import { FileManagerAPI } from './api.js';

export class FileManager {
    constructor() {
        this.currentPath = '.';
        this.token = localStorage.getItem('token');
        this.contextMenuTarget = null;
        this.confirmModal = null;
        this.createDirModal = null;
        this.api = new FileManagerAPI(this.token);
        this.init();
    }

    init() {
        this.cacheDomElements();
        this.bindEvents();
        this.initModals();
        this.checkAuth();
    }

    cacheDomElements() {
        Object.keys(SELECTORS).forEach(key => {
            const selector = SELECTORS[key];
            domCache[selector] = document.querySelector(selector);
        });
    }

    bindEvents() {
        // Логин
        domCache[SELECTORS.LOGIN_FORM]?.addEventListener('submit', (e) => this.handleLogin(e));

        // Логаут
        document.getElementById('logoutBtn')?.addEventListener('click', () => this.logout());

        // Загрузка файлов
        domCache[SELECTORS.UPLOAD_FORM]?.addEventListener('submit', (e) => this.handleUpload(e));

        // Создание директории
        document.getElementById('createDirForm')?.addEventListener('submit', (e) => this.handleCreateDir(e));

        // Подтверждение удаления
        document.getElementById('confirmDelete')?.addEventListener('click', () => this.confirmDelete());

        // Кнопка создания директории
        domCache[SELECTORS.CREATE_DIR_BTN]?.addEventListener('click', () => this.showCreateDirDialog());

        // Контекстное меню
        domCache[SELECTORS.CONTEXT_MENU]?.addEventListener('click', (e) => this.handleContextMenuClick(e));

        // Закрытие контекстного меню
        document.addEventListener('click', (e) => this.closeContextMenu(e));

        // Предотвращение стандартного контекстного меню
        document.addEventListener('contextmenu', (e) => e.preventDefault());

        // Делегирование событий
        document.addEventListener('click', (e) => this.handleEventDelegation(e));
    }

    initModals() {
        this.confirmModal = new bootstrap.Modal(document.getElementById('confirmModal'));
        this.createDirModal = new bootstrap.Modal(document.getElementById('createDirModal'));
    }

    checkAuth() {
        if (this.token) {
            this.showMainSection();
            this.loadFiles(this.currentPath);
        } else {
            this.showLoginSection();
        }
    }

    // Авторизация
    async handleLogin(e) {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        try {
            const data = await this.api.login(username, password);
            if (data.token) {
                this.token = data.token;
                localStorage.setItem('token', this.token);
                this.api = new FileManagerAPI(this.token);
                this.showMainSection();
                this.loadFiles('.');
            }
        } catch (error) {
            showError(error.message);
        }
    }

    logout() {
        localStorage.removeItem('token');
        this.token = null;
        this.api = new FileManagerAPI(null);
        this.showLoginSection();
    }

    // Работа с файлами
    async loadFiles(path = '.') {
        if (!this.token) {
            this.showLoginSection();
            return;
        }

        try {
            const files = await this.api.getFiles(path);
            this.updatePath(path);
            this.renderFiles(files);
        } catch (error) {
            this.handleError(error);
        }
    }

    renderFiles(files) {
        const filesContainer = domCache[SELECTORS.FILES];
        if (!filesContainer) return;

        if (!Array.isArray(files) || files.length === 0) {
            filesContainer.innerHTML = '<p class="text-muted">Нет файлов</p>';
            return;
        }

        const fragment = document.createDocumentFragment();

        files.forEach(file => {
            const div = document.createElement('div');
            div.className = 'd-flex justify-content-between align-items-center border-bottom py-2 file-item';

            if (file.name !== '..') {
                div.setAttribute('data-is-dir', file.isDir);
                div.setAttribute('data-path', file.path || file.name);
                div.setAttribute('oncontextmenu', `fileManager.showContextMenu(event, this)`);
            }

            div.innerHTML = this.getFileHTML(file);
            fragment.appendChild(div);
        });

        filesContainer.innerHTML = '';
        filesContainer.appendChild(fragment);
    }

    getFileHTML(file) {
        const isDirectory = file.isDir;
        const icon = isDirectory ?
            '<i class="bi bi-folder me-2 text-success"></i>' :
            '<i class="bi bi-file-earmark me-2"></i>';

        const size = !isDirectory ? formatFileSize(file.size) : '';

        const filePath = (file.path || file.name).replace(/\\/g, '/');
        const encodedPath = PathUtils.encodePath(filePath);

        return `
            <div>
                ${icon}
                ${isDirectory ?
            `<span class="clickable" onclick="fileManager.loadFiles('${filePath}')">${file.name}</span>` :
            `<a href="#" class="text-decoration-none" onclick="return fileManager.handleDownload(event, '${encodedPath}')">${file.name}</a>`
        }
            </div>
            ${size ? `<span class="text-muted small">${size}</span>` : ''}
        `;
    }

    // Загрузка файлов
    async handleUpload(e) {
        e.preventDefault();

        const formData = new FormData(e.target);
        const progressBar = e.target.querySelector('.progress');
        const progressBarInner = progressBar.querySelector('.progress-bar');
        const submitButton = e.target.querySelector('button[type="submit"]');

        if (!formData.get('path')) {
            formData.set('path', '.');
        }

        progressBar.classList.remove('d-none');
        submitButton.disabled = true;

        try {
            const data = await this.api.uploadFile(formData);
            if (data.message) {
                this.loadFiles(this.currentPath);
                e.target.reset();
                showToast('Успех', 'Файл успешно загружен');
            }
        } catch (error) {
            this.handleError(error);
        } finally {
            progressBar.classList.add('d-none');
            progressBarInner.style.width = '0%';
            submitButton.disabled = false;
        }
    }

    // Скачивание файлов
    async handleDownload(event, filePath) {
        event.preventDefault();

        if (!this.token) {
            this.showLoginSection();
            return false;
        }

        try {
            const blob = await this.api.downloadFile(filePath);
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', decodeURIComponent(PathUtils.getFileName(filePath)));

            document.body.appendChild(link);
            link.click();

            window.URL.revokeObjectURL(url);
            document.body.removeChild(link);
        } catch (error) {
            this.handleError(error);
        }

        return false;
    }

    // Работа с директориями
    async handleCreateDir(e) {
        e.preventDefault();

        const dirname = document.getElementById('dirname').value.trim();
        const path = document.getElementById('createDirPath').value;

        if (!dirname) {
            showError('Введите имя директории');
            return;
        }

        try {
            const data = await this.api.createDirectory(dirname, path);
            if (data.message) {
                this.loadFiles(this.currentPath);
                document.getElementById('dirname').value = '';
                this.createDirModal.hide();
                showToast('Успех', 'Директория успешно создана');
            }
        } catch (error) {
            this.handleError(error);
        }
    }

    showCreateDirDialog() {
        document.getElementById('createDirPath').value = this.currentPath;
        document.getElementById('dirname').value = '';
        this.createDirModal.show();
    }

    // Удаление
    async deleteFile(path) {
        try {
            const data = await this.api.deleteFile(path);
            if (data.message) {
                showToast('Успех', 'Файл успешно удален');
                this.loadFiles(this.currentPath);
            }
        } catch (error) {
            this.handleError(error);
        }
    }

    async deleteDirectory(path) {
        try {
            const data = await this.api.deleteDirectory(path);
            if (data.message) {
                showToast('Успех', 'Директория успешно удалена');
                this.loadFiles(this.currentPath);
            }
        } catch (error) {
            this.handleError(error);
        }
    }

    // Контекстное меню
    showContextMenu(e, item) {
        e.preventDefault();
        this.contextMenuTarget = item;

        const menuItems = domCache[SELECTORS.CONTEXT_MENU].querySelectorAll('.context-menu-item');
        const isDir = item.dataset.isDir === 'true';

        menuItems.forEach(menuItem => {
            menuItem.style.display = isDir || menuItem.dataset.action === 'delete' ? 'flex' : 'none';
        });

        const contextMenu = domCache[SELECTORS.CONTEXT_MENU];
        contextMenu.style.display = 'block';

        // Позиционирование
        const clickX = e.clientX;
        const clickY = e.clientY;
        const screenW = window.innerWidth;
        const screenH = window.innerHeight;
        const menuW = contextMenu.offsetWidth;
        const menuH = contextMenu.offsetHeight;

        contextMenu.style.left = Math.min(clickX, screenW - menuW) + 'px';
        contextMenu.style.top = Math.min(clickY, screenH - menuH) + 'px';
    }

    handleContextMenuClick(e) {
        const action = e.target.closest('[data-action]')?.dataset.action;
        if (!action) return;

        switch(action) {
            case 'create':
                this.showCreateDirDialog();
                this.closeContextMenu();
                break;
            case 'delete':
                if (!this.contextMenuTarget) return;
                this.prepareDelete();
                this.closeContextMenu();
                break;
        }
    }

    prepareDelete() {
        const path = this.contextMenuTarget.dataset.path;
        const name = this.contextMenuTarget.textContent.trim();
        const isDir = this.contextMenuTarget.dataset.isDir === 'true';

        document.getElementById('itemType').textContent = isDir ? 'директорию' : 'файл';
        document.getElementById('itemToDelete').textContent = name;

        document.getElementById('confirmDelete').onclick = () => {
            if (isDir) {
                this.deleteDirectory(path);
            } else {
                this.deleteFile(path);
            }
            this.confirmModal.hide();
        };

        this.confirmModal.show();
    }

    confirmDelete() {
        if (!this.contextMenuTarget) return;
        const path = this.contextMenuTarget.dataset.path;
        const isDir = this.contextMenuTarget.dataset.isDir === 'true';

        if (isDir) {
            this.deleteDirectory(path);
        } else {
            this.deleteFile(path);
        }
    }

    closeContextMenu(e) {
        if (!e || !domCache[SELECTORS.CONTEXT_MENU].contains(e.target)) {
            domCache[SELECTORS.CONTEXT_MENU].style.display = 'none';
        }
    }

    // Делегирование событий
    handleEventDelegation(e) {
        // Обработка кликов по хлебным крошкам
        if (e.target.matches('#currentPath .clickable')) {
            const path = e.target.dataset.path;
            if (path) this.loadFiles(path);
        }
    }

    // Вспомогательные методы
    updatePath(path) {
        this.currentPath = path;
        const parts = path === '.' ? [] : path.split('/');

        let breadcrumbs = '<span class="clickable" data-path="." onclick="fileManager.loadFiles(\'.\')">Корневая директория</span>';
        let currentLink = '';

        for (let i = 0; i < parts.length; i++) {
            currentLink = parts.slice(0, i + 1).join('/');
            breadcrumbs += ' / ';

            if (i === parts.length - 1) {
                breadcrumbs += `<span>${parts[i]}</span>`;
            } else {
                breadcrumbs += `<span class="clickable" data-path="${currentLink}" onclick="fileManager.loadFiles('${currentLink}')">${parts[i]}</span>`;
            }
        }

        document.getElementById('currentPath').innerHTML = breadcrumbs;
        document.getElementById('uploadPath').value = path;
        document.getElementById('createDirPath').value = path;
    }

    showMainSection() {
        domCache[SELECTORS.LOGIN_SECTION].style.display = 'none';
        domCache[SELECTORS.MAIN_SECTION].style.display = 'block';
    }

    showLoginSection() {
        domCache[SELECTORS.LOGIN_SECTION].style.display = 'block';
        domCache[SELECTORS.MAIN_SECTION].style.display = 'none';
    }

    handleError(error) {
        console.error('Error:', error);
        if (error.message === 'unauthorized') {
            this.showLoginSection();
        } else {
            showError(error.message || 'Произошла ошибка');
        }
    }
}