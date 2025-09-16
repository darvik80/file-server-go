// Класс для управления приложениями
class AppManager {
    constructor() {
        this.apps = [];
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

    // Обработка регистрации приложения
    async handleRegisterApp(e) {
        e.preventDefault();

        const name = document.getElementById('appName').value;
        const description = document.getElementById('appDescription').value;

        // Получаем выбранные разрешения
        const permissions = [];
        const readPermission = document.getElementById('readPermission');
        const writePermission = document.getElementById('writePermission');

        if (readPermission && readPermission.checked) {
            permissions.push('read');
        }
        if (writePermission && writePermission.checked) {
            permissions.push('write');
        }

        if (!name) {
            showNotification(translator.get('fillAppName'), 'error');
            return;
        }

        try {
            const response = await fetch('/register-app', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify({ name, description, permissions })
            });

            if (!response.ok) {
                if (response.status === 401) {
                    fileManager.logout();
                    return;
                }
                throw new Error('Ошибка при регистрации приложения');
            }

            const data = await response.json();
            showNotification(translator.get('appRegistered'), 'success');

            // Показываем данные приложения
            this.showAppCredentials(data);

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('registerAppModal'));
            if (modal) {
                modal.hide();
            }

            // Очищаем форму
            document.getElementById('registerAppForm').reset();
            // Сбрасываем чекбоксы разрешений
            if (readPermission) readPermission.checked = true;
            if (writePermission) writePermission.checked = false;

            // Обновляем список приложений
            this.refreshApplications();
        } catch (error) {
            handleError(error, translator.get('registerAppError'));
        }
    }

    // Показать учетные данные приложения
    showAppCredentials(app) {
        // Формируем строку разрешений
        let permissionsStr = '';
        if (app.permissions && Array.isArray(app.permissions)) {
            permissionsStr = app.permissions.map(perm => {
                switch (perm) {
                    case 'read':
                        return translator.get('readPermission');
                    case 'write':
                        return translator.get('writePermission');
                    default:
                        return perm;
                }
            }).join(', ');
        }

        const credentialsHTML = `
            <div class="alert alert-info mt-3">
                <h5>${translator.get('appCredentials')}</h5>
                <p><strong>${translator.get('appName')}:</strong> ${app.name}</p>
                <p><strong>${translator.get('accessKeyId')}:</strong> ${app.access_key_id}</p>
                <p><strong>${translator.get('accessKeySecret')}:</strong> ${app.access_key_secret}</p>
                <p><strong>${translator.get('permissions')}:</strong> ${permissionsStr}</p>
                <p class="text-muted">${translator.get('saveCredentials')}</p>
            </div>
        `;

        // Добавляем после формы регистрации
        const form = document.getElementById('registerAppForm');
        if (form) {
            form.insertAdjacentHTML('afterend', credentialsHTML);

            // Автоматически удаляем сообщение через 30 секунд
            setTimeout(() => {
                const alert = form.nextElementSibling;
                if (alert && alert.classList.contains('alert')) {
                    alert.remove();
                }
            }, 30000);
        }
    }

    // Обновить список приложений
    async refreshApplications() {
        try {
            showLoading(document.getElementById('appsList'));

            const response = await fetch('/applications', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (!response.ok) {
                if (response.status === 401) {
                    fileManager.logout();
                    return;
                }
                throw new Error('Ошибка при получении списка приложений');
            }

            const data = await response.json();
            this.apps = data.applications || [];
            this.renderApplications();
        } catch (error) {
            handleError(error, translator.get('refreshAppsError'));
            document.getElementById('appsList').innerHTML = `<div class="text-center py-4">${translator.get('error')}</div>`;
        }
    }

    // Рендеринг списка приложений
    renderApplications() {
        const container = document.getElementById('appsList');

        if (this.apps.length === 0) {
            container.innerHTML = `<div class="text-center py-4">${translator.get('noApps')}</div>`;
            return;
        }

        const appsHTML = this.apps.map(app => this.createAppHTML(app)).join('');
        container.innerHTML = `
            <div class="table-responsive">
                <table class="table table-striped table-hover">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>${translator.get('appName')}</th>
                            <th>${translator.get('accessKeyId')}</th>
                            <th>${translator.get('description')}</th>
                            <th>${translator.get('permissions')}</th>
                            <th>${translator.get('createdAt')}</th>
                            <th>${translator.get('actions')}</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${appsHTML}
                    </tbody>
                </table>
            </div>
        `;

        // Применяем переводы к новым элементам
        translator.applyTranslations();
    }

    // Создание HTML для приложения
    createAppHTML(app) {
        // Формируем строку разрешений
        let permissionsStr = '';
        if (app.permissions && Array.isArray(app.permissions)) {
            permissionsStr = app.permissions.map(perm => {
                switch (perm) {
                    case 'read':
                        return translator.get('readPermission');
                    case 'write':
                        return translator.get('writePermission');
                    default:
                        return perm;
                }
            }).join(', ');
        }

        return `
            <tr>
                <td>${app.id}</td>
                <td>${app.name}</td>
                <td>${app.access_key_id}</td>
                <td>${app.description || ''}</td>
                <td>${permissionsStr}</td>
                <td>${app.created_at ? new Date(app.created_at).toLocaleDateString() : ''}</td>
                <td>
                    <button class="btn btn-danger btn-sm" onclick="appManager.deleteApplication(${app.id})" data-i18n="delete">Удалить</button>
                </td>
            </tr>
        `;
    }

    // Удалить приложение
    async deleteApplication(id) {
        if (!confirm(translator.get('confirmDeleteApp'))) {
            return;
        }

        try {
            const response = await fetch(`/delete-app?id=${id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (!response.ok) {
                if (response.status === 401) {
                    fileManager.logout();
                    return;
                }
                throw new Error('Ошибка при удалении приложения');
            }

            showNotification(translator.get('appDeleted'), 'success');
            this.refreshApplications();
        } catch (error) {
            handleError(error, translator.get('deleteAppError'));
        }
    }
}

// Создаем глобальный экземпляр менеджера приложений
window.appManager = new AppManager();