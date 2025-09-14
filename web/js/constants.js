// Константы и селекторы
export const SELECTORS = {
    ERROR_SECTION: '#errorSection',
    ERROR_MESSAGE: '#errorMessage',
    LOGIN_SECTION: '#loginSection',
    MAIN_SECTION: '#mainSection',
    CURRENT_PATH: '#currentPath',
    FILES: '#files',
    UPLOAD_FORM: '#uploadForm',
    LOGIN_FORM: '#loginForm',
    CONTEXT_MENU: '#contextMenu',
    CREATE_DIR_BTN: '#createDirBtn'
};

// Кэш DOM элементов
export const domCache = {};

// Утилиты для работы с путями
export const PathUtils = {
    normalize(path) {
        return decodeURIComponent(path).replace(/\\/g, '/');
    },
    getFileName(path) {
        return path.split('/').pop();
    },
    getParentPath(path) {
        const parts = path.split('/');
        parts.pop();
        return parts.join('/') || '.';
    },
    encodePath(path) {
        return encodeURIComponent(this.normalize(path));
    }
};