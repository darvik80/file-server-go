// Утилиты уведомлений

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

// Показать загрузку
function showLoading(element, show = true) {
    if (show) {
        element.innerHTML = '<div class="loading">Загрузка...</div>';
    }
}

window.showNotification = showNotification;
window.showLoading = showLoading;