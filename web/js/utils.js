import { domCache, SELECTORS } from './constants.js';

// Показать ошибку
export function showError(message) {
    const errorSection = domCache[SELECTORS.ERROR_SECTION];
    const errorMessage = domCache[SELECTORS.ERROR_MESSAGE];

    if (errorSection && errorMessage) {
        errorMessage.textContent = message;
        errorSection.style.display = 'block';
    }
}

// Обработка ошибок запросов
export async function handleFetchError(response) {
    if (!response.ok) {
        if (response.status === 401) {
            throw new Error('unauthorized');
        }
        const data = await response.json();
        throw new Error(data.error || 'Произошла ошибка');
    }
    return response.json();
}

// Форматирование размера файла
export function formatFileSize(bytes) {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

// Показать уведомление
export function showToast(title, message) {
    const toastHTML = `
        <div class="toast align-items-center text-white bg-success border-0" role="alert">
            <div class="d-flex">
                <div class="toast-body">
                    <strong>${title}</strong><br>${message}
                </div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        </div>
    `;

    const toastContainer = document.createElement('div');
    toastContainer.className = 'toast-container position-fixed bottom-0 end-0 p-3';
    toastContainer.innerHTML = toastHTML;
    document.body.appendChild(toastContainer);

    const toast = new bootstrap.Toast(toastContainer.querySelector('.toast'));
    toast.show();

    toastContainer.querySelector('.toast').addEventListener('hidden.bs.toast', () => {
        toastContainer.remove();
    });
}