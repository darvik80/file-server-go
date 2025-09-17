// Главный файл приложения

// Инициализация приложения
document.addEventListener('DOMContentLoaded', function() {
    console.log('Файловый сервер запущен');
    
    // Создаем экземпляры менеджеров
    window.fileManager = new FileManager();
    window.userManager = new UserManager();
    window.appManager = new AppManager();
    
    // Обработчик закрытия модального окна
    const modal = document.getElementById('previewModal');
    if (modal) {
        // Закрытие по клику на фон
        modal.addEventListener('click', function(e) {
            if (e.target === modal) {
                fileManager.closePreview();
            }
        });
        
        // Закрытие по Escape
        document.addEventListener('keydown', function(e) {
            if (e.key === 'Escape' && modal.style.display === 'block') {
                fileManager.closePreview();
            }
        });
    }
    
    // Обработчики событий для вкладок
    const usersTab = document.getElementById('users-tab');
    const appsTab = document.getElementById('apps-tab');

    if (usersTab) {
        usersTab.addEventListener('shown.bs.tab', function() {
            if (window.userManager) {
                window.userManager.refreshUsers();
            }
        });
    }

    if (appsTab) {
        appsTab.addEventListener('shown.bs.tab', function() {
            if (window.appManager) {
                window.appManager.refreshApplications();
            }
        });
    }

    // Показываем приветственное сообщение
    setTimeout(() => {
        showNotification('Добро пожаловать в файловый сервер!', 'success');
    }, 500);
});

// Глобальные функции для использования в HTML
window.refreshFiles = function() {
    if (window.fileManager) {
        window.fileManager.refreshFiles();
    }
};

window.toggleView = function() {
    if (window.fileManager) {
        window.fileManager.toggleView();
    }
};

// Обработка ошибок JavaScript
window.addEventListener('error', function(e) {
    console.error('JavaScript Error:', e.error);
    handleError('Произошла ошибка в приложении', 'error');
});

// Обработка необработанных промисов
window.addEventListener('unhandledrejection', function(e) {
    console.error('Unhandled Promise Rejection:', e.reason);
    handleError('Произошла ошибка при выполнении операции', 'error');
    e.preventDefault();
});