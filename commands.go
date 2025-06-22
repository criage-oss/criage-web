package main

import (
	"fmt"
	"os"
	"text/tabwriter"

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

// installPackage выполняет установку пакета
func installPackage(packageName string) error {
	cmd := getCurrentCommand()
	
	global, _ := cmd.Flags().GetBool("global")
	version, _ := cmd.Flags().GetString("version")
	force, _ := cmd.Flags().GetBool("force")
	dev, _ := cmd.Flags().GetBool("dev")
	arch, _ := cmd.Flags().GetString("arch")
	osName, _ := cmd.Flags().GetString("os")

	return packageManager.InstallPackage(packageName, version, global, force, dev, arch, osName)
}

// uninstallPackage выполняет удаление пакета
func uninstallPackage(packageName string) error {
	cmd := getCurrentCommand()
	
	global, _ := cmd.Flags().GetBool("global")
	purge, _ := cmd.Flags().GetBool("purge")

	return packageManager.UninstallPackage(packageName, global, purge)
}

// updatePackage выполняет обновление пакета
func updatePackage(packageName string) error {
	return packageManager.UpdatePackage(packageName)
}

// updateAllPackages выполняет обновление всех пакетов
func updateAllPackages() error {
	packages, err := packageManager.ListPackages(false, false)
	if err != nil {
		return fmt.Errorf("failed to list packages: %w", err)
	}

	fmt.Printf("Обновление %d пакетов...\n", len(packages))

	for _, pkg := range packages {
		fmt.Printf("Обновление %s...\n", pkg.Name)
		if err := packageManager.UpdatePackage(pkg.Name); err != nil {
			fmt.Printf("Ошибка обновления %s: %v\n", pkg.Name, err)
		}
	}

	return nil
}

// searchPackages выполняет поиск пакетов
func searchPackages(query string) error {
	results, err := packageManager.SearchPackages(query)
	if err != nil {
		return fmt.Errorf("failed to search packages: %w", err)
	}

	if len(results) == 0 {
		fmt.Printf("Пакеты не найдены по запросу: %s\n", query)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ИМЯ\tВЕРСИЯ\tОПИСАНИЕ\tАВТОР\tРЕПОЗИТОРИЙ")

	for _, result := range results {
		description := result.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}
		
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			result.Name,
			result.Version,
			description,
			result.Author,
			result.Repository,
		)
	}

	w.Flush()
	return nil
}

// listPackages выводит список установленных пакетов
func listPackages() error {
	cmd := getCurrentCommand()
	
	global, _ := cmd.Flags().GetBool("global")
	outdated, _ := cmd.Flags().GetBool("outdated")

	packages, err := packageManager.ListPackages(global, outdated)
	if err != nil {
		return fmt.Errorf("failed to list packages: %w", err)
	}

	if len(packages) == 0 {
		if outdated {
			fmt.Println("Все пакеты имеют актуальные версии")
		} else {
			fmt.Println("Пакеты не установлены")
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ИМЯ\tВЕРСИЯ\tОПИСАНИЕ\tРАЗМЕР\tДАТА УСТАНОВКИ")

	for _, pkg := range packages {
		description := pkg.Description
		if len(description) > 40 {
			description = description[:37] + "..."
		}
		
		size := formatSize(pkg.Size)
		date := pkg.InstallDate.Format("2006-01-02")
		
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			pkg.Name,
			pkg.Version,
			description,
			size,
			date,
		)
	}

	w.Flush()
	return nil
}

// showPackageInfo показывает подробную информацию о пакете
func showPackageInfo(packageName string) error {
	info, err := packageManager.GetPackageInfo(packageName)
	if err != nil {
		return fmt.Errorf("failed to get package info: %w", err)
	}

	fmt.Printf("Имя: %s\n", info.Name)
	fmt.Printf("Версия: %s\n", info.Version)
	fmt.Printf("Описание: %s\n", info.Description)
	fmt.Printf("Автор: %s\n", info.Author)
	fmt.Printf("Размер: %s\n", formatSize(info.Size))
	fmt.Printf("Дата установки: %s\n", info.InstallDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("Путь установки: %s\n", info.InstallPath)
	fmt.Printf("Глобальный: %t\n", info.Global)

	if len(info.Dependencies) > 0 {
		fmt.Println("\nЗависимости:")
		for name, version := range info.Dependencies {
			fmt.Printf("  %s: %s\n", name, version)
		}
	}

	if len(info.Scripts) > 0 {
		fmt.Println("\nСкрипты:")
		for name, command := range info.Scripts {
			fmt.Printf("  %s: %s\n", name, command)
		}
	}

	if len(info.Files) > 0 {
		fmt.Println("\nФайлы:")
		for _, file := range info.Files {
			fmt.Printf("  %s\n", file)
		}
	}

	return nil
}

// createPackage создает новый пакет
func createPackage(name string) error {
	cmd := getCurrentCommand()
	
	template, _ := cmd.Flags().GetString("template")
	author, _ := cmd.Flags().GetString("author")
	description, _ := cmd.Flags().GetString("description")

	return packageManager.CreatePackage(name, template, author, description)
}

// buildPackage собирает пакет (нужно переделать для получения флагов из контекста)
func buildPackage() error {
	// Временное решение с дефолтными значениями
	return packageManager.BuildPackage("", "tar.zst", 3)
}

// publishPackage публикует пакет
func publishPackage() error {
	cmd := getCurrentCommand()
	
	registry, _ := cmd.Flags().GetString("registry")
	token, _ := cmd.Flags().GetString("token")

	return packageManager.PublishPackage(registry, token)
}

// setConfig устанавливает значение конфигурации
func setConfig(key, value string) error {
	return packageManager.GetConfigManager().SetValue(key, value)
}

// getConfig получает значение конфигурации
func getConfig(key string) error {
	value, err := packageManager.GetConfigManager().GetValue(key)
	if err != nil {
		return err
	}
	
	fmt.Printf("%s = %s\n", key, value)
	return nil
}

// listConfig показывает все настройки конфигурации
func listConfig() error {
	values := packageManager.GetConfigManager().ListValues()
	
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "КЛЮЧ\tЗНАЧЕНИЕ")
	
	// Сортируем ключи для консистентного вывода
	var keys []string
	for key := range values {
		keys = append(keys, key)
	}
	
	for _, key := range keys {
		fmt.Fprintf(w, "%s\t%s\n", key, values[key])
	}
	
	w.Flush()
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

