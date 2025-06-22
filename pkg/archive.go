package pkg

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/ulikunitz/xz"
)

// ArchiveManager управляет созданием и извлечением архивов
type ArchiveManager struct {
	config       *Config
	zstdEncoder  *zstd.Encoder
	zstdDecoder  *zstd.Decoder
	encoderMutex sync.Mutex
	decoderMutex sync.Mutex
}

// NewArchiveManager создает новый менеджер архивов
func NewArchiveManager(config *Config) (*ArchiveManager, error) {
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