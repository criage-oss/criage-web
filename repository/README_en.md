# Criage Repository Server

Repository server for storing and managing Criage packages.

## Features

- 📦 **Package Storage** - automatic indexing of uploaded packages
- 🔍 **Package Search** - fast search by name, description, author
- 📊 **Statistics** - tracking downloads and package popularity
- 🌐 **REST API** - full-featured API for integrations
- 🚀 **Web Interface** - simple web interface for browsing packages
- 🔒 **Security** - authentication for package uploads

## Installation and Setup

### Requirements

- Go 1.21 or higher
- Git

### Build and Run

```bash
# Clone repository
git clone <repository-url>
cd criage/repository

# Install dependencies
go mod tidy

# Build server
go build -o criage-repository

# Run with default configuration
./criage-repository

# Or specify configuration file
./criage-repository -config /path/to/config.json
```

Server will be available at: `http://localhost:8080`

## Configuration

On first run, a configuration file `config.json` is created:

```json
{
  "port": 8080,
  "storage_path": "./packages",
  "index_path": "./index.json",
  "upload_token": "your-secret-token",
  "max_file_size": 104857600,
  "allowed_formats": [
    "criage",
    "tar.zst",
    "tar.lz4", 
    "tar.xz",
    "tar.gz",
    "zip"
  ],
  "enable_cors": true,
  "log_level": "info"
}
```

### Configuration Parameters

- `port` - HTTP server port (default 8080)
- `storage_path` - path to packages directory (./packages)
- `index_path` - path to index file (./index.json)
- `upload_token` - token for package uploads
- `max_file_size` - maximum upload file size in bytes
- `allowed_formats` - allowed archive formats
- `enable_cors` - enable CORS headers
- `log_level` - logging level

## API Endpoints

### Repository Information

```
GET /api/v1/
```

Returns repository information, package count, and supported formats.

### Package List

```
GET /api/v1/packages?page=1&limit=20
```

Returns list of all packages with pagination.

### Package Information

```
GET /api/v1/packages/{name}
```

Returns detailed package information with all versions.

### Version Information

```
GET /api/v1/packages/{name}/{version}
```

Returns information about specific package version.

### Package Search

```
GET /api/v1/search?q={query}&limit=20
```

Searches packages by name, description, keywords.

### Package Download

```
GET /api/v1/download/{name}/{version}/{filename}
```

Downloads package file. Automatically increments download counter.

### Package Upload

```
POST /api/v1/upload
Headers: Authorization: Bearer {token}
Content-Type: multipart/form-data
```

Uploads new package. Requires authorization token.

### Statistics

```
GET /api/v1/stats
```

Returns repository statistics: popular packages, download counts, breakdown by licenses and authors.

### Index Refresh

```
POST /api/v1/refresh
Headers: Authorization: Bearer {token}
```

Forces package index update.

## Usage

### Upload package via curl

```bash
curl -X POST \
  -H "Authorization: Bearer your-secret-token" \
  -F "package=@test-package-1.0.0.criage" \
  http://localhost:8080/api/v1/upload
```

### Search packages

```bash
curl "http://localhost:8080/api/v1/search?q=test"
```

### Download package

```bash
curl -O "http://localhost:8080/api/v1/download/test-package/1.0.0/test-package-1.0.0.criage"
```

## Integration with criage

To publish packages to repository use command:

```bash
criage publish --registry http://localhost:8080 --token your-secret-token
```

To install packages from repository:

```bash
criage repo add myrepo http://localhost:8080
criage install package-name --repo myrepo
```

## Data Structure

### Repository Index

Server automatically maintains JSON index of all packages in `index.json` file:

```json
{
  "last_updated": "2024-01-15T10:30:45Z",
  "total_packages": 25,
  "packages": {
    "package-name": {
      "name": "package-name",
      "description": "Package description",
      "author": "Author Name",
      "license": "MIT",
      "versions": [
        {
          "version": "1.0.0",
          "files": [
            {
              "os": "linux",
              "arch": "amd64",
              "format": "criage",
              "filename": "package-1.0.0.criage",
              "size": 1024,
              "checksum": "sha256:..."
            }
          ]
        }
      ]
    }
  },
  "statistics": {
    "total_downloads": 1250,
    "packages_by_license": {
      "MIT": 15,
      "Apache-2.0": 8,
      "GPL-3.0": 2
    },
    "packages_by_author": {
      "John Doe": 5,
      "Jane Smith": 3
    },
    "popular_packages": [
      "web-framework",
      "database-driver",
      "logging-lib"
    ]
  }
}
```

### Metadata Extraction

Server automatically extracts metadata from uploaded criage packages:

- Package manifest (`criage.yaml`)
- Build manifest (`build.json`)
- Compression information
- Dependencies and descriptions

## Deployment

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o criage-repository

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/criage-repository .
COPY --from=builder /app/web ./web/

EXPOSE 8080
CMD ["./criage-repository"]
```

### Systemd Service

```ini
[Unit]
Description=Criage Repository Server
After=network.target

[Service]
Type=simple
User=criage
WorkingDirectory=/opt/criage-repository
ExecStart=/opt/criage-repository/criage-repository -config /etc/criage/config.json
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name packages.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Handle large package uploads
        client_max_body_size 100M;
    }
}
```

## Web Interface

The repository includes a built-in web interface accessible at the server root URL. Features include:

- **Package Browser** - browse all available packages
- **Search Interface** - search packages with filters
- **Statistics Dashboard** - view download stats and popular packages
- **Package Details** - detailed package information with download links
- **Responsive Design** - works on desktop and mobile devices

### Web Interface Features

- Real-time search with autocomplete
- Package filtering by license, author, format
- Download statistics visualization
- Package dependency graphs
- Mobile-responsive design

## Monitoring

Server logs all HTTP requests and package operations. Logs include:

- HTTP requests with execution time
- New package uploads
- Indexing errors
- Download statistics
- Authentication attempts

### Log Format

```
2024-01-15T10:30:45Z [INFO] HTTP 200 GET /api/v1/packages?page=1 - 45ms
2024-01-15T10:31:12Z [INFO] Package uploaded: test-package-1.0.0.criage (1024 bytes)
2024-01-15T10:31:30Z [INFO] Download: test-package/1.0.0/test-package-1.0.0.criage
2024-01-15T10:32:01Z [ERROR] Failed to extract metadata from package: invalid format
```

## Security

- Authentication via Bearer tokens
- Upload file size limits
- File format validation
- CORS support for web interfaces
- Rate limiting (configurable)
- Input sanitization and validation

### Security Best Practices

1. **Use strong tokens** - Generate cryptographically secure upload tokens
2. **HTTPS deployment** - Always use HTTPS in production
3. **File validation** - Server validates package formats and checksums
4. **Access logs** - Monitor access patterns for suspicious activity
5. **Regular updates** - Keep dependencies updated

## Performance

- Asynchronous index updates
- Package metadata caching
- Efficient index search
- HTTP Keep-Alive support
- Concurrent request handling
- Background cleanup tasks

### Performance Tuning

- **Index caching** - In-memory index for fast searches
- **Parallel uploads** - Handle multiple uploads simultaneously
- **Compression** - Serve compressed responses when supported
- **Static file caching** - Cache web interface assets

## Development

### Building from Source

```bash
git clone https://github.com/Zu-Krein/criage.git
cd criage/repository

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build
go build -o criage-repository

# Run in development mode
go run . -config config.json
```

### API Testing

```bash
# Test repository info
curl http://localhost:8080/api/v1/

# Test package search
curl "http://localhost:8080/api/v1/search?q=test"

# Test upload (requires token)
curl -X POST \
  -H "Authorization: Bearer test-token" \
  -F "package=@test-package.criage" \
  http://localhost:8080/api/v1/upload
```

## Troubleshooting

### Common Issues

1. **Port already in use** - Check if another service is using the configured port
2. **Permission denied** - Ensure write permissions for packages directory
3. **Upload failures** - Check file size limits and allowed formats
4. **Index corruption** - Delete index.json to force regeneration

### Debug Mode

Run with debug logging:

```bash
./criage-repository -config config.json -log-level debug
```

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## Support

- 📧 Email: <support@criage.ru>
- 🐛 Issues: <https://github.com/Zu-Krein/criage/issues>
- 📖 Documentation: <https://docs.criage.ru>

## License

MIT License - see LICENSE file for details.
