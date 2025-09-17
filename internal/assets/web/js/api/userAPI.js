// API для работы с пользователями
class UserAPI {
    constructor(token = null) {
        this.baseURL = '';
        this.token = token || localStorage.getItem('token'); // Получаем токен из localStorage если не передан
    }

    // Установить токен авторизации
    setToken(token) {
        this.token = token;
        // Сохраняем токен в localStorage
        if (token) {
            localStorage.setItem('token', token);
        } else {
            localStorage.removeItem('token');
        }
    }

    // Получить заголовки с авторизацией
    getAuthHeaders() {
        const headers = {
            'Content-Type': 'application/json'
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        return headers;
    }

    // Получить список пользователей
    async getUsers() {
        try {
            const response = await fetch('/users', {
                method: 'GET',
                headers: this.getAuthHeaders()
            });

            if (!response.ok) {
                if (response.status === 401) {
                    throw new Error('unauthorized');
                }
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            // Получаем данные из ответа
            const data = await response.json();
            // Возвращаем массив пользователей, а не весь объект
            return data.users || [];
        } catch (error) {
            console.error('Error fetching users:', error);
            throw error;
        }
    }

    // Создать пользователя
    async createUser(userData) {
        try {
            const response = await fetch('/users', {
                method: 'POST',
                headers: this.getAuthHeaders(),
                body: JSON.stringify(userData)
            });

            if (!response.ok) {
                if (response.status === 401) {
                    throw new Error('unauthorized');
                }
                const data = await response.json();
                throw new Error(data.error || `HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error('Error creating user:', error);
            throw error;
        }
    }

    // Изменить пароль пользователя
    async changePassword(passwordData) {
        try {
            const response = await fetch('/change-password', {
                method: 'POST',
                headers: this.getAuthHeaders(),
                body: JSON.stringify(passwordData)
            });

            if (!response.ok) {
                if (response.status === 401) {
                    throw new Error('unauthorized');
                }
                const data = await response.json();
                throw new Error(data.error || `HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error('Error changing password:', error);
            throw error;
        }
    }

    // Удалить пользователя
    async deleteUser(id) {
        try {
            const response = await fetch(`/users/${id}`, {
                method: 'DELETE',
                headers: this.getAuthHeaders()
            });

            if (!response.ok) {
                if (response.status === 401) {
                    throw new Error('unauthorized');
                }
                const data = await response.json();
                throw new Error(data.error || `HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error('Error deleting user:', error);
            throw error;
        }
    }
}

// Экспортируем класс для использования в других модулях
if (typeof module !== 'undefined' && module.exports) {
    module.exports = UserAPI;
} else if (typeof window !== 'undefined') {
    window.UserAPI = UserAPI;
}