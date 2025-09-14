import { FileManager } from './fileManager.js';

// Глобальная переменная для доступа из HTML
window.fileManager = null;

// Глобальные функции для HTML атрибутов
window.showContextMenu = function(e, item) {
    if (window.fileManager) {
        window.fileManager.showContextMenu(e, item);
    }
};

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    window.fileManager = new FileManager();
});