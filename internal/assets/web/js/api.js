// API для работы с файловым сервером с поддержкой авторизации

class FileServerAPI {
    constructor(token = null) {
        this.baseURL = '';
        this.token = token;
    }

    // Установить токен авторизации
    setToken(token) {
        this.token = token;
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

    // Получить список файлов
    async getFiles() {
        try {
            const headers = {};
            if (this.token) {
                headers['Authorization'] = `Bearer ${this.token}`;
            }

            const response = await fetch('/files', {
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
            console.error('Error fetching files:', error);
            throw error;
        }
    }

    // Загрузить файл
    async uploadFile(file, onProgress = null) {
        return new Promise((resolve, reject) => {
            const formData = new FormData();
            formData.append('file', file);

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

            const response = await fetch(`/download/${encodeURIComponent(filename)}`, {
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

// Создаем глобальный экземпляр API
window.fileAPI = new FileServerAPI();