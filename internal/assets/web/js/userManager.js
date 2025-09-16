// Класс для управления пользователями
class UserManager {
    constructor() {
        this.users = [];
        this.currentUserRole = null;
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadCurrentUserRole();
    }

    bindEvents() {
        // Обработчик формы создания пользователя
        const createUserForm = document.getElementById('createUserForm');
        if (createUserForm) {
            createUserForm.addEventListener('submit', (e) => this.handleCreateUser(e));
        }

        // Обработчик формы изменения пароля
        const changePasswordForm = document.getElementById('changePasswordForm');
        if (changePasswordForm) {
            changePasswordForm.addEventListener('submit', (e) => this.handleChangePassword(e));
        }
    }

    // Загрузка роли текущего пользователя
    async loadCurrentUserRole() {
        try {
            const response = await fetch('/user-info', {
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
                throw new Error('Ошибка при получении информации о пользователе');
            }

            const data = await response.json();
            this.currentUserRole = data.role;
        } catch (error) {
            console.error('Error loading user role:', error);
            this.currentUserRole = null;
        }
    }

    // Показать модальное окно создания пользователя
    showCreateUserModal() {
        const modal = new bootstrap.Modal(document.getElementById('createUserModal'));
        modal.show();
    }

    // Показать модальное окно изменения пароля
    showChangePasswordModal() {
        const modal = new bootstrap.Modal(document.getElementById('changePasswordModal'));
        modal.show();

        // Если пользователь - администратор, скрываем поле старого пароля
        const oldPasswordGroup = document.getElementById('oldPasswordGroup');
        if (oldPasswordGroup) {
            oldPasswordGroup.style.display = this.currentUserRole === 'admin' ? 'none' : 'block';
        }
    }

    // Обработка создания пользователя
    async handleCreateUser(e) {
        e.preventDefault();

        const username = document.getElementById('newUsername').value;
        const password = document.getElementById('newUserPassword').value;
        const role = document.getElementById('userRole').value;

        if (!username || !password) {
            showNotification(translator.get('fillAllFields'), 'error');
            return;
        }

        try {
            const response = await fetch('/create-user', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify({ username, password, role })
            });

            if (!response.ok) {
                if (response.status === 401) {
                    fileManager.logout();
                    return;
                }
                throw new Error('Ошибка при создании пользователя');
            }

            const data = await response.json();
            showNotification(data.message || translator.get('userCreated'), 'success');

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('createUserModal'));
            if (modal) {
                modal.hide();
            }

            // Очищаем форму
            document.getElementById('createUserForm').reset();

            // Обновляем список пользователей
            this.refreshUsers();
        } catch (error) {
            handleError(error, translator.get('createUserError'));
        }
    }

    // Обработка изменения пароля
    async handleChangePassword(e) {
        e.preventDefault();

        const username = document.getElementById('changePasswordUsername').value;
        const oldPassword = document.getElementById('oldPassword').value;
        const newPassword = document.getElementById('changePasswordNewPassword').value;

        if (!username || !newPassword) {
            showNotification(translator.get('fillAllFields'), 'error');
            return;
        }

        // Для администратора старый пароль не обязателен
        const requestData = {
            username: username,
            new_password: newPassword
        };

        // Для обычных пользователей требуется старый пароль
        if (this.currentUserRole !== 'admin') {
            if (!oldPassword) {
                showNotification(translator.get('fillAllFields'), 'error');
                return;
            }
            requestData.old_password = oldPassword;
        }

        try {
            const response = await fetch('/change-password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify(requestData)
            });

            if (!response.ok) {
                if (response.status === 401) {
                    fileManager.logout();
                    return;
                }
                throw new Error('Ошибка при изменении пароля');
            }

            const data = await response.json();
            showNotification(data.message || translator.get('passwordChanged'), 'success');

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('changePasswordModal'));
            if (modal) {
                modal.hide();
            }

            // Очищаем форму
            document.getElementById('changePasswordForm').reset();
        } catch (error) {
            handleError(error, translator.get('changePasswordError'));
        }
    }

    // Обновить список пользователей
    async refreshUsers() {
        try {
            showLoading(document.getElementById('usersList'));

            const response = await fetch('/users', {
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
                throw new Error('Ошибка при получении списка пользователей');
            }

            const data = await response.json();
            this.users = data.users || [];
            this.renderUsers();
        } catch (error) {
            handleError(error, translator.get('refreshUsersError'));
            document.getElementById('usersList').innerHTML = `<div class="text-center py-4">${translator.get('error')}</div>`;
        }
    }

    // Рендеринг списка пользователей
    renderUsers() {
        const container = document.getElementById('usersList');

        if (this.users.length === 0) {
            container.innerHTML = `<div class="text-center py-4">${translator.get('noUsers')}</div>`;
            return;
        }

        const usersHTML = this.users.map(user => this.createUserHTML(user)).join('');
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

    // Создание HTML для пользователя
    createUserHTML(user) {
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

    // Удалить пользователя
    async deleteUser(id) {
        if (!confirm(translator.get('confirmDeleteUser'))) {
            return;
        }

        try {
            const response = await fetch(`/delete-user?id=${id}`, {
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
                throw new Error('Ошибка при удалении пользователя');
            }

            showNotification(translator.get('userDeleted'), 'success');
            this.refreshUsers();
        } catch (error) {
            handleError(error, translator.get('deleteUserError'));
        }
    }
}

// Создаем глобальный экземпляр менеджера пользователей
window.userManager = new UserManager();