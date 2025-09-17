// Менеджер пользователей
// Убираем import, так как файлы подключаются через HTML

// import { renderUsers } from '../components/userRenderer.js';

class UserManager {
    constructor() {
        this.users = [];
        this.token = localStorage.getItem('token');
        this.init();
    }

    init() {
        this.bindEvents();
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

    // Показать модальное окно создания пользователя
    showCreateUserModal() {
        const modal = new bootstrap.Modal(document.getElementById('createUserModal'));
        modal.show();
    }

    // Показать модальное окно изменения пароля
    showChangePasswordModal() {
        const modal = new bootstrap.Modal(document.getElementById('changePasswordModal'));
        modal.show();

        // Если пользователь администратор, скрываем поле старого пароля
        const userRole = fileManager.userRole;
        const oldPasswordGroup = document.getElementById('oldPasswordGroup');

        if (userRole === 'admin' && oldPasswordGroup) {
            oldPasswordGroup.style.display = 'none';
        } else if (oldPasswordGroup) {
            oldPasswordGroup.style.display = 'block';
        }

        // Очищаем поле имени пользователя
        const usernameField = document.getElementById('changePasswordUsername');
        if (usernameField) {
            usernameField.value = '';
        }
    }

    // Показать модальное окно изменения пароля для конкретного пользователя
    showChangePasswordModalForUser(username) {
        const modal = new bootstrap.Modal(document.getElementById('changePasswordModal'));
        modal.show();

        // Если пользователь администратор, скрываем поле старого пароля
        const userRole = fileManager.userRole;
        const oldPasswordGroup = document.getElementById('oldPasswordGroup');

        if (userRole === 'admin' && oldPasswordGroup) {
            oldPasswordGroup.style.display = 'none';
        } else if (oldPasswordGroup) {
            oldPasswordGroup.style.display = 'block';
        }

        // Предзаполняем поле имени пользователя
        const usernameField = document.getElementById('changePasswordUsername');
        if (usernameField) {
            usernameField.value = username;
            // Делаем поле только для чтения, если администратор редактирует пароль другого пользователя
            if (userRole === 'admin') {
                usernameField.readOnly = true;
            }
        }
    }

    // Обновление списка пользователей
    async refreshUsers() {
        if (!this.token) return;

        try {
            showLoading(document.getElementById('usersList'));
            this.users = await userAPI.getUsers();
            this.renderUsers();
        } catch (error) {
            if (error.message === 'unauthorized') {
                fileManager.logout();
                return;
            }
            handleError(error, translator.get('refreshUsersError'));
            document.getElementById('usersList').innerHTML = `<div class="text-center py-4">${translator.get('error')}</div>`;
        }
    }

    // Рендеринг списка пользователей
    renderUsers() {
        const container = document.getElementById('usersList');
        // Используем глобальную функцию
        renderUsers(container, this.users);
    }

    // Создание пользователя
    async handleCreateUser(e) {
        e.preventDefault();

        const username = document.getElementById('newUsername').value;
        const password = document.getElementById('newUserPassword').value;
        const role = document.getElementById('userRole').value;

        if (!username || !password || !role) {
            showNotification(translator.get('fillAllFields'), 'error');
            return;
        }

        try {
            await userAPI.createUser({
                username: username,
                password: password,
                role: role
            });

            showNotification(translator.get('userCreated'), 'success');

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
            if (error.message === 'unauthorized') {
                fileManager.logout();
                return;
            }
            handleError(error, translator.get('createUserError'));
        }
    }

    // Изменение пароля
    async handleChangePassword(e) {
        e.preventDefault();

        const username = document.getElementById('changePasswordUsername').value;
        const oldPassword = document.getElementById('oldPassword').value;
        const newPassword = document.getElementById('changePasswordNewPassword').value;

        if (!username || !newPassword) {
            showNotification(translator.get('fillAllFields'), 'error');
            return;
        }

        // Для администратора старый пароль не требуется
        const userRole = fileManager.userRole;
        if (userRole !== 'admin' && !oldPassword) {
            showNotification(translator.get('fillAllFields'), 'error');
            return;
        }

        try {
            const requestData = {
                username: username,
                new_password: newPassword
            };

            // Добавляем старый пароль, если это не администратор
            if (userRole !== 'admin') {
                requestData.old_password = oldPassword;
            }

            await userAPI.changePassword(requestData);

            showNotification(translator.get('passwordChanged'), 'success');

            // Закрываем модальное окно
            const modal = bootstrap.Modal.getInstance(document.getElementById('changePasswordModal'));
            if (modal) {
                modal.hide();
            }

            // Очищаем форму и сбрасываем состояние поля имени пользователя
            document.getElementById('changePasswordForm').reset();
            const usernameField = document.getElementById('changePasswordUsername');
            if (usernameField) {
                usernameField.readOnly = false; // Сбрасываем readonly состояние
            }
        } catch (error) {
            if (error.message === 'unauthorized') {
                fileManager.logout();
                return;
            }
            handleError(error, translator.get('changePasswordError'));
        }
    }

    // Удаление пользователя
    async deleteUser(id) {
        if (!confirm(translator.get('confirmDeleteUser'))) return;

        try {
            await userAPI.deleteUser(id);
            showNotification(translator.get('userDeleted'), 'success');
            this.refreshUsers();
        } catch (error) {
            if (error.message === 'unauthorized') {
                fileManager.logout();
                return;
            }
            handleError(error, translator.get('deleteUserError'));
        }
    }
}