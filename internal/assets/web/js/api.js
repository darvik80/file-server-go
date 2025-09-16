// API для работы с файловым сервером с поддержкой авторизации

class FileServerAPI {
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
        // Обновляем токен во всех API экземплярах
        this.updateAllAPITokens(token);
    }

    // Обновить токен во всех API экземплярах
    updateAllAPITokens(token) {
        if (window.userAPI) {
            window.userAPI.setToken(token);
        }
        if (window.appAPI) {
            window.appAPI.setToken(token);
        }

        // Также обновляем токен в менеджерах
        if (window.userManager) {
            window.userManager.token = token;
        }
        if (window.appManager) {
            window.appManager.token = token;
        }
    }

    // Авторизация
    async login(username, password) {
        try {
            const response = await fetch('/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });

            if (!response.ok) {
                const data = await response.json();
                throw new Error(data.error || 'Ошибка авторизации');
            }

            const data = await response.json();
            // Устанавливаем токен после успешной авторизации
            if (data.token) {
                this.setToken(data.token);
            }
            return data;
        } catch (error) {
            console.error('Login error:', error);
            throw error;
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

    // Получить информацию о текущем пользователе
    async getUserInfo() {
        try {
            const headers = {};
            if (this.token) {
                headers['Authorization'] = `Bearer ${this.token}`;
            }

            const response = await fetch('/user-info', {
                method: 'GET',
                headers: headers
            });

            if (!response.ok) {
                if (response.status === 401) {
                    throw new Error('unauthorized');
                }
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error('Error fetching user info:', error);
            throw error;
        }
    }

    // Получить список файлов
    async getFiles(path = null) {
        try {
            let url = '/files';
            const headers = {};

            if (this.token) {
                headers['Authorization'] = `Bearer ${this.token}`;
            }

            // Если путь задан, отправляем его в теле POST запроса
            let response;
            if (path) {
                response = await fetch(url, {
                    method: 'POST',
                    headers: {
                        ...headers,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ path: path })
                });
            } else {
                // Если путь не задан, используем GET запрос
                response = await fetch(url, {
                    method: 'GET',
                    headers: headers
                });
            }

            if (!response.ok) {
                if (response.status === 401) {
                    throw new Error('unauthorized');
                }
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
        } catch (error) {
            console.error('Error fetching files:', error);
            throw error;
        }
    }

    // Загрузить файл
    async uploadFile(file, path = null, onProgress = null) {
        return new Promise((resolve, reject) => {
            const formData = new FormData();
            formData.append('file', file);

            // Добавляем путь, если он задан
            if (path) {
                formData.append('path', path);
            }

            // Отладочный вывод для проверки передачи параметра path
            console.log('Uploading file with path:', path);
            for (let pair of formData.entries()) {
                console.log(pair[0] + ': ' + pair[1]);
            }

            const xhr = new XMLHttpRequest();

            // Отслеживание прогресса
            if (onProgress) {
                xhr.upload.addEventListener('progress', (e) => {
                    if (e.lengthComputable) {
                        const percentComplete = (e.loaded / e.total) * 100;
                        onProgress(percentComplete);
                    }
                });
            }

            xhr.addEventListener('load', () => {
                if (xhr.status === 200) {
                    try {
                        const response = JSON.parse(xhr.responseText);
                        resolve(response);
                    } catch (e) {
                        resolve({ message: 'Файл успешно загружен' });
                    }
                } else if (xhr.status === 401) {
                    reject(new Error('unauthorized'));
                } else {
                    reject(new Error(`Ошибка загрузки: ${xhr.status}`));
                }
            });

            xhr.addEventListener('error', () => {
                reject(new Error('Ошибка сети при загрузке файла'));
            });

            xhr.open('POST', '/upload');

            // Добавляем токен авторизации
            if (this.token) {
                xhr.setRequestHeader('Authorization', `Bearer ${this.token}`);
            }

            xhr.send(formData);
        });
    }

    // Скачать файл
    async downloadFile(filename) {
        try {
            const headers = {};
            if (this.token) {
                headers['Authorization'] = `Bearer ${this.token}`;
            }

            // Убедимся, что путь к файлу правильно закодирован
            const encodedFilename = filename.split('/').map(part => encodeURIComponent(part)).join('/');

            const response = await fetch(`/download/${encodedFilename}`, {
                method: 'GET',
                headers: headers
            });

            if (!response.ok) {
                if (response.status === 401) {
                    throw new Error('unauthorized');
                }
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.blob();
        } catch (error) {
            console.error('Error downloading file:', error);
            throw error;
        }
    }

    // Получить информацию о файле (для предварительного просмотра)
    async getFileInfo(filename) {
        try {
            // Для простоты возвращаем базовую информацию
            return {
                name: filename,
                canPreview: canPreview(filename)
            };
        } catch (error) {
            console.error('Error getting file info:', error);
            throw new Error('Ошибка при получении информации о файле');
        }
    }
}

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

// Создаем глобальные экземпляры API с учетом токена из localStorage
window.fileAPI = new FileServerAPI();
window.userAPI = new UserAPI();
window.appAPI = new AppAPI();