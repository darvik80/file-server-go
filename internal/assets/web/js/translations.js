// Система переводов для поддержки нескольких языков
const translations = {
    'ru': {
        // Общие термины
        'title': 'Файловый сервер',
        'login': 'Войти',
        'logout': 'Выйти',
        'username': 'Имя пользователя',
        'password': 'Пароль',
        'upload': 'Загрузить файл',
        'files': 'Файлы',
        'refresh': 'Обновить',
        'view': 'Вид',
        'search': 'Поиск файлов...',
        'sortBy': 'Сортировать по',
        'name': 'Имя',
        'size': 'Размер',
        'date': 'Дата',
        'download': 'Скачать',
        'preview': 'Просмотр',
        'noFiles': 'Нет загруженных файлов',
        'noFilesFound': 'Файлы не найдены',
        'loading': 'Загрузка файлов...',
        'uploadProgress': 'Загрузка {filename}... ({current}/{total})',
        'uploadSuccess': 'Файл "{filename}" успешно загружен',
        'downloadSuccess': 'Файл "{filename}" успешно скачан',
        'previewTitle': 'Предварительный просмотр: {filename}',
        'previewNotAvailable': 'Предварительный просмотр недоступен для этого типа файла.',
        'unauthorized': 'Необходима авторизация',
        'error': 'Ошибка',
        'info': 'Информация',
        'success': 'Успех',
        'uploadAreaText': 'Перетащите файлы сюда или нажмите для выбора',
        'selectFiles': 'Выбрать файлы',
        'loginTitle': 'Авторизация',
        'loginSuccess': 'Успешная авторизация!',
        'logoutSuccess': 'Вы вышли из системы',
        'fileTooLarge': 'Файл "{filename}" слишком большой. Максимальный размер: {maxSize}',
        'uploadError': 'Ошибка загрузки файла "{filename}"',
        'downloadError': 'Ошибка скачивания файла "{filename}"',
        'previewError': 'Ошибка предварительного просмотра',
        'refreshError': 'Ошибка при обновлении списка файлов',
        'viewToggle': 'Переключен вид: {view}',
        'listView': 'Список',
        'gridView': 'Сетка'
    },
    'en': {
        // General terms
        'title': 'File Server',
        'login': 'Login',
        'logout': 'Logout',
        'username': 'Username',
        'password': 'Password',
        'upload': 'Upload File',
        'files': 'Files',
        'refresh': 'Refresh',
        'view': 'View',
        'search': 'Search files...',
        'sortBy': 'Sort by',
        'name': 'Name',
        'size': 'Size',
        'date': 'Date',
        'download': 'Download',
        'preview': 'Preview',
        'noFiles': 'No files uploaded',
        'noFilesFound': 'No files found',
        'loading': 'Loading files...',
        'uploadProgress': 'Uploading {filename}... ({current}/{total})',
        'uploadSuccess': 'File "{filename}" uploaded successfully',
        'downloadSuccess': 'File "{filename}" downloaded successfully',
        'previewTitle': 'Preview: {filename}',
        'previewNotAvailable': 'Preview not available for this file type.',
        'unauthorized': 'Authorization required',
        'error': 'Error',
        'info': 'Info',
        'success': 'Success',
        'uploadAreaText': 'Drag files here or click to select',
        'selectFiles': 'Select Files',
        'loginTitle': 'Login',
        'loginSuccess': 'Login successful!',
        'logoutSuccess': 'You have been logged out',
        'fileTooLarge': 'File "{filename}" is too large. Maximum size: {maxSize}',
        'uploadError': 'Error uploading file "{filename}"',
        'downloadError': 'Error downloading file "{filename}"',
        'previewError': 'Preview error',
        'refreshError': 'Error refreshing file list',
        'viewToggle': 'View switched to: {view}',
        'listView': 'List',
        'gridView': 'Grid'
    }
};

// Класс для управления переводами
class Translator {
    constructor() {
        this.currentLanguage = localStorage.getItem('language') || 'ru';
        this.translations = translations;
    }

    // Установить язык
    setLanguage(lang) {
        if (this.translations[lang]) {
            this.currentLanguage = lang;
            localStorage.setItem('language', lang);
            this.applyTranslations();
        }
    }

    // Получить перевод
    get(key, params = {}) {
        let text = this.translations[this.currentLanguage][key] || key;
        
        // Заменить параметры в строке
        Object.keys(params).forEach(param => {
            text = text.replace(`{${param}}`, params[param]);
        });
        
        return text;
    }

    // Применить переводы ко всем элементам с атрибутом data-i18n
    applyTranslations() {
        const elements = document.querySelectorAll('[data-i18n]');
        elements.forEach(element => {
            const key = element.getAttribute('data-i18n');
            const params = {};
            
            // Получить параметры из data-i18n-params
            const paramsAttr = element.getAttribute('data-i18n-params');
            if (paramsAttr) {
                try {
                    Object.assign(params, JSON.parse(paramsAttr));
                } catch (e) {
                    console.error('Invalid JSON in data-i18n-params:', paramsAttr);
                }
            }
            
            element.textContent = this.get(key, params);
        });
        
        // Обновить title страницы
        document.title = this.get('title');
    }
}

// Создаем глобальный экземпляр переводчика
const translator = new Translator();