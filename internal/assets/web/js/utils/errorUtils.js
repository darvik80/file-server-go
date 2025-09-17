// Утилиты обработки ошибок

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

// Вспомогательная функция для экранирования HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

window.handleError = handleError;
window.debounce = debounce;
window.escapeHtml = escapeHtml;