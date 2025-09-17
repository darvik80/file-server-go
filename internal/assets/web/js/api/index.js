// Индексный файл для экспорта всех API классов

// Проверяем, что классы доступны в window
if (typeof window.FileServerAPI === 'undefined' ||
    typeof window.UserAPI === 'undefined' ||
    typeof window.AppAPI === 'undefined') {
    console.error('API classes are not available in window object');
}

// Создаем глобальные экземпляры API с учетом токена из localStorage
if (typeof window !== 'undefined') {
    window.fileAPI = new window.FileServerAPI();
    window.userAPI = new window.UserAPI();
    window.appAPI = new window.AppAPI();
}