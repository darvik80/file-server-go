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
        'users': 'Пользователи',
        'applications': 'Приложения',
        'refresh': 'Обновить',
        'view': 'Вид',
        'search': 'Поиск файлов...',
        'sortBy': 'Сортировать по',
        'name': 'Имя',
        'size': 'Размер',
        'date': 'Дата',
        'download': 'Скачать',
        'preview': 'Просмотр',
        'delete': 'Удалить',
        'actions': 'Действия',
        'createdAt': 'Дата создания',
        'description': 'Описание',
        'appName': 'Название приложения',
        'appDescription': 'Описание приложения',
        'accessKeyId': 'Ключ доступа',
        'accessKeySecret': 'Секретный ключ',
        'addUser': 'Добавить пользователя',
        'registerApp': 'Зарегистрировать приложение',
        'noFiles': 'Нет загруженных файлов',
        'noUsers': 'Нет пользователей',
        'noApps': 'Нет зарегистрированных приложений',
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
        'gridView': 'Сетка',
        'fillAllFields': 'Заполните все поля',
        'fillAppName': 'Введите название приложения',
        'userCreated': 'Пользователь успешно создан',
        'appRegistered': 'Приложение успешно зарегистрировано',
        'userDeleted': 'Пользователь успешно удален',
        'appDeleted': 'Приложение успешно удалено',
        'createUserError': 'Ошибка при создании пользователя',
        'registerAppError': 'Ошибка при регистрации приложения',
        'deleteUserError': 'Ошибка при удалении пользователя',
        'deleteAppError': 'Ошибка при удалении приложения',
        'refreshUsersError': 'Ошибка при обновлении списка пользователей',
        'refreshAppsError': 'Ошибка при обновлении списка приложений',
        'confirmDeleteUser': 'Вы уверены, что хотите удалить этого пользователя?',
        'confirmDeleteApp': 'Вы уверены, что хотите удалить это приложение?',
        'appCredentials': 'Учетные данные приложения',
        'saveCredentials': 'Сохраните эти данные - они больше не будут показаны',
        'role': 'Роль',
        'permissions': 'Разрешения',
        'adminRole': 'Администратор',
        'writerRole': 'Писатель',
        'readerRole': 'Читатель',
        'readPermission': 'Чтение',
        'writePermission': 'Запись',
        'ru': 'Русский',
        'en': 'English',
        'uploadNotAllowed': 'Пользователям с ролью "Читатель" запрещено загружать файлы',
        // Добавленные переводы для функции создания директорий
        'createDir': 'Создать папку',
        'dirName': 'Имя папки',
        'dirCreated': 'Папка успешно создана',
        'createDirError': 'Ошибка при создании папки'
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
        'users': 'Users',
        'applications': 'Applications',
        'refresh': 'Refresh',
        'view': 'View',
        'search': 'Search files...',
        'sortBy': 'Sort by',
        'name': 'Name',
        'size': 'Size',
        'date': 'Date',
        'download': 'Download',
        'preview': 'Preview',
        'delete': 'Delete',
        'actions': 'Actions',
        'createdAt': 'Created At',
        'description': 'Description',
        'appName': 'Application Name',
        'appDescription': 'Application Description',
        'accessKeyId': 'Access Key ID',
        'accessKeySecret': 'Access Key Secret',
        'addUser': 'Add User',
        'registerApp': 'Register Application',
        'noFiles': 'No files uploaded',
        'noUsers': 'No users',
        'noApps': 'No applications registered',
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
        'gridView': 'Grid',
        'fillAllFields': 'Please fill in all fields',
        'fillAppName': 'Please enter application name',
        'userCreated': 'User created successfully',
        'appRegistered': 'Application registered successfully',
        'userDeleted': 'User deleted successfully',
        'appDeleted': 'Application deleted successfully',
        'createUserError': 'Error creating user',
        'registerAppError': 'Error registering application',
        'deleteUserError': 'Error deleting user',
        'deleteAppError': 'Error deleting application',
        'refreshUsersError': 'Error refreshing users list',
        'refreshAppsError': 'Error refreshing applications list',
        'confirmDeleteUser': 'Are you sure you want to delete this user?',
        'confirmDeleteApp': 'Are you sure you want to delete this application?',
        'appCredentials': 'Application Credentials',
        'saveCredentials': 'Save these credentials - they will not be shown again',
        'role': 'Role',
        'permissions': 'Permissions',
        'adminRole': 'Administrator',
        'writerRole': 'Writer',
        'readerRole': 'Reader',
        'readPermission': 'Read',
        'writePermission': 'Write',
        'ru': 'Russian',
        'en': 'English',
        'uploadNotAllowed': 'Users with "Reader" role are not allowed to upload files',
        // Added translations for directory creation function
        'createDir': 'Create Directory',
        'dirName': 'Directory Name',
        'dirCreated': 'Directory created successfully',
        'createDirError': 'Error creating directory'
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