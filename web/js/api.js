import { handleFetchError, showError, showToast } from './utils.js';
import { PathUtils } from './constants.js';

export class FileManagerAPI {
    constructor(token) {
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
            showError(error.message);
            throw error;
        }
    }

    // Получить список файлов
    async getFiles(path = '.') {
        const normalizedPath = PathUtils.normalize(path);

        try {
            const response = await fetch(`/files?path=${encodeURIComponent(normalizedPath)}`, {
                method: 'GET',
                headers: { 'Authorization': `Bearer ${this.token}` }
            });

            return await handleFetchError(response);
        } catch (error) {
            throw error;
        }
    }

    // Загрузить файл
    async uploadFile(formData) {
        try {
            const response = await fetch('/upload', {
                method: 'POST',
                headers: { 'Authorization': `Bearer ${this.token}` },
                body: formData
            });

            return await handleFetchError(response);
        } catch (error) {
            throw error;
        }
    }

    // Скачать файл
    async downloadFile(filePath) {
        try {
            const response = await fetch(`/download/${filePath}`, {
                method: 'GET',
                headers: { 'Authorization': `Bearer ${this.token}` }
            });

            if (!response.ok) throw new Error('Ошибка скачивания файла');
            return await response.blob();
        } catch (error) {
            throw error;
        }
    }

    // Создать директорию
    async createDirectory(dirname, path) {
        try {
            const response = await fetch('/create-dir', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${this.token}`
                },
                body: JSON.stringify({ dirname, path })
            });

            return await handleFetchError(response);
        } catch (error) {
            throw error;
        }
    }

    // Удалить файл
    async deleteFile(path) {
        try {
            const response = await fetch('/delete-file?path=' + encodeURIComponent(path), {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${this.token}` }
            });

            return await handleFetchError(response);
        } catch (error) {
            throw error;
        }
    }

    // Удалить директорию
    async deleteDirectory(path) {
        try {
            const response = await fetch('/delete-dir?path=' + encodeURIComponent(path), {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${this.token}` }
            });

            return await handleFetchError(response);
        } catch (error) {
            throw error;
        }
    }
}