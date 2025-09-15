// Константы приложения
const CONSTANTS = {
    API_BASE: '',
    MAX_FILE_SIZE: 10 * 1024 * 1024, // 10MB
    SUPPORTED_PREVIEW_TYPES: ['image/jpeg', 'image/png', 'image/gif', 'text/plain', 'application/json']
};

// Утилиты для работы с путями
const PathUtils = {
    normalize(path) {
        return decodeURIComponent(path).replace(/\\/g, '/');
    },
    getFileName(path) {
        return path.split('/').pop();
    },
    getParentPath(path) {
        const parts = path.split('/');
        parts.pop();
        return parts.join('/') || '.';
    },
    encodePath(path) {
        return encodeURIComponent(this.normalize(path));
    },
    getFileExtension(filename) {
        return filename.split('.').pop().toLowerCase();
    },
    getFileIcon(filename, isDir) {
        if (isDir) return '📁';
        
        const ext = this.getFileExtension(filename);
        const iconMap = {
            'pdf': '📄',
            'doc': '📝', 'docx': '📝',
            'xls': '📊', 'xlsx': '📊',
            'ppt': '📈', 'pptx': '📈',
            'zip': '📦', 'rar': '📦', '7z': '📦',
            'jpg': '🖼️', 'jpeg': '🖼️', 'png': '🖼️', 'gif': '🖼️',
            'mp4': '🎥', 'avi': '🎥', 'mov': '🎥',
            'mp3': '🎵', 'wav': '🎵', 'flac': '🎵',
            'txt': '📄', 'md': '📄',
            'js': '📜', 'html': '📜', 'css': '📜', 'json': '📜'
        };
        
        return iconMap[ext] || '📄';
    }
};

// Экспорт для глобального использования
window.CONSTANTS = CONSTANTS;
window.PathUtils = PathUtils;