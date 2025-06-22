package main

import (
	"time"
)

// PackageEntry представляет запись о пакете в индексе репозитория
type PackageEntry struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Author        string         `json:"author"`
	License       string         `json:"license"`
	Homepage      string         `json:"homepage"`
	Repository    string         `json:"repository"`
	Keywords      []string       `json:"keywords"`
	Versions      []VersionEntry `json:"versions"`
	LatestVersion string         `json:"latest_version"`
	Downloads     int64          `json:"downloads"`
	Updated       time.Time      `json:"updated"`
}

// VersionEntry представляет конкретную версию пакета
type VersionEntry struct {
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Dependencies map[string]string `json:"dependencies"`
	DevDeps      map[string]string `json:"dev_dependencies"`
	Files        []FileEntry       `json:"files"`
	Size         int64             `json:"size"`
	Checksum     string            `json:"checksum"`
	Uploaded     time.Time         `json:"uploaded"`
	Downloads    int64             `json:"downloads"`
}

// FileEntry представляет файл пакета для разных платформ
type FileEntry struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Format   string `json:"format"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Checksum string `json:"checksum"`
}

// RepositoryIndex представляет индекс всех пакетов в репозитории
type RepositoryIndex struct {
	LastUpdated   time.Time                `json:"last_updated"`
	TotalPackages int                      `json:"total_packages"`
	Packages      map[string]*PackageEntry `json:"packages"`
	Statistics    *Statistics              `json:"statistics"`
}

// Statistics статистика репозитория
type Statistics struct {
	TotalDownloads    int64          `json:"total_downloads"`
	PackagesByLicense map[string]int `json:"packages_by_license"`
	PackagesByAuthor  map[string]int `json:"packages_by_author"`
	PopularPackages   []string       `json:"popular_packages"`
}

// SearchResult результат поиска пакетов
type SearchResult struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Downloads   int64     `json:"downloads"`
	Updated     time.Time `json:"updated"`
	Score       float64   `json:"score"`
}

// UploadRequest запрос на загрузку пакета
type UploadRequest struct {
	Token       string `json:"token"`
	PackageName string `json:"package_name"`
	Version     string `json:"version"`
	OS          string `json:"os"`
	Arch        string `json:"arch"`
	Format      string `json:"format"`
}

// ApiResponse стандартный ответ API
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Config конфигурация сервера репозитория
type Config struct {
	Port           int      `json:"port"`
	StoragePath    string   `json:"storage_path"`
	IndexPath      string   `json:"index_path"`
	UploadToken    string   `json:"upload_token"`
	MaxFileSize    int64    `json:"max_file_size"`
	AllowedFormats []string `json:"allowed_formats"`
	EnableCORS     bool     `json:"enable_cors"`
	LogLevel       string   `json:"log_level"`
}
