package main

import (
	"fmt"
	"os"

	"criage/pkg"

	"github.com/spf13/cobra"
)

var packageManager *pkg.PackageManager

func init() {
	var err error
	packageManager, err = pkg.NewPackageManager()
	if err != nil {
		fmt.Printf("Ошибка инициализации пакетного менеджера: %v\n", err)
		os.Exit(1)
	}
}

// installPackage устанавливает пакет
func installPackage(packageName string) error {
	return packageManager.InstallPackage(packageName, "", false, false, false, "", "")
}

// uninstallPackage удаляет пакет
func uninstallPackage(packageName string) error {
	return packageManager.UninstallPackage(packageName, false, false)
}

// updatePackage обновляет пакет
func updatePackage(packageName string) error {
	return packageManager.UpdatePackage(packageName)
}

// updateAllPackages обновляет все пакеты
func updateAllPackages() error {
	packages, err := packageManager.ListPackages(false, true)
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		if err := packageManager.UpdatePackage(pkg.Name); err != nil {
			fmt.Printf("Не удалось обновить %s: %v\n", pkg.Name, err)
		}
	}
	return nil
}

// searchPackages выполняет поиск пакетов
func searchPackages(query string) error {
	results, err := packageManager.SearchPackages(query)
	if err != nil {
		return err
	}

	fmt.Printf("Найдено %d пакетов:\n", len(results))
	for _, result := range results {
		fmt.Printf("- %s (%s): %s\n", result.Name, result.Version, result.Description)
	}
	return nil
}

// listPackages показывает список установленных пакетов
func listPackages() error {
	packages, err := packageManager.ListPackages(false, false)
	if err != nil {
		return err
	}

	fmt.Printf("Установлено %d пакетов:\n", len(packages))
	for _, pkg := range packages {
		fmt.Printf("- %s (%s)\n", pkg.Name, pkg.Version)
	}
	return nil
}

// showPackageInfo показывает информацию о пакете
func showPackageInfo(packageName string) error {
	info, err := packageManager.GetPackageInfo(packageName)
	if err != nil {
		return err
	}

	fmt.Printf("Название: %s\n", info.Name)
	fmt.Printf("Версия: %s\n", info.Version)
	fmt.Printf("Описание: %s\n", info.Description)
	fmt.Printf("Автор: %s\n", info.Author)
	fmt.Printf("Установлен: %s\n", info.InstallDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("Размер: %d байт\n", info.Size)
	return nil
}

// createPackage создает новый пакет
func createPackage(name string) error {
	return packageManager.CreatePackage(name, "basic", "", "")
}

// publishPackage публикует пакет
func publishPackage() error {
	return packageManager.PublishPackage("", "")
}

// showArchiveMetadata показывает метаданные архива
func showArchiveMetadata(archivePath string) error {
	archiveManager, err := pkg.NewArchiveManager(pkg.DefaultConfig(), version)
	if err != nil {
		return fmt.Errorf("failed to create archive manager: %w", err)
	}
	defer archiveManager.Close()

	format := archiveManager.DetectFormat(archivePath)
	metadata, err := archiveManager.ExtractMetadataFromArchive(archivePath, format)
	if err != nil {
		return fmt.Errorf("failed to extract metadata: %w", err)
	}

	fmt.Printf("=== Метаданные архива %s ===\n", archivePath)
	fmt.Printf("Тип сжатия: %s\n", metadata.CompressionType)
	fmt.Printf("Создан: %s\n", metadata.CreatedAt)
	fmt.Printf("Создано с помощью: %s\n", metadata.CreatedBy)

	if metadata.PackageManifest != nil {
		fmt.Printf("\n=== Манифест пакета ===\n")
		fmt.Printf("Название: %s\n", metadata.PackageManifest.Name)
		fmt.Printf("Версия: %s\n", metadata.PackageManifest.Version)
		fmt.Printf("Описание: %s\n", metadata.PackageManifest.Description)
		fmt.Printf("Автор: %s\n", metadata.PackageManifest.Author)
		fmt.Printf("Лицензия: %s\n", metadata.PackageManifest.License)
		if len(metadata.PackageManifest.Dependencies) > 0 {
			fmt.Printf("Зависимости:\n")
			for name, version := range metadata.PackageManifest.Dependencies {
				fmt.Printf("  - %s: %s\n", name, version)
			}
		}
	}

	if metadata.BuildManifest != nil {
		fmt.Printf("\n=== Манифест сборки ===\n")
		fmt.Printf("Скрипт сборки: %s\n", metadata.BuildManifest.BuildScript)
		fmt.Printf("Выходная директория: %s\n", metadata.BuildManifest.OutputDir)
		fmt.Printf("Формат сжатия: %s (уровень %d)\n",
			metadata.BuildManifest.Compression.Format,
			metadata.BuildManifest.Compression.Level)
		if len(metadata.BuildManifest.Targets) > 0 {
			fmt.Printf("Целевые платформы:\n")
			for _, target := range metadata.BuildManifest.Targets {
				fmt.Printf("  - %s/%s\n", target.OS, target.Arch)
			}
		}
	}

	return nil
}

// setConfig устанавливает значение конфигурации
func setConfig(key, value string) error {
	fmt.Printf("Установка конфигурации %s = %s\n", key, value)
	// Здесь будет реализация установки конфигурации
	return nil
}

// getConfig получает значение конфигурации
func getConfig(key string) error {
	fmt.Printf("Получение значения конфигурации для ключа: %s\n", key)
	// Здесь будет реализация получения конфигурации
	return nil
}

// listConfig показывает все настройки
func listConfig() error {
	fmt.Println("Список всех настроек конфигурации:")
	// Здесь будет реализация показа всех настроек
	return nil
}

// getCurrentCommand возвращает текущую команду для доступа к флагам
func getCurrentCommand() *cobra.Command {
	// Это вспомогательная функция для получения текущей команды
	// В реальной реализации нужно передавать команду через контекст
	return &cobra.Command{}
}

// formatSize форматирует размер в человеко-читаемый формат
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
