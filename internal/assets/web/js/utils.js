// Утилиты приложения

// Показать уведомление
function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.innerHTML = `
        <div class="notification-content">
            <span>${message}</span>
            <button class="notification-close" onclick="this.parentElement.parentElement.remove()">&times;</button>
        </div>
    `;
    
    const container = document.getElementById('notifications');
    container.appendChild(notification);
    
    // Автоматическое удаление через 5 секунд
    setTimeout(() => {
        if (notification.parentElement) {
            notification.remove();
        }
    }, 5000);
}

// Форматирование размера файла
function formatFileSize(bytes) {
    if (bytes === 0) return '0 B';
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
    return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`;
}

// Форматирование даты
function formatDate(dateString) {
    if (!dateString) return 'Неизвестно';
    
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) {
        return 'Сегодня ' + date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
    } else if (diffDays === 1) {
        return 'Вчера ' + date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
    } else if (diffDays < 7) {
        return `${diffDays} дн. назад`;
    } else {
        return date.toLocaleDateString('ru-RU');
    }
}

// Обработка ошибок
function handleError(error, context = '') {
    console.error('Error:', error);
    let message = 'Произошла ошибка';
    
    if (error.message) {
        message = error.message;
    } else if (typeof error === 'string') {
        message = error;
    }
    
    if (context) {
        message = `${context}: ${message}`;
    }
    
    showNotification(message, 'error');
}

// Показать загрузку
function showLoading(element, show = true) {
    if (show) {
        element.innerHTML = '<div class="loading">Загрузка...</div>';
    }
}

// Debounce функция
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Проверка типа файла для предварительного просмотра
// Теперь сервер правильно устанавливает Content-Type, поэтому разрешаем предварительный просмотр для всех файлов
function canPreview(filename) {
    // Всегда разрешаем предварительный просмотр, так как сервер теперь правильно устанавливает Content-Type
    return true;
}

// Экспорт функций в глобальную область
window.showNotification = showNotification;
window.formatFileSize = formatFileSize;
window.formatDate = formatDate;
window.handleError = handleError;
window.showLoading = showLoading;
window.debounce = debounce;
window.canPreview = canPreview;