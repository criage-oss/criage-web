package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:     "criage",
		Short:   "Высокопроизводительный пакетный менеджер",
		Long:    "Criage - быстрый и эффективный пакетный менеджер для управления пакетами и архивами",
		Version: version,
	}

	// Команды управления пакетами
	rootCmd.AddCommand(
		newInstallCmd(),
		newUninstallCmd(),
		newUpdateCmd(),
		newSearchCmd(),
		newListCmd(),
		newInfoCmd(),
		newCreateCmd(),
		newBuildCmd(),
		newPublishCmd(),
		newConfigCmd(),
		newMetadataCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Команда установки пакетов
func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [package]",
		Short: "Установить пакет",
		Long:  "Установить пакет из репозитория или локального файла",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return installPackage(args[0])
		},
	}

	cmd.Flags().BoolP("global", "g", false, "Установить пакет глобально")
	cmd.Flags().StringP("version", "v", "", "Версия пакета для установки")
	cmd.Flags().BoolP("force", "f", false, "Принудительная установка")
	cmd.Flags().BoolP("dev", "d", false, "Установить dev зависимости")
	cmd.Flags().StringP("arch", "a", "", "Архитектура (x86_64, arm64)")
	cmd.Flags().StringP("os", "o", "", "Операционная система")

	return cmd
}

// Команда удаления пакетов
func newUninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall [package]",
		Short: "Удалить пакет",
		Long:  "Удалить установленный пакет",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return uninstallPackage(args[0])
		},
	}

	cmd.Flags().BoolP("global", "g", false, "Удалить глобальный пакет")
	cmd.Flags().BoolP("purge", "p", false, "Полное удаление с конфигурацией")

	return cmd
}

// Команда обновления пакетов
func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [package]",
		Short: "Обновить пакет",
		Long:  "Обновить пакет до последней версии",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return updateAllPackages()
			}
			return updatePackage(args[0])
		},
	}

	cmd.Flags().BoolP("global", "g", false, "Обновить глобальные пакеты")
	cmd.Flags().BoolP("all", "a", false, "Обновить все пакеты")

	return cmd
}

// Команда поиска пакетов
func newSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search [query]",
		Short: "Найти пакеты",
		Long:  "Найти пакеты в репозитории",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return searchPackages(args[0])
		},
	}
}

// Команда списка пакетов
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Показать установленные пакеты",
		Long:  "Показать список установленных пакетов",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listPackages()
		},
	}

	cmd.Flags().BoolP("global", "g", false, "Показать глобальные пакеты")
	cmd.Flags().BoolP("outdated", "o", false, "Показать устаревшие пакеты")

	return cmd
}

// Команда информации о пакете
func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info [package]",
		Short: "Информация о пакете",
		Long:  "Показать подробную информацию о пакете",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showPackageInfo(args[0])
		},
	}
}

// Команда создания пакета
func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Создать новый пакет",
		Long:  "Создать новый пакет с базовой структурой",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createPackage(args[0])
		},
	}

	cmd.Flags().StringP("template", "t", "basic", "Шаблон пакета")
	cmd.Flags().StringP("author", "a", "", "Автор пакета")
	cmd.Flags().StringP("description", "d", "", "Описание пакета")

	return cmd
}

// Команда сборки пакета
func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Собрать пакет",
		Long:  "Собрать пакет из исходников",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			format, _ := cmd.Flags().GetString("format")
			compression, _ := cmd.Flags().GetInt("compression")

			return packageManager.BuildPackage(output, format, compression)
		},
	}

	cmd.Flags().StringP("output", "o", "", "Выходной файл")
	cmd.Flags().StringP("format", "f", "tar.zst", "Формат архива")
	cmd.Flags().IntP("compression", "c", 3, "Уровень сжатия")

	return cmd
}

// Команда публикации пакета
func newPublishCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Опубликовать пакет",
		Long:  "Опубликовать пакет в репозитории",
		RunE: func(cmd *cobra.Command, args []string) error {
			return publishPackage()
		},
	}

	cmd.Flags().StringP("registry", "r", "", "URL репозитория")
	cmd.Flags().StringP("token", "t", "", "Токен авторизации")

	return cmd
}

// Команда конфигурации
func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Управление конфигурацией",
		Long:  "Управление настройками criage",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "set [key] [value]",
			Short: "Установить значение конфигурации",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return setConfig(args[0], args[1])
			},
		},
		&cobra.Command{
			Use:   "get [key]",
			Short: "Получить значение конфигурации",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return getConfig(args[0])
			},
		},
		&cobra.Command{
			Use:   "list",
			Short: "Показать все настройки",
			RunE: func(cmd *cobra.Command, args []string) error {
				return listConfig()
			},
		},
	)

	return cmd
}

// Команда просмотра метаданных архива
func newMetadataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "metadata [archive]",
		Short: "Показать метаданные архива",
		Long:  "Извлечь и показать встроенные метаданные из архива пакета",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showArchiveMetadata(args[0])
		},
	}
}
