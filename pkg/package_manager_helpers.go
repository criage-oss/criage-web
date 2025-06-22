package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// loadInstalledPackages загружает информацию об установленных пакетах
func (pm *PackageManager) loadInstalledPackages() error {
	// Загружаем из локальных директорий
	localPath := pm.configManager.GetConfig().LocalPath
	globalPath := pm.configManager.GetConfig().GlobalPath
	
	// Загружаем локальные пакеты
	if err := pm.loadPackagesFromDir(localPath, false); err != nil {
		fmt.Printf("Предупреждение: failed to load local packages: %v\n", err)
	}
	
	// Загружаем глобальные пакеты
	if err := pm.loadPackagesFromDir(globalPath, true); err != nil {
		fmt.Printf("Предупреждение: failed to load global packages: %v\n", err)
	}
	
	return nil
}

// loadPackagesFromDir загружает пакеты из директории
func (pm *PackageManager) loadPackagesFromDir(dir string, global bool) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil // Директория не существует, это нормально
	}
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		packageDir := filepath.Join(dir, entry.Name())
		infoPath := filepath.Join(packageDir, ".criage", "package.json")
		
		if _, err := os.Stat(infoPath); os.IsNotExist(err) {
			continue
		}
		
		var info PackageInfo
		data, err := os.ReadFile(infoPath)
		if err != nil {
			continue
		}
		
		if err := json.Unmarshal(data, &info); err != nil {
			continue
		}
		
		pm.packagesMutex.Lock()
		pm.installedPackages[info.Name] = &info
		pm.packagesMutex.Unlock()
	}
	
	return nil
}

// savePackageInfo сохраняет информацию о пакете
func (pm *PackageManager) savePackageInfo(info *PackageInfo) error {
	infoDir := filepath.Join(info.InstallPath, ".criage")
	if err := os.MkdirAll(infoDir, 0755); err != nil {
		return err
	}
	
	infoPath := filepath.Join(infoDir, "package.json")
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(infoPath, data, 0644)
}

// removePackageInfo удаляет информацию о пакете
func (pm *PackageManager) removePackageInfo(packageName string) error {
	info, exists := pm.getInstalledPackage(packageName)
	if !exists {
		return nil
	}
	
	infoPath := filepath.Join(info.InstallPath, ".criage", "package.json")
	return os.Remove(infoPath)
}

// copyFiles копирует файлы из исходной директории в целевую
func (pm *PackageManager) copyFiles(srcDir, dstDir string, files []string) error {
	for _, pattern := range files {
		matches, err := filepath.Glob(filepath.Join(srcDir, pattern))
		if err != nil {
			return err
		}
		
		for _, src := range matches {
			rel, err := filepath.Rel(srcDir, src)
			if err != nil {
				return err
			}
			
			dst := filepath.Join(dstDir, rel)
			
			if err := pm.copyFile(src, dst); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// copyFile копирует отдельный файл
func (pm *PackageManager) copyFile(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	
	if srcInfo.IsDir() {
		return os.MkdirAll(dst, srcInfo.Mode())
	}
	
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	if _, err := srcFile.WriteTo(dstFile); err != nil {
		return err
	}
	
	return os.Chmod(dst, srcInfo.Mode())
}

// calculateDirSize вычисляет размер директории
func (pm *PackageManager) calculateDirSize(dir string) int64 {
	var size int64
	
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	
	return size
}

// searchInRepository выполняет поиск в репозитории
func (pm *PackageManager) searchInRepository(repo Repository, query string) ([]SearchResult, error) {
	// Простая реализация поиска
	url := fmt.Sprintf("%s/search?q=%s", repo.URL, query)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	if repo.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+repo.AuthToken)
	}
	
	resp, err := pm.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed: HTTP %d", resp.StatusCode)
	}
	
	var results []SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	
	// Устанавливаем репозиторий для каждого результата
	for i := range results {
		results[i].Repository = repo.Name
	}
	
	return results, nil
}

// executeBuildScript выполняет скрипт сборки
func (pm *PackageManager) executeBuildScript(manifest *BuildManifest) error {
	cmd := exec.Command("sh", "-c", manifest.BuildScript)
	
	// Устанавливаем переменные окружения
	env := os.Environ()
	for key, value := range manifest.BuildEnv {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	cmd.Env = env
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// uploadPackage загружает пакет в репозиторий
func (pm *PackageManager) uploadPackage(registryURL, token, archivePath string, manifest *PackageManifest) error {
	// Простая реализация загрузки
	url := fmt.Sprintf("%s/upload", registryURL)
	
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Package-Name", manifest.Name)
	req.Header.Set("X-Package-Version", manifest.Version)
	
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	
	resp, err := pm.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed: HTTP %d", resp.StatusCode)
	}
	
	return nil
}

// GetConfigManager возвращает менеджер конфигурации
func (pm *PackageManager) GetConfigManager() *ConfigManager {
	return pm.configManager
}