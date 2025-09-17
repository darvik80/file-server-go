// API для работы с приложениями
class AppAPI {
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

    // Получить список приложений
    async getApplications() {
        try {
            const response = await fetch('/applications', {
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
            // Возвращаем массив приложений, а не весь объект
            return data.applications || [];
        } catch (error) {
            console.error('Error fetching applications:', error);
            throw error;
        }
    }

    // Зарегистрировать приложение
    async registerApplication(appData) {
        try {
            const response = await fetch('/applications', {
                method: 'POST',
                headers: this.getAuthHeaders(),
                body: JSON.stringify(appData)
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
            console.error('Error registering application:', error);
            throw error;
        }
    }

    // Удалить приложение
    async deleteApplication(id) {
        try {
            const response = await fetch(`/applications/${id}`, {
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
            console.error('Error deleting application:', error);
            throw error;
        }
    }
}

// Экспортируем класс для использования в других модулях
if (typeof module !== 'undefined' && module.exports) {
    module.exports = AppAPI;
} else if (typeof window !== 'undefined') {
    window.AppAPI = AppAPI;
}