// Компоненты для рендеринга приложений

// Создание HTML для приложения
function createAppHTML(app) {
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

// Рендеринг списка приложений
function renderApplications(container, apps) {
    if (apps.length === 0) {
        container.innerHTML = `<div class="text-center py-4">${translator.get('noApps')}</div>`;
        return;
    }

    const appsHTML = apps.map(app => createAppHTML(app)).join('');
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

// Показать учетные данные приложения
function showAppCredentials(app) {
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

// Делаем функции глобально доступными
window.createAppHTML = createAppHTML;
window.renderApplications = renderApplications;
window.showAppCredentials = showAppCredentials;