package pkg

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/ulikunitz/xz"
)

// PackageMetadata представляет метаданные пакета для встраивания в архив
type PackageMetadata struct {
	PackageManifest *PackageManifest `json:"package_manifest"`
	BuildManifest   *BuildManifest   `json:"build_manifest"`
	CompressionType string           `json:"compression_type"`
	CreatedAt       string           `json:"created_at"`
	CreatedBy       string           `json:"created_by"`
	Checksum        string           `json:"checksum,omitempty"`
}

// ArchiveManager управляет созданием и извлечением архивов
type ArchiveManager struct {
	config       *Config
	zstdEncoder  *zstd.Encoder
	zstdDecoder  *zstd.Decoder
	encoderMutex sync.Mutex
	decoderMutex sync.Mutex
	version      string
}

// NewArchiveManager создает новый менеджер архивов
func NewArchiveManager(config *Config, version string) (*ArchiveManager, error) {
	// Создаем zstd encoder с оптимальными настройками
	encoder, err := zstd.NewWriter(nil,
		zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(config.Compression.Level)),
		zstd.WithEncoderConcurrency(config.Parallel),
		zstd.WithWindowSize(1<<20), // 1MB window для лучшего сжатия
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create zstd encoder: %w", err)
	}

	// Создаем zstd decoder
	decoder, err := zstd.NewReader(nil,
		zstd.WithDecoderConcurrency(config.Parallel),
		zstd.WithDecoderMaxWindow(1<<26), // 64MB max window
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create zstd decoder: %w", err)
	}

	return &ArchiveManager{
		config:      config,
		zstdEncoder: encoder,
		zstdDecoder: decoder,
		version:     version,
	}, nil
}

// CreateArchive создает архив из указанных файлов
func (am *ArchiveManager) CreateArchive(srcPath, dstPath string, format ArchiveFormat, files []string, exclude []string) error {
	switch format {
	case FormatTarZst:
		return am.createTarZst(srcPath, dstPath, files, exclude)
	case FormatTarLz4:
		return am.createTarLz4(srcPath, dstPath, files, exclude)
	case FormatTarXz:
		return am.createTarXz(srcPath, dstPath, files, exclude)
	case FormatTarGz:
		return am.createTarGz(srcPath, dstPath, files, exclude)
	case FormatZip:
		return am.createZip(srcPath, dstPath, files, exclude)
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}
}

// ExtractArchive извлекает архив
func (am *ArchiveManager) ExtractArchive(srcPath, dstPath string, format ArchiveFormat) error {
	switch format {
	case FormatTarZst:
		return am.extractTarZst(srcPath, dstPath)
	case FormatTarLz4:
		return am.extractTarLz4(srcPath, dstPath)
	case FormatTarXz:
		return am.extractTarXz(srcPath, dstPath)
	case FormatTarGz:
		return am.extractTarGz(srcPath, dstPath)
	case FormatZip:
		return am.extractZip(srcPath, dstPath)
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}
}

// DetectFormat определяет формат архива по расширению
func (am *ArchiveManager) DetectFormat(filename string) ArchiveFormat {
	switch {
	case strings.HasSuffix(filename, ".tar.zst"):
		return FormatTarZst
	case strings.HasSuffix(filename, ".tar.lz4"):
		return FormatTarLz4
	case strings.HasSuffix(filename, ".tar.xz"):
		return FormatTarXz
	case strings.HasSuffix(filename, ".tar.gz") || strings.HasSuffix(filename, ".tgz"):
		return FormatTarGz
	case strings.HasSuffix(filename, ".zip"):
		return FormatZip
	default:
		return FormatTarZst // По умолчанию
	}
}

// CreateArchiveWithMetadata создает архив с встроенными метаданными
func (am *ArchiveManager) CreateArchiveWithMetadata(srcPath, dstPath string, format ArchiveFormat, files []string, exclude []string, metadata *PackageMetadata) error {
	metadata.CompressionType = string(format)
	metadata.CreatedAt = time.Now().Format(time.RFC3339)
	metadata.CreatedBy = "criage/" + am.version // версия будет передаваться из main.go

	switch format {
	case FormatTarZst, FormatTarLz4, FormatTarXz, FormatTarGz:
		return am.createTarWithMetadata(srcPath, dstPath, format, files, exclude, metadata)
	case FormatZip:
		return am.createZipWithMetadata(srcPath, dstPath, files, exclude, metadata)
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}
}

// ExtractMetadataFromArchive извлекает метаданные из архива
func (am *ArchiveManager) ExtractMetadataFromArchive(archivePath string, format ArchiveFormat) (*PackageMetadata, error) {
	switch format {
	case FormatTarZst, FormatTarLz4, FormatTarXz, FormatTarGz:
		return am.extractTarMetadata(archivePath, format)
	case FormatZip:
		return am.extractZipMetadata(archivePath)
	default:
		return nil, fmt.Errorf("unsupported archive format: %s", format)
	}
}

// createTarZst создает TAR архив со сжатием Zstandard
func (am *ArchiveManager) createTarZst(srcPath, dstPath string, files []string, exclude []string) error {
	outFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	am.encoderMutex.Lock()
	am.zstdEncoder.Reset(outFile)
	defer func() {
		am.zstdEncoder.Close()
		am.encoderMutex.Unlock()
	}()

	tw := tar.NewWriter(am.zstdEncoder)
	defer tw.Close()

	return am.addFilesToTar(tw, srcPath, files, exclude)
}

// createTarLz4 создает TAR архив со сжатием LZ4
func (am *ArchiveManager) createTarLz4(srcPath, dstPath string, files []string, exclude []string) error {
	outFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	lz4Writer := lz4.NewWriter(outFile)
	defer lz4Writer.Close()

	// LZ4 не поддерживает настройку уровня сжатия в этой версии библиотеки

	tw := tar.NewWriter(lz4Writer)
	defer tw.Close()

	return am.addFilesToTar(tw, srcPath, files, exclude)
}

// createTarXz создает TAR архив со сжатием XZ
func (am *ArchiveManager) createTarXz(srcPath, dstPath string, files []string, exclude []string) error {
	outFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	xzWriter, err := xz.NewWriter(outFile)
	if err != nil {
		return fmt.Errorf("failed to create xz writer: %w", err)
	}
	defer xzWriter.Close()

	tw := tar.NewWriter(xzWriter)
	defer tw.Close()

	return am.addFilesToTar(tw, srcPath, files, exclude)
}

// createTarGz создает TAR архив со сжатием Gzip
func (am *ArchiveManager) createTarGz(srcPath, dstPath string, files []string, exclude []string) error {
	outFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	gzWriter, err := gzip.NewWriterLevel(outFile, am.config.Compression.Level)
	if err != nil {
		return fmt.Errorf("failed to create gzip writer: %w", err)
	}
	defer gzWriter.Close()

	tw := tar.NewWriter(gzWriter)
	defer tw.Close()

	return am.addFilesToTar(tw, srcPath, files, exclude)
}

// createZip создает ZIP архив
func (am *ArchiveManager) createZip(srcPath, dstPath string, files []string, exclude []string) error {
	outFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// ZIP использует встроенное сжатие

	return am.addFilesToZip(zipWriter, srcPath, files, exclude)
}

// addFilesToTar добавляет файлы в TAR архив
func (am *ArchiveManager) addFilesToTar(tw *tar.Writer, srcPath string, files []string, exclude []string) error {
	excludeMap := make(map[string]bool)
	for _, ex := range exclude {
		excludeMap[ex] = true
	}

	for _, file := range files {
		fullPath := filepath.Join(srcPath, file)

		err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Проверяем исключения
			relPath, _ := filepath.Rel(srcPath, path)
			if excludeMap[relPath] {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			// Создаем заголовок
			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}

			header.Name = relPath

			// Записываем заголовок
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			// Если это файл, копируем содержимое
			if info.Mode().IsRegular() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				_, err = io.Copy(tw, file)
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// addFilesToZip добавляет файлы в ZIP архив
func (am *ArchiveManager) addFilesToZip(zw *zip.Writer, srcPath string, files []string, exclude []string) error {
	excludeMap := make(map[string]bool)
	for _, ex := range exclude {
		excludeMap[ex] = true
	}

	for _, file := range files {
		fullPath := filepath.Join(srcPath, file)

		err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Проверяем исключения
			relPath, _ := filepath.Rel(srcPath, path)
			if excludeMap[relPath] {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			// Пропускаем директории
			if info.IsDir() {
				return nil
			}

			// Создаем файл в архиве
			zipFile, err := zw.Create(relPath)
			if err != nil {
				return err
			}

			// Копируем содержимое файла
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(zipFile, file)
			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// extractTarZst извлекает TAR архив со сжатием Zstandard
func (am *ArchiveManager) extractTarZst(srcPath, dstPath string) error {
	inFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer inFile.Close()

	am.decoderMutex.Lock()
	am.zstdDecoder.Reset(inFile)
	defer am.decoderMutex.Unlock()

	tr := tar.NewReader(am.zstdDecoder)
	return am.extractTarContent(tr, dstPath)
}

// extractTarLz4 извлекает TAR архив со сжатием LZ4
func (am *ArchiveManager) extractTarLz4(srcPath, dstPath string) error {
	inFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer inFile.Close()

	lz4Reader := lz4.NewReader(inFile)
	tr := tar.NewReader(lz4Reader)
	return am.extractTarContent(tr, dstPath)
}

// extractTarXz извлекает TAR архив со сжатием XZ
func (am *ArchiveManager) extractTarXz(srcPath, dstPath string) error {
	inFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer inFile.Close()

	xzReader, err := xz.NewReader(inFile)
	if err != nil {
		return fmt.Errorf("failed to create xz reader: %w", err)
	}

	tr := tar.NewReader(xzReader)
	return am.extractTarContent(tr, dstPath)
}

// extractTarGz извлекает TAR архив со сжатием Gzip
func (am *ArchiveManager) extractTarGz(srcPath, dstPath string) error {
	inFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer inFile.Close()

	gzReader, err := gzip.NewReader(inFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tr := tar.NewReader(gzReader)
	return am.extractTarContent(tr, dstPath)
}

// extractZip извлекает ZIP архив
func (am *ArchiveManager) extractZip(srcPath, dstPath string) error {
	zipReader, err := zip.OpenReader(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if err := am.extractZipFile(file, dstPath); err != nil {
			return err
		}
	}

	return nil
}

// extractTarContent извлекает содержимое TAR архива
func (am *ArchiveManager) extractTarContent(tr *tar.Reader, dstPath string) error {
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dstPath, header.Name)

		// Проверяем, что путь находится внутри целевой директории
		if !strings.HasPrefix(target, filepath.Clean(dstPath)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			outFile, err := os.Create(target)
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()

			if err := os.Chmod(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		}
	}

	return nil
}

// extractZipFile извлекает отдельный файл из ZIP архива
func (am *ArchiveManager) extractZipFile(file *zip.File, dstPath string) error {
	target := filepath.Join(dstPath, file.Name)

	// Проверяем, что путь находится внутри целевой директории
	if !strings.HasPrefix(target, filepath.Clean(dstPath)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", file.Name)
	}

	if file.FileInfo().IsDir() {
		return os.MkdirAll(target, 0755)
	}

	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

// Close закрывает ресурсы менеджера архивов
func (am *ArchiveManager) Close() error {
	am.encoderMutex.Lock()
	defer am.encoderMutex.Unlock()

	am.decoderMutex.Lock()
	defer am.decoderMutex.Unlock()

	if am.zstdEncoder != nil {
		am.zstdEncoder.Close()
	}
	if am.zstdDecoder != nil {
		am.zstdDecoder.Close()
	}

	return nil
}

// createTarWithMetadata создает TAR архив с метаданными в PAX extended headers
func (am *ArchiveManager) createTarWithMetadata(srcPath, dstPath string, format ArchiveFormat, files []string, exclude []string, metadata *PackageMetadata) error {
	outFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Создаем компрессор в зависимости от формата
	var compressor io.WriteCloser
	switch format {
	case FormatTarZst:
		am.encoderMutex.Lock()
		am.zstdEncoder.Reset(outFile)
		compressor = am.zstdEncoder
		defer func() {
			am.zstdEncoder.Close()
			am.encoderMutex.Unlock()
		}()
	case FormatTarLz4:
		compressor = lz4.NewWriter(outFile)
		defer compressor.Close()
	case FormatTarXz:
		xzWriter, err := xz.NewWriter(outFile)
		if err != nil {
			return fmt.Errorf("failed to create xz writer: %w", err)
		}
		compressor = xzWriter
		defer compressor.Close()
	case FormatTarGz:
		gzWriter, err := gzip.NewWriterLevel(outFile, am.config.Compression.Level)
		if err != nil {
			return fmt.Errorf("failed to create gzip writer: %w", err)
		}
		compressor = gzWriter
		defer compressor.Close()
	}

	tw := tar.NewWriter(compressor)
	defer tw.Close()

	// Добавляем метаданные как PAX extended header в начале архива
	if err := am.addMetadataToPaxHeader(tw, metadata); err != nil {
		return fmt.Errorf("failed to add metadata: %w", err)
	}

	return am.addFilesToTar(tw, srcPath, files, exclude)
}

// addMetadataToPaxHeader добавляет метаданные в PAX extended header
func (am *ArchiveManager) addMetadataToPaxHeader(tw *tar.Writer, metadata *PackageMetadata) error {
	// Сериализуем метаданные в JSON
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Создаем PAX extended header
	paxHeaders := map[string]string{
		"criage.metadata": string(metadataJSON),
		"criage.version":  am.version,
	}

	// Добавляем информацию о манифестах если они есть
	if metadata.PackageManifest != nil {
		manifestJSON, _ := json.Marshal(metadata.PackageManifest)
		paxHeaders["criage.package_manifest"] = string(manifestJSON)
	}

	if metadata.BuildManifest != nil {
		buildJSON, _ := json.Marshal(metadata.BuildManifest)
		paxHeaders["criage.build_manifest"] = string(buildJSON)
	}

	// Создаем специальный файл с PAX headers
	hdr := &tar.Header{
		Name:       ".criage_metadata",
		Mode:       0644,
		Size:       0,
		ModTime:    time.Now(),
		Typeflag:   tar.TypeReg,
		PAXRecords: paxHeaders,
	}

	return tw.WriteHeader(hdr)
}

// createZipWithMetadata создает ZIP архив с метаданными в комментарии
func (am *ArchiveManager) createZipWithMetadata(srcPath, dstPath string, files []string, exclude []string, metadata *PackageMetadata) error {
	outFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Сериализуем метаданные в JSON и добавляем как комментарий к ZIP
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	zipWriter.SetComment(string(metadataJSON))

	// Также добавляем метаданные как отдельный файл в архив
	metadataFile, err := zipWriter.Create(".criage_metadata.json")
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}

	if _, err := metadataFile.Write(metadataJSON); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return am.addFilesToZip(zipWriter, srcPath, files, exclude)
}

// extractTarMetadata извлекает метаданные из TAR архива
func (am *ArchiveManager) extractTarMetadata(archivePath string, format ArchiveFormat) (*PackageMetadata, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	// Создаем декомпрессор в зависимости от формата
	var reader io.Reader
	switch format {
	case FormatTarZst:
		am.decoderMutex.Lock()
		am.zstdDecoder.Reset(file)
		reader = am.zstdDecoder
		defer am.decoderMutex.Unlock()
	case FormatTarLz4:
		reader = lz4.NewReader(file)
	case FormatTarXz:
		xzReader, err := xz.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("failed to create xz reader: %w", err)
		}
		reader = xzReader
	case FormatTarGz:
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	tr := tar.NewReader(reader)

	// Ищем метаданные в первых записях архива
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar header: %w", err)
		}

		// Проверяем PAX records на наличие метаданных
		if hdr.PAXRecords != nil {
			if metadataStr, exists := hdr.PAXRecords["criage.metadata"]; exists {
				var metadata PackageMetadata
				if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
					return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
				}
				return &metadata, nil
			}
		}

		// Если это файл метаданных
		if hdr.Name == ".criage_metadata" || hdr.Name == ".criage_metadata.json" {
			data, err := io.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("failed to read metadata file: %w", err)
			}

			var metadata PackageMetadata
			if err := json.Unmarshal(data, &metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
			return &metadata, nil
		}

		// Останавливаемся после нескольких записей, если метаданные не найдены
		if hdr.Name != ".criage_metadata" && !strings.HasPrefix(hdr.Name, ".") {
			break
		}
	}

	return nil, fmt.Errorf("metadata not found in archive")
}

// extractZipMetadata извлекает метаданные из ZIP архива
func (am *ArchiveManager) extractZipMetadata(archivePath string) (*PackageMetadata, error) {
	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer zipReader.Close()

	// Сначала пробуем извлечь из комментария ZIP
	if zipReader.Comment != "" {
		var metadata PackageMetadata
		if err := json.Unmarshal([]byte(zipReader.Comment), &metadata); err == nil {
			return &metadata, nil
		}
	}

	// Ищем файл метаданных
	for _, file := range zipReader.File {
		if file.Name == ".criage_metadata.json" {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open metadata file: %w", err)
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("failed to read metadata: %w", err)
			}

			var metadata PackageMetadata
			if err := json.Unmarshal(data, &metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
			return &metadata, nil
		}
	}

	return nil, fmt.Errorf("metadata not found in zip archive")
}
