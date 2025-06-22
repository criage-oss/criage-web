package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"criage/pkg"
)

// IndexManager управляет индексом пакетов
type IndexManager struct {
	config      *Config
	index       *RepositoryIndex
	storagePath string
	indexPath   string
}

// NewIndexManager создает новый менеджер индекса
func NewIndexManager(config *Config) (*IndexManager, error) {
	im := &IndexManager{
		config:      config,
		storagePath: config.StoragePath,
		indexPath:   config.IndexPath,
	}

	// Загружаем существующий индекс или создаем новый
	if err := im.loadIndex(); err != nil {
		log.Printf("Создание нового индекса: %v", err)
		im.index = &RepositoryIndex{
			LastUpdated:   time.Now(),
			TotalPackages: 0,
			Packages:      make(map[string]*PackageEntry),
			Statistics: &Statistics{
				TotalDownloads:    0,
				PackagesByLicense: make(map[string]int),
				PackagesByAuthor:  make(map[string]int),
				PopularPackages:   []string{},
			},
		}
	}

	// Сканируем директорию с пакетами для обновления индекса
	if err := im.scanPackages(); err != nil {
		return nil, fmt.Errorf("failed to scan packages: %w", err)
	}

	return im, nil
}

// loadIndex загружает индекс из файла
func (im *IndexManager) loadIndex() error {
	data, err := os.ReadFile(im.indexPath)
	if err != nil {
		return err
	}

	im.index = &RepositoryIndex{}
	return json.Unmarshal(data, im.index)
}

// saveIndex сохраняет индекс в файл
func (im *IndexManager) saveIndex() error {
	im.index.LastUpdated = time.Now()
	im.index.TotalPackages = len(im.index.Packages)

	// Обновляем статистику
	im.updateStatistics()

	data, err := json.MarshalIndent(im.index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(im.indexPath, data, 0644)
}

// scanPackages сканирует директорию с пакетами и обновляет индекс
func (im *IndexManager) scanPackages() error {
	return filepath.Walk(im.storagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории и не-архивные файлы
		if info.IsDir() || !im.isPackageFile(path) {
			return nil
		}

		// Проверяем, есть ли уже этот файл в индексе
		if im.isFileIndexed(path, info) {
			return nil
		}

		log.Printf("Найден новый пакет: %s", path)
		return im.addPackageFromFile(path)
	})
}

// isPackageFile проверяет, является ли файл пакетом
func (im *IndexManager) isPackageFile(filename string) bool {
	for _, format := range im.config.AllowedFormats {
		if strings.HasSuffix(strings.ToLower(filename), "."+format) {
			return true
		}
	}
	return false
}

// isFileIndexed проверяет, индексирован ли уже файл
func (im *IndexManager) isFileIndexed(path string, info os.FileInfo) bool {
	filename := filepath.Base(path)

	for _, pkg := range im.index.Packages {
		for _, version := range pkg.Versions {
			for _, file := range version.Files {
				if file.Filename == filename && file.Size == info.Size() {
					return true
				}
			}
		}
	}
	return false
}

// addPackageFromFile добавляет пакет в индекс из файла
func (im *IndexManager) addPackageFromFile(filePath string) error {
	// Вычисляем контрольную сумму
	checksum, err := im.calculateChecksum(filePath)
	if err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	// Получаем информацию о файле
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// Пытаемся извлечь метаданные из архива
	metadata, err := im.extractMetadata(filePath)
	if err != nil {
		log.Printf("Предупреждение: не удалось извлечь метаданные из %s: %v", filePath, err)
		// Пытаемся извлечь информацию из имени файла
		return im.addPackageFromFilename(filePath, checksum, fileInfo.Size())
	}

	// Парсим информацию о платформе из имени файла
	os, arch, format, err := im.parseFilename(filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to parse filename: %w", err)
	}

	// Добавляем в индекс
	return im.addPackageToIndex(metadata, filePath, checksum, fileInfo.Size(), os, arch, format)
}

// extractMetadata извлекает метаданные из архива пакета
func (im *IndexManager) extractMetadata(filePath string) (*pkg.PackageMetadata, error) {
	// Создаем архивный менеджер
	config := pkg.DefaultConfig()
	archiveManager, err := pkg.NewArchiveManager(config, "1.0.0")
	if err != nil {
		return nil, err
	}
	defer archiveManager.Close()

	// Определяем формат и извлекаем метаданные
	format := archiveManager.DetectFormat(filePath)
	return archiveManager.ExtractMetadataFromArchive(filePath, format)
}

// addPackageFromFilename добавляет пакет, если метаданные недоступны
func (im *IndexManager) addPackageFromFilename(filePath, checksum string, size int64) error {
	filename := filepath.Base(filePath)

	// Пытаемся распарсить имя файла (package-version-os-arch.format)
	parts := strings.Split(filename, "-")
	if len(parts) < 2 {
		return fmt.Errorf("invalid filename format: %s", filename)
	}

	name := parts[0]
	version := parts[1]

	os, arch, format, err := im.parseFilename(filename)
	if err != nil {
		return err
	}

	// Создаем базовую структуру метаданных
	metadata := &pkg.PackageMetadata{
		PackageManifest: &pkg.PackageManifest{
			Name:    name,
			Version: version,
		},
		BuildManifest: &pkg.BuildManifest{
			Name:    name,
			Version: version,
		},
	}

	return im.addPackageToIndex(metadata, filePath, checksum, size, os, arch, format)
}

// parseFilename парсит имя файла для извлечения платформы
func (im *IndexManager) parseFilename(filename string) (string, string, string, error) {
	// Определяем формат
	var format string
	for _, f := range im.config.AllowedFormats {
		if strings.HasSuffix(strings.ToLower(filename), "."+f) {
			format = f
			filename = strings.TrimSuffix(filename, "."+f)
			break
		}
	}

	if format == "" {
		return "", "", "", fmt.Errorf("unknown format")
	}

	// Пытаемся найти OS и Arch в имени файла
	parts := strings.Split(filename, "-")
	if len(parts) >= 4 {
		// Формат: package-version-os-arch
		return parts[len(parts)-2], parts[len(parts)-1], format, nil
	}

	// По умолчанию
	return "linux", "amd64", format, nil
}

// addPackageToIndex добавляет пакет в индекс
func (im *IndexManager) addPackageToIndex(metadata *pkg.PackageMetadata, filePath, checksum string, size int64, osName, arch, format string) error {
	if metadata.PackageManifest == nil {
		return fmt.Errorf("package manifest is nil")
	}

	manifest := metadata.PackageManifest
	name := manifest.Name
	version := manifest.Version

	// Получаем или создаем запись пакета
	pkg, exists := im.index.Packages[name]
	if !exists {
		pkg = &PackageEntry{
			Name:        name,
			Description: manifest.Description,
			Author:      manifest.Author,
			License:     manifest.License,
			Homepage:    manifest.Homepage,
			Repository:  manifest.Repository,
			Keywords:    manifest.Keywords,
			Versions:    []VersionEntry{},
			Downloads:   0,
			Updated:     time.Now(),
		}
		im.index.Packages[name] = pkg
	}

	// Проверяем, есть ли уже эта версия
	var versionEntry *VersionEntry
	for i := range pkg.Versions {
		if pkg.Versions[i].Version == version {
			versionEntry = &pkg.Versions[i]
			break
		}
	}

	// Создаем новую версию, если её нет
	if versionEntry == nil {
		newVersion := VersionEntry{
			Version:      version,
			Description:  manifest.Description,
			Dependencies: manifest.Dependencies,
			DevDeps:      manifest.DevDeps,
			Files:        []FileEntry{},
			Size:         size,
			Checksum:     checksum,
			Uploaded:     time.Now(),
			Downloads:    0,
		}
		pkg.Versions = append(pkg.Versions, newVersion)
		versionEntry = &pkg.Versions[len(pkg.Versions)-1]
	}

	// Добавляем файл
	fileEntry := FileEntry{
		OS:       osName,
		Arch:     arch,
		Format:   format,
		Filename: filepath.Base(filePath),
		Size:     size,
		Checksum: checksum,
	}
	versionEntry.Files = append(versionEntry.Files, fileEntry)

	// Обновляем последнюю версию
	pkg.LatestVersion = im.getLatestVersion(pkg.Versions)
	pkg.Updated = time.Now()

	// Сохраняем индекс
	return im.saveIndex()
}

// getLatestVersion определяет последнюю версию пакета
func (im *IndexManager) getLatestVersion(versions []VersionEntry) string {
	if len(versions) == 0 {
		return ""
	}

	// Простая сортировка по времени загрузки
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Uploaded.After(versions[j].Uploaded)
	})

	return versions[0].Version
}

// calculateChecksum вычисляет SHA256 хеш файла
func (im *IndexManager) calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// updateStatistics обновляет статистику репозитория
func (im *IndexManager) updateStatistics() {
	stats := im.index.Statistics

	// Сбрасываем счетчики
	stats.PackagesByLicense = make(map[string]int)
	stats.PackagesByAuthor = make(map[string]int)
	stats.TotalDownloads = 0

	// Подсчитываем статистику
	for _, pkg := range im.index.Packages {
		if pkg.License != "" {
			stats.PackagesByLicense[pkg.License]++
		}
		if pkg.Author != "" {
			stats.PackagesByAuthor[pkg.Author]++
		}
		stats.TotalDownloads += pkg.Downloads
	}

	// Обновляем популярные пакеты
	im.updatePopularPackages()
}

// updatePopularPackages обновляет список популярных пакетов
func (im *IndexManager) updatePopularPackages() {
	type packageDownloads struct {
		name      string
		downloads int64
	}

	var packages []packageDownloads
	for name, pkg := range im.index.Packages {
		packages = append(packages, packageDownloads{name, pkg.Downloads})
	}

	// Сортируем по количеству скачиваний
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].downloads > packages[j].downloads
	})

	// Берем топ-10
	im.index.Statistics.PopularPackages = []string{}
	limit := 10
	if len(packages) < limit {
		limit = len(packages)
	}

	for i := 0; i < limit; i++ {
		im.index.Statistics.PopularPackages = append(
			im.index.Statistics.PopularPackages,
			packages[i].name,
		)
	}
}

// GetIndex возвращает текущий индекс
func (im *IndexManager) GetIndex() *RepositoryIndex {
	return im.index
}

// SearchPackages выполняет поиск пакетов
func (im *IndexManager) SearchPackages(query string) []SearchResult {
	var results []SearchResult
	query = strings.ToLower(query)

	for _, pkg := range im.index.Packages {
		score := 0.0

		// Поиск в названии (высокий приоритет)
		if strings.Contains(strings.ToLower(pkg.Name), query) {
			score += 10.0
		}

		// Поиск в описании
		if strings.Contains(strings.ToLower(pkg.Description), query) {
			score += 5.0
		}

		// Поиск в ключевых словах
		for _, keyword := range pkg.Keywords {
			if strings.Contains(strings.ToLower(keyword), query) {
				score += 3.0
			}
		}

		// Поиск в авторе
		if strings.Contains(strings.ToLower(pkg.Author), query) {
			score += 2.0
		}

		if score > 0 {
			results = append(results, SearchResult{
				Name:        pkg.Name,
				Version:     pkg.LatestVersion,
				Description: pkg.Description,
				Author:      pkg.Author,
				Downloads:   pkg.Downloads,
				Updated:     pkg.Updated,
				Score:       score,
			})
		}
	}

	// Сортируем по релевантности
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// IncrementDownload увеличивает счетчик скачиваний
func (im *IndexManager) IncrementDownload(packageName, version string) error {
	pkg, exists := im.index.Packages[packageName]
	if !exists {
		return fmt.Errorf("package not found: %s", packageName)
	}

	pkg.Downloads++

	// Увеличиваем счетчик для конкретной версии
	for i := range pkg.Versions {
		if pkg.Versions[i].Version == version {
			pkg.Versions[i].Downloads++
			break
		}
	}

	return im.saveIndex()
}
