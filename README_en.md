<div align="center">
  <img src="logo.png" alt="Criage Logo" width="200">
  
# Criage - High-Performance Package Manager
  
  Criage is a modern package manager written in Go, providing fast installation, updates, and package management with support for various compression formats.
  
  🇬🇧 English Version | [🇷🇺 Русская версия](README.md)
</div>

## Features

### Core Functions

- 🚀 **High Performance** - uses fast compression algorithms (Zstandard, LZ4)
- 📦 **Unified Package Extension** - all packages use `.criage` extension with embedded metadata about compression type
- 🔧 **Dependency Management** - automatic dependency resolution and installation
- 🌐 **Multiple Repositories** - support for multiple package sources
- 🎯 **Cross-Platform** - supports Linux, macOS, Windows
- ⚡ **Parallel Operations** - multithreaded processing for acceleration

### Package Management

- Install and remove packages
- Update to latest versions
- Search packages in repositories
- View package information
- Global and local installation

### Package Development

- Create new packages from templates
- Build packages with customizable scripts
- Publish to repositories
- Lifecycle hooks (pre/post install/remove)
- Build manifests

## Installation

### From Source

```bash
git clone https://github.com/Zu-Krein/criage.git
cd criage
go build -o criage
sudo mv criage /usr/local/bin/
```

### Verify Installation

```bash
criage --version
```

## Usage

### Basic Commands

#### Installing Packages

```bash
# Install package
criage install package-name

# Install specific version
criage install package-name --version 1.2.3

# Global installation
criage install package-name --global

# Install with dev dependencies
criage install package-name --dev
```

#### Removing Packages

```bash
# Remove package
criage uninstall package-name

# Complete removal with configuration
criage uninstall package-name --purge
```

#### Updating Packages

```bash
# Update specific package
criage update package-name

# Update all packages
criage update --all
```

#### Search and Information

```bash
# Find packages
criage search keyword

# Show installed packages
criage list

# Show only outdated packages
criage list --outdated

# Detailed package information
criage info package-name
```

### Package Development

#### Creating New Package

```bash
# Create package from basic template
criage create my-package --author "Your Name" --description "Package description"
```

#### Building Package

```bash
# Build with default settings (creates .criage file)
criage build

# Specify compression type and level
criage build --format tar.zst --compression 6 --output my-package-1.0.0.criage
```

#### Publishing Package

```bash
# Publish to repository
criage publish --registry https://packages.example.com --token YOUR_TOKEN
```

### Configuration

#### View Settings

```bash
# Show all settings
criage config list

# Get specific setting value
criage config get cache_path
```

#### Change Settings

```bash
# Change cache path
criage config set cache_path /custom/cache/path

# Change default compression level
criage config set compression.level 6

# Change number of parallel threads
criage config set parallel 8
```

## Project Structure

```
criage/
├── main.go              # Main entry point
├── commands.go          # CLI command implementations
├── go.mod               # Go module
├── go.sum               # Dependencies
└── pkg/                 # Core packages
    ├── types.go         # Data structures
    ├── archive.go       # Archive operations
    ├── config.go        # Configuration management
    ├── package_manager.go        # Main package manager logic
    └── package_manager_helpers.go # Helper functions
```

## File Formats

### Package Manifest (criage.yaml)

```yaml
name: my-package
version: 1.0.0
description: Example package
author: Your Name
license: MIT
homepage: https://github.com/user/my-package
repository: https://github.com/user/my-package

keywords:
  - utility
  - tool

dependencies:
  some-lib: ^1.0.0

dev_dependencies:
  test-framework: ^2.0.0

scripts:
  build: make build
  test: make test
  install: make install

files:
  - "bin/*"
  - "lib/*"
  - "README.md"

exclude:
  - "*.log"
  - ".git"
  - "node_modules"

arch:
  - amd64
  - arm64

os:
  - linux  
  - darwin
  - windows

hooks:
  pre_install:
    - echo "Installing package..."
  post_install:
    - echo "Package installed successfully"
```

### Build Configuration (build.json)

```json
{
  "name": "my-package",
  "version": "1.0.0",
  "build_script": "make build",
  "build_env": {
    "CGO_ENABLED": "0",
    "GOOS": "linux"
  },
  "output_dir": "./dist",
  "include_files": ["bin/*", "lib/*"],
  "exclude_files": ["*.log", "test/*"],
  "compression": {
    "format": "tar.zst",
    "level": 3
  },
  "targets": [
    {"os": "linux", "arch": "amd64"},
    {"os": "linux", "arch": "arm64"},
    {"os": "darwin", "arch": "amd64"},
    {"os": "windows", "arch": "amd64"}
  ]
}
```

## Embedding Metadata in Archives

Criage supports embedding package metadata (`criage.yaml` and `build.json`) directly into archives. This allows getting package information without needing to extract it.

### Supported Formats

#### TAR Archives (tar.zst, tar.lz4, tar.xz, tar.gz)

- Uses **PAX Extended Headers** - standard mechanism for storing additional metadata
- Compatible with most modern archivers
- Metadata stored in fields `criage.metadata`, `criage.package_manifest`, `criage.build_manifest`

#### ZIP Archives

- Uses **ZIP Comment** for basic metadata
- Additionally creates `.criage_metadata.json` file inside archive
- Full backward compatibility

### Embedded Data

- **Package Manifest** (`criage.yaml`) - name, version, dependencies, author
- **Build Manifest** (`build.json`) - build settings, target platforms
- **Compression Type** - format and compression level
- **Creation Metadata** - date, criage version
- **Checksums** - for integrity verification

### Usage Examples

#### Creating Archive with Metadata

```bash
# Build package with automatic metadata embedding
criage build --format tar.zst --compression 6

# Result: test-package-1.0.0.tar.zst with embedded metadata
```

#### Viewing Archive Metadata

```bash
# Show all archive metadata
criage metadata test-package-1.0.0.tar.zst

# Example output:
# === Archive Metadata test-package-1.0.0.tar.zst ===
# Compression Type: tar.zst
# Created: 2024-01-15T10:30:45Z
# Created with: criage/1.0.0
# 
# === Package Manifest ===
# Name: test-package
# Version: 1.0.0
# Description: Test package
# Author: Developer Name
# License: MIT
# Dependencies:
#   - some-lib: ^1.0.0
# 
# === Build Manifest ===
# Build Script: echo Building...
# Output Directory: ./build
# Compression Format: tar.zst (level 6)
# Target Platforms:
#   - linux/amd64
#   - linux/arm64
```

### Benefits of Metadata Embedding

1. **Self-Sufficiency** - archive contains all necessary information
2. **Fast Access** - no need to extract for information retrieval
3. **Standards Compliance** - uses standard archive format mechanisms
4. **Compatibility** - works with any archivers supporting PAX
5. **Security** - built-in checksums for integrity verification

### Technical Details

#### Metadata Structure

```json
{
  "package_manifest": {
    "name": "test-package",
    "version": "1.0.0",
    "dependencies": {...}
  },
  "build_manifest": {
    "build_script": "echo Building...",
    "compression": {...}
  },
  "compression_type": "tar.zst",
  "created_at": "2024-01-15T10:30:45Z",
  "created_by": "criage/1.0.0",
  "checksum": "sha256:..."
}
```

#### Location in Archive

- **TAR**: PAX Extended Headers at beginning of archive
- **ZIP**: Archive comment + separate `.criage_metadata.json` file

## Performance

Criage is optimized for maximum performance:

- **Zstandard compression** - up to 3x faster than gzip with better compression
- **LZ4 compression** - extremely fast compression/decompression
- **Parallel processing** - utilizes all available CPU cores
- **Smart caching** - avoids repeated downloads
- **Efficient dependency resolution** - minimizes network requests

## Compression Format Comparison

| Format | Compression Speed | Decompression Speed | Size | Use Case |
|--------|------------------|-------------------|------|----------|
| tar.zst | Medium | Very Fast | Excellent | Default |
| tar.lz4 | Very Fast | Very Fast | Average | Fast operations |
| tar.xz | Slow | Medium | Excellent | Minimal size |
| tar.gz | Medium | Medium | Good | Compatibility |
| zip | Medium | Fast | Good | Windows compatibility |

## Development

### Requirements

- Go 1.21 or higher
- Git

### Building from Source

```bash
git clone https://github.com/Zu-Krein/criage.git
cd criage
go mod tidy
go build -o criage
```

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

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
