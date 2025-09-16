// Компоненты для рендеринга пользователей

// Создание HTML для пользователя
function createUserHTML(user) {
    // Определяем отображаемое имя роли
    let roleDisplay = user.role;
    switch (user.role) {
        case 'admin':
            roleDisplay = translator.get('adminRole');
            break;
        case 'writer':
            roleDisplay = translator.get('writerRole');
            break;
        case 'reader':
            roleDisplay = translator.get('readerRole');
            break;
    }

    return `
        <tr>
            <td>${user.id}</td>
            <td>${user.username}</td>
            <td>${roleDisplay}</td>
            <td>${user.created_at ? new Date(user.created_at).toLocaleDateString() : ''}</td>
            <td>
                <button class="btn btn-danger btn-sm" onclick="userManager.deleteUser(${user.id})" data-i18n="delete">Удалить</button>
            </td>
        </tr>
    `;
}

// Рендеринг списка пользователей
function renderUsers(container, users) {
    if (users.length === 0) {
        container.innerHTML = `<div class="text-center py-4">${translator.get('noUsers')}</div>`;
        return;
    }

    const usersHTML = users.map(user => createUserHTML(user)).join('');
    container.innerHTML = `
        <div class="table-responsive">
            <table class="table table-striped table-hover">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>${translator.get('username')}</th>
                        <th>${translator.get('role')}</th>
                        <th>${translator.get('createdAt')}</th>
                        <th>${translator.get('actions')}</th>
                    </tr>
                </thead>
                <tbody>
                    ${usersHTML}
                </tbody>
            </table>
        </div>
    `;

    // Применяем переводы к новым элементам
    translator.applyTranslations();
}

// Делаем функции глобально доступными
window.createUserHTML = createUserHTML;
window.renderUsers = renderUsers;