# File Server Go

[Russian Version](README.md)

Simple file server in Go with authentication support for managing files and directories.

## Features

- 📤 File upload via web interface or API
- 📥 File download
- 📋 View list of files and directories
- 📁 Create and delete directories
- 🔐 JWT authentication
- 🌐 Modern web interface
- 🔒 CORS support
- 📝 Request logging

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
```

## Configuration

Application is configured via environment variables:

- `PORT` - server port (default: 8080)
- `UPLOAD_DIR` - directory for uploaded files (default: ./uploads)
- `JWT_SECRET` - secret key for JWT tokens
- `AUTH_USERNAME` - username for authentication
- `AUTH_PASSWORD` - password for authentication

Example:
```bash
export PORT=3000
export UPLOAD_DIR=/path/to/uploads
export JWT_SECRET=your-secret-key
export AUTH_USERNAME=admin
export AUTH_PASSWORD=password
make run
```

## API

### Authentication

To get JWT token:
```bash
# Get token
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# Response will contain token:
# {"token":"your-jwt-token"}
```

### Working with Files

#### Upload file via curl
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

# Direct file upload (raw upload)
curl -X PUT http://localhost:8080/upload/raw/filename.txt \
  -H "Authorization: Bearer your-jwt-token" \
  --data-binary @path/to/your/file.txt
```

#### Upload file via wget
```bash
# Get token and save to variable
TOKEN=$(wget -qO- --post-data='{"username":"admin","password":"password"}' \
  --header='Content-Type: application/json' \
  http://localhost:8080/login | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Upload file
wget --method=POST \
  --header="Authorization: Bearer $TOKEN" \
  --body-file=/path/to/your/file.txt \
  http://localhost:8080/upload/raw/filename.txt
```

#### List files and directories
```bash
# Get list of files in root directory
curl http://localhost:8080/files \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json"

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
  -H "Authorization: Bearer your-jwt-token" \

# Delete file
curl -X DELETE http://localhost:8080/delete-file?path=path/to/file.txt\
  -H "Authorization: Bearer your-jwt-token" \
```

## Project Structure

```
file-server-go/
├── cmd/server/          # Application entry point
├── internal/           # Private application code
│   ├── auth/          # Authentication and JWT
│   ├── config/        # Configuration
│   ├── handlers/      # HTTP handlers
│   ├── middleware/    # HTTP middleware
│   └── server/        # HTTP server
├── web/              # Web interface
│   └── templates/    # HTML templates
├── uploads/          # Uploaded files
├── Makefile         # Make commands
└── README.md        # Documentation
```

## Security

- All API requests require JWT authentication
- Path validation to prevent path traversal attacks
- File size limit (default 32MB)
- MIME type validation
- Sanitization of file and directory names

## Development

For development with automatic reload, install [Air](https://github.com/cosmtrek/air):

```bash
go install github.com/cosmtrek/air@latest
make dev
```

## License

MIT License