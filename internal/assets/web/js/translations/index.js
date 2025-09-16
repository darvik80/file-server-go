// Система переводов для поддержки нескольких языков
// Убираем import, так как файлы подключаются через HTML

// const translations = {
//     'ru': ruTranslations,
//     'en': enTranslations
// };

// Класс для управления переводами
class Translator {
    constructor() {
        this.currentLanguage = localStorage.getItem('language') || 'ru';
        // Используем глобальные объекты
        this.translations = {
            'ru': window.ruTranslations,
            'en': window.enTranslations
        };
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
window.translator = translator;