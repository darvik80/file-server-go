// Менеджер приложений
// Убираем import, так как файлы подключаются через HTML

// import { renderApplications, showAppCredentials } from '../components/appRenderer.js';

class AppManager {
    constructor() {
        this.apps = [];
        this.token = localStorage.getItem('token');
        this.init();
    }

    init() {
        this.bindEvents();
    }

    bindEvents() {
        // Обработчик формы регистрации приложения
        const registerAppForm = document.getElementById('registerAppForm');
        if (registerAppForm) {
            registerAppForm.addEventListener('submit', (e) => this.handleRegisterApp(e));
        }
    }

    // Показать модальное окно регистрации приложения
    showRegisterAppModal() {
        const modal = new bootstrap.Modal(document.getElementById('registerAppModal'));
        modal.show();
    }

    // Обновление списка приложений
    async refreshApplications() {
        if (!this.token) return;

        try {
            showLoading(document.getElementById('appsList'));
            this.apps = await appAPI.getApplications();
            this.renderApplications();
        } catch (error) {
            if (error.message === 'unauthorized') {
                fileManager.logout();
                return;
            }
            handleError(error, translator.get('refreshAppsError'));
            document.getElementById('appsList').innerHTML = `<div class="text-center py-4">${translator.get('error')}</div>`;
        }
    }

    // Рендеринг списка приложений
    renderApplications() {
        const container = document.getElementById('appsList');
        // Используем глобальную функцию
        renderApplications(container, this.apps);
    }

    // Регистрация приложения
    async handleRegisterApp(e) {
        e.preventDefault();

        const name = document.getElementById('appName').value;
        const description = document.getElementById('appDescription').value;

        if (!name) {
            showNotification(translator.get('fillAppName'), 'error');
            return;
        }

        // Собираем разрешения
        const permissions = [];
        const readPermission = document.getElementById('readPermission');
        const writePermission = document.getElementById('writePermission');

        if (readPermission && readPermission.checked) {
            permissions.push('read');
        }
        if (writePermission && writePermission.checked) {
            permissions.push('write');
        }

        try {
            const app = await appAPI.registerApplication({
                name: name,
                description: description,
                permissions: permissions
            });

            showNotification(translator.get('appRegistered'), 'success');
            // Используем глобальную функцию
            showAppCredentials(app);

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('registerAppModal'));
            if (modal) {
                modal.hide();
            }

            // Очищаем форму
            document.getElementById('registerAppForm').reset();

            // Обновляем список приложений
            this.refreshApplications();
        } catch (error) {
            if (error.message === 'unauthorized') {
                fileManager.logout();
                return;
            }
            handleError(error, translator.get('registerAppError'));
        }
    }

    // Удаление приложения
    async deleteApplication(id) {
        if (!confirm(translator.get('confirmDeleteApp'))) return;

        try {
            await appAPI.deleteApplication(id);
            showNotification(translator.get('appDeleted'), 'success');
            this.refreshApplications();
        } catch (error) {
            if (error.message === 'unauthorized') {
                fileManager.logout();
                return;
            }
            handleError(error, translator.get('deleteAppError'));
        }
    }
}

// Делаем appManager глобально доступным
window.appManager = new AppManager();