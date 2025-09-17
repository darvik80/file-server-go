# File Server Go

[Russian Version](README.md)

Simple file server in Go with authentication support, user management and applications for managing files and directories.

## Features

- 📤 File upload via web interface or API
- 📥 File download with automatic content type detection
- 📋 View list of files and directories
- 📁 Create and delete directories
- 🔐 JWT authentication for users
- 🔑 API keys for applications (accessKeyId/accessKeySecret)
- 👥 User management system with roles
- 🎛️ Application management system
- 🌐 Modern web interface with localization support (RU/EN)
- 🔒 CORS support
- 📝 Request logging
- 🛡️ Security: path validation, file size limits

## User Roles

- **Admin** - full access to all functions, user and application management
- **Writer** - upload, create, delete files and directories
- **Reader** - view and download files only

## Quick Start

### Installation and Launch

```bash
# Clone repository
git clone <your-repo-url>
cd file-server-go

# Install dependencies
make deps

# Run server
make run
```

Server will be available at: http://localhost:8080

### Using Make Commands

```bash
make build      # Build application
make run        # Run application
make test       # Run tests
make fmt        # Format code
make vet        # Check code
make clean      # Clean build
make build-all  # Build for all platforms
make dev        # Run with hot reload (requires air)
```

## Configuration

Application is configured via environment variables:

- `PORT` - server port (default: 8080)
- `UPLOAD_DIR` - directory for uploaded files (default: ./uploads)
- `JWT_KEY` - secret key for JWT tokens (generated automatically)
- `ADMIN_USERNAME` - admin username (default: admin)
- `ADMIN_PASSWORD` - admin password (default: admin)

Example:
```bash
export PORT=3000
export UPLOAD_DIR=/path/to/uploads
export JWT_KEY=your-secret-key
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD=secure-password
make run
```

## API

### User Authentication

To get JWT token:
```bash
# Get token
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Response will contain token:
# {"token":"your-jwt-token"}
```

### Application Management (Admin only)

#### Register new application
```bash
# Create application with read and write permissions
curl -X POST http://localhost:8080/register-app \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Application",
    "description": "Application for file management",
    "permissions": ["read", "write"]
  }'

# Response contains credentials:
# {
#   "access_key_id": "ak_xxxxxxxxxxxxxxxx",
#   "access_key_secret": "sk_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
#   "name": "My Application",
#   "description": "Application for file management",
#   "permissions": ["read", "write"]
# }
```

#### Get list of applications
```bash
curl -X GET http://localhost:8080/applications \
  -H "Authorization: Bearer your-jwt-token"
```

#### Delete application
```bash
curl -X DELETE http://localhost:8080/delete-app?id=1 \
  -H "Authorization: Bearer your-jwt-token"
```

### File Upload via Application API

#### Upload file using accessKeyId and accessKeySecret
```bash
# Direct file upload (recommended method)
curl -X POST http://localhost:8080/upload/raw/ \
  -H "Authorization: App ak_your_access_key_id:sk_your_access_key_secret" \
  -H "Content-Type: application/octet-stream" \
  --data-binary @path/to/your/file.txt \
  --get --data-urlencode "filename=uploaded_file.txt"

# Upload with TTL (file will be deleted after 3600 seconds)
curl -X POST http://localhost:8080/upload/raw/ \
  -H "Authorization: App ak_your_access_key_id:sk_your_access_key_secret" \
  -H "Content-Type: application/octet-stream" \
  --data-binary @path/to/your/file.txt \
  --get --data-urlencode "filename=temp_file.txt" \
  --data-urlencode "ttl=3600"

# Upload via PUT method
curl -X PUT http://localhost:8080/upload/raw/my_document.pdf \
  -H "Authorization: App ak_your_access_key_id:sk_your_access_key_secret" \
  --data-binary @document.pdf
```

#### Examples in different programming languages

**Python:**
```python
import requests

# Application credentials
access_key_id = "ak_your_access_key_id"
access_key_secret = "sk_your_access_key_secret"

# Upload file
with open('file.txt', 'rb') as f:
    response = requests.post(
        'http://localhost:8080/upload/raw/',
        headers={
            'Authorization': f'App {access_key_id}:{access_key_secret}',
            'Content-Type': 'application/octet-stream'
        },
        params={'filename': 'uploaded_file.txt'},
        data=f
    )
    
print(response.json())
```

**JavaScript (Node.js):**
```javascript
const fs = require('fs');
const axios = require('axios');

const accessKeyId = 'ak_your_access_key_id';
const accessKeySecret = 'sk_your_access_key_secret';

const fileData = fs.readFileSync('file.txt');

axios.post('http://localhost:8080/upload/raw/', fileData, {
    headers: {
        'Authorization': `App ${accessKeyId}:${accessKeySecret}`,
        'Content-Type': 'application/octet-stream'
    },
    params: {
        filename: 'uploaded_file.txt'
    }
}).then(response => {
    console.log(response.data);
}).catch(error => {
    console.error(error.response.data);
});
```

**Go:**
```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "os"
)

func main() {
    accessKeyID := "ak_your_access_key_id"
    accessKeySecret := "sk_your_access_key_secret"
    
    file, err := os.Open("file.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    var buf bytes.Buffer
    io.Copy(&buf, file)
    
    req, err := http.NewRequest("POST", "http://localhost:8080/upload/raw/?filename=uploaded_file.txt", &buf)
    if err != nil {
        panic(err)
    }
    
    req.Header.Set("Authorization", fmt.Sprintf("App %s:%s", accessKeyID, accessKeySecret))
    req.Header.Set("Content-Type", "application/octet-stream")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    fmt.Println("Status:", resp.Status)
}
```

### Working with Files via JWT Token

#### Upload file via form
```bash
# Upload file to root directory
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer your-jwt-token" \
  -F "file=@path/to/your/file.txt"

# Upload file to specific directory
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer your-jwt-token" \
  -F "file=@path/to/your/file.txt" \
  -F "path=directory/subdirectory"
```

#### List files and directories
```bash
# Get list of files in root directory
curl http://localhost:8080/files \
  -H "Authorization: Bearer your-jwt-token"

# Get list of files in specific directory
curl -X POST http://localhost:8080/files \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{"path": "directory/subdirectory"}'
```

#### Download file
```bash
# Download file
curl -O -H "Authorization: Bearer your-jwt-token" \
  http://localhost:8080/download/path/to/file.txt
```

#### Working with directories
```bash
# Create directory
curl -X POST http://localhost:8080/create-dir \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "dirname": "new_directory",
    "path": "parent_directory"
  }'

# Delete directory
curl -X DELETE http://localhost:8080/delete-dir?path=directory/to/delete \
  -H "Authorization: Bearer your-jwt-token"

# Delete file
curl -X DELETE http://localhost:8080/delete-file?path=path/to/file.txt \
  -H "Authorization: Bearer your-jwt-token"
```

### User Management (Admin only)

#### Create user
```bash
curl -X POST http://localhost:8080/create-user \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "password": "password123",
    "role": "writer"
  }'
```

#### Get list of users
```bash
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer your-jwt-token"
```

#### Change password
```bash
# Admin can change password for any user
curl -X POST http://localhost:8080/change-password \
  -H "Authorization: Bearer admin-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "targetuser",
    "new_password": "newpassword123"
  }'

# Regular user can only change their own password
curl -X POST http://localhost:8080/change-password \
  -H "Authorization: Bearer user-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "currentuser",
    "old_password": "oldpassword",
    "new_password": "newpassword123"
  }'
```

#### Delete user
```bash
curl -X DELETE http://localhost:8080/delete-user?id=2 \
  -H "Authorization: Bearer your-jwt-token"
```

## Project Structure

```
file-server-go/
├── cmd/server/                    # Application entry point
├── internal/                     # Private application code
│   ├── assets/                   # Embedded resources (web interface)
│   ├── auth/                     # Authentication and JWT
│   ├── config/                   # Configuration
│   ├── database/                 # Database operations
│   ├── handlers/                 # HTTP handlers
│   │   ├── app_handler.go        # Application management
│   │   ├── auth_handler.go       # Authentication
│   │   ├── file_handler.go       # File operations
│   │   └── user_handler.go       # User management
│   ├── middleware/               # HTTP middleware
│   ├── models/                   # Data models
│   ├── server/                   # HTTP server
│   ├── services/                 # Business logic
│   └── utils/                    # Utility functions
├── uploads/                      # Uploaded files
├── Dockerfile                    # Docker configuration
├── Makefile                      # Make commands
└── README.md                     # Documentation
```

## Security

- All API requests require authentication (JWT token or API keys)
- Path validation to prevent path traversal attacks
- File size limit (default 32MB)
- MIME type validation for safe content display
- Sanitization of file and directory names
- Role-based access control system
- Secure storage of passwords and API keys
- Centralized security checks through utility functions

## Database

Project uses SQLite database with the following tables:
- `users` - system users
- `applications` - registered applications
- `files` - file metadata (for TTL functionality)

## Docker

```bash
# Build image
docker build -t file-server-go .

# Run container
docker run -p 8080:8080 -v $(pwd)/uploads:/app/uploads file-server-go
```

## Development

For development with automatic reload, install [Air](https://github.com/cosmtrek/air):

```bash
go install github.com/cosmtrek/air@latest
make dev
```

## License

MIT License