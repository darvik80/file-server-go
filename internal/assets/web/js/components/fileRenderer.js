// Компоненты для рендеринга файлов

// Создание HTML для файла
function createFileHTML(file, currentView, currentPath) {
    const icon = PathUtils.getFileIcon(file.name, file.isDir);
    const size = formatFileSize(file.size || 0);
    const canPreviewFile = canPreview(file.name);

    if (currentView === 'grid') {
        // Режим сетки
        return `
            <div class="file-item col-md-4 mb-3">
                <div class="card h-100">
                    <div class="card-body d-flex flex-column">
                        <div class="text-center mb-3">
                            <div class="file-icon" style="font-size: 2rem;">${icon}</div>
                        </div>
                        <h5 class="card-title text-truncate">${file.name}</h5>
                        <p class="card-text text-muted">${size}</p>
                        <div class="mt-auto">
                            <div class="file-actions d-flex justify-content-center gap-2">
                                ${file.name === '..' ? 
                                    `<button class="btn btn-secondary btn-sm" onclick="fileManager.navigateToFolder('${file.path}')">${translator.get('back')}</button>` :
                                    file.isDir ? 
                                        `<button class="btn btn-primary btn-sm" onclick="fileManager.navigateToFolder('${file.path}')">${translator.get('open')}</button>` : 
                                        `${canPreviewFile ? `<button class="btn btn-secondary btn-sm" onclick="fileManager.previewFile('${file.name}')">${translator.get('preview')}</button>` : ''}
                                        <button class="btn btn-primary btn-sm" onclick="fileManager.downloadFile('${file.name}')">${translator.get('download')}</button>`}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
    } else {
        // Режим списка (по умолчанию)
        return `
            <div class="file-item">
                <div class="file-icon">${icon}</div>
                <div class="file-info">
                    <div class="file-name">${file.name}</div>
                    <div class="file-meta">${size}</div>
                </div>
                <div class="file-actions">
                    ${file.name === '..' ? 
                        `<button class="btn btn-secondary" onclick="fileManager.navigateToFolder('${file.path}')">${translator.get('back')}</button>` :
                        file.isDir ? 
                            `<button class="btn btn-primary" onclick="fileManager.navigateToFolder('${file.path}')">${translator.get('open')}</button>` : 
                            `${canPreviewFile ? `<button class="btn btn-secondary" onclick="fileManager.previewFile('${file.name}')">${translator.get('preview')}</button>` : ''}
                            <button class="btn btn-primary" onclick="fileManager.downloadFile('${file.name}')">${translator.get('download')}</button>`}
                </div>
            </div>
        `;
    }
}

// Рендеринг списка файлов
function renderFiles(container, files, searchTerm, currentView) {
    if (files.length === 0) {
        container.innerHTML = searchTerm ?
            `<div class="text-center py-4">${translator.get('noFilesFound')}</div>` :
            `<div class="text-center py-4">${translator.get('noFiles')}</div>`;
        return;
    }

    // Устанавливаем класс контейнера в зависимости от режима отображения
    if (currentView === 'grid') {
        container.className = 'row g-3';
    } else {
        container.className = 'files-list';
    }

    const filesHTML = files.map(file => createFileHTML(file, currentView)).join('');
    container.innerHTML = filesHTML;

    // Применяем переводы к новым элементам
    translator.applyTranslations();
}

// Обновление навигационной цепочки (breadcrumb)
function updateBreadcrumb(breadcrumb, currentPath) {
    if (!breadcrumb) return;

    // Очищаем breadcrumb
    breadcrumb.innerHTML = '';

    // Если мы в корневой директории, показываем только "Файлы"
    if (currentPath === '.' || currentPath === '/') {
        breadcrumb.innerHTML = `<li class="breadcrumb-item active" aria-current="page">${translator.get('files')}</li>`;
        return;
    }

    // Разбиваем путь на части
    const pathParts = currentPath.split('/');

    // Добавляем ссылку на корневую директорию
    const rootItem = document.createElement('li');
    rootItem.className = 'breadcrumb-item';
    rootItem.innerHTML = `<a href="#" onclick="fileManager.navigateToFolder('.')">${translator.get('files')}</a>`;
    breadcrumb.appendChild(rootItem);

    // Добавляем промежуточные директории
    let pathSoFar = '';
    for (let i = 0; i < pathParts.length; i++) {
        if (pathParts[i] === '') continue;

        if (pathSoFar === '') {
            pathSoFar = pathParts[i];
        } else {
            pathSoFar += '/' + pathParts[i];
        }

        const item = document.createElement('li');
        item.className = 'breadcrumb-item';

        // Для последнего элемента делаем его активным (без ссылки)
        if (i === pathParts.length - 1) {
            item.className += ' active';
            item.setAttribute('aria-current', 'page');
            item.textContent = pathParts[i];
        } else {
            item.innerHTML = `<a href="#" onclick="fileManager.navigateToFolder('${pathSoFar}')">${pathParts[i]}</a>`;
        }

        breadcrumb.appendChild(item);
    }
}

// Делаем функции глобально доступными
window.createFileHTML = createFileHTML;
window.renderFiles = renderFiles;
window.updateBreadcrumb = updateBreadcrumb;