package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// ApiServer представляет API сервер репозитория
type ApiServer struct {
	indexManager *IndexManager
	config       *Config
	router       *mux.Router
}

// NewApiServer создает новый API сервер
func NewApiServer(indexManager *IndexManager, config *Config) *ApiServer {
	server := &ApiServer{
		indexManager: indexManager,
		config:       config,
		router:       mux.NewRouter(),
	}

	server.setupRoutes()
	return server
}

// setupRoutes настраивает маршруты API
func (s *ApiServer) setupRoutes() {
	// Middleware для CORS если включен
	if s.config.EnableCORS {
		s.router.Use(s.corsMiddleware)
	}

	// Middleware для логирования
	s.router.Use(s.loggingMiddleware)

	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Информация о репозитории
	api.HandleFunc("/", s.handleRepositoryInfo).Methods("GET")
	api.HandleFunc("/stats", s.handleStats).Methods("GET")

	// Управление пакетами
	api.HandleFunc("/packages", s.handleListPackages).Methods("GET")
	api.HandleFunc("/packages/{name}", s.handleGetPackage).Methods("GET")
	api.HandleFunc("/packages/{name}/{version}", s.handleGetPackageVersion).Methods("GET")

	// Поиск
	api.HandleFunc("/search", s.handleSearchPackages).Methods("GET")

	// Скачивание
	api.HandleFunc("/download/{name}/{version}/{file}", s.handleDownload).Methods("GET")

	// Загрузка (требует токен)
	api.HandleFunc("/upload", s.handleUpload).Methods("POST")

	// Обновление индекса (требует токен)
	api.HandleFunc("/refresh", s.handleRefreshIndex).Methods("POST")

	// Статические файлы для веб интерфейса
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
}

// corsMiddleware добавляет CORS заголовки
func (s *ApiServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware логирует HTTP запросы
func (s *ApiServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// handleRepositoryInfo возвращает информацию о репозитории
func (s *ApiServer) handleRepositoryInfo(w http.ResponseWriter, r *http.Request) {
	index := s.indexManager.GetIndex()

	info := map[string]interface{}{
		"name":           "Criage Package Repository",
		"version":        "1.0.0",
		"last_updated":   index.LastUpdated,
		"total_packages": index.TotalPackages,
		"formats":        s.config.AllowedFormats,
	}

	s.sendJSONResponse(w, http.StatusOK, true, info, "")
}

// handleStats возвращает статистику репозитория
func (s *ApiServer) handleStats(w http.ResponseWriter, r *http.Request) {
	index := s.indexManager.GetIndex()
	s.sendJSONResponse(w, http.StatusOK, true, index.Statistics, "")
}

// handleListPackages возвращает список всех пакетов
func (s *ApiServer) handleListPackages(w http.ResponseWriter, r *http.Request) {
	index := s.indexManager.GetIndex()

	// Параметры пагинации
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Получаем список пакетов
	var packages []*PackageEntry
	for _, pkg := range index.Packages {
		packages = append(packages, pkg)
	}

	// Пагинация
	start := (page - 1) * limit
	end := start + limit
	if start > len(packages) {
		start = len(packages)
	}
	if end > len(packages) {
		end = len(packages)
	}

	result := map[string]interface{}{
		"packages":    packages[start:end],
		"total":       len(packages),
		"page":        page,
		"limit":       limit,
		"total_pages": (len(packages) + limit - 1) / limit,
	}

	s.sendJSONResponse(w, http.StatusOK, true, result, "")
}

// handleGetPackage возвращает информацию о конкретном пакете
func (s *ApiServer) handleGetPackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	index := s.indexManager.GetIndex()
	pkg, exists := index.Packages[name]
	if !exists {
		s.sendJSONResponse(w, http.StatusNotFound, false, nil, "Package not found")
		return
	}

	s.sendJSONResponse(w, http.StatusOK, true, pkg, "")
}

// handleGetPackageVersion возвращает информацию о конкретной версии пакета
func (s *ApiServer) handleGetPackageVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	index := s.indexManager.GetIndex()
	pkg, exists := index.Packages[name]
	if !exists {
		s.sendJSONResponse(w, http.StatusNotFound, false, nil, "Package not found")
		return
	}

	for _, v := range pkg.Versions {
		if v.Version == version {
			s.sendJSONResponse(w, http.StatusOK, true, v, "")
			return
		}
	}

	s.sendJSONResponse(w, http.StatusNotFound, false, nil, "Version not found")
}

// handleSearchPackages выполняет поиск пакетов
func (s *ApiServer) handleSearchPackages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		s.sendJSONResponse(w, http.StatusBadRequest, false, nil, "Search query is required")
		return
	}

	results := s.indexManager.SearchPackages(query)

	// Ограничиваем результаты
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	if len(results) > limit {
		results = results[:limit]
	}

	response := map[string]interface{}{
		"query":   query,
		"results": results,
		"total":   len(results),
	}

	s.sendJSONResponse(w, http.StatusOK, true, response, "")
}

// handleDownload обрабатывает скачивание пакетов
func (s *ApiServer) handleDownload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]
	filename := vars["file"]

	// Проверяем, существует ли пакет в индексе
	index := s.indexManager.GetIndex()
	pkg, exists := index.Packages[name]
	if !exists {
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}

	// Найдем версию и файл
	var fileEntry *FileEntry
	for _, v := range pkg.Versions {
		if v.Version == version {
			for _, f := range v.Files {
				if f.Filename == filename {
					fileEntry = &f
					break
				}
			}
			break
		}
	}

	if fileEntry == nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Путь к файлу
	filePath := filepath.Join(s.config.StoragePath, filename)

	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found on disk", http.StatusNotFound)
		return
	}

	// Увеличиваем счетчик загрузок
	go func() {
		if err := s.indexManager.IncrementDownload(name, version); err != nil {
			log.Printf("Failed to increment download counter: %v", err)
		}
	}()

	// Устанавливаем заголовки
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", strconv.FormatInt(fileEntry.Size, 10))
	w.Header().Set("X-Checksum", fileEntry.Checksum)

	// Отправляем файл
	http.ServeFile(w, r, filePath)
}

// handleUpload обрабатывает загрузку новых пакетов
func (s *ApiServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	// Проверяем токен
	token := r.Header.Get("Authorization")
	if token != "Bearer "+s.config.UploadToken {
		s.sendJSONResponse(w, http.StatusUnauthorized, false, nil, "Invalid token")
		return
	}

	// Проверяем размер файла
	r.Body = http.MaxBytesReader(w, r.Body, s.config.MaxFileSize)

	// Парсим multipart form
	if err := r.ParseMultipartForm(s.config.MaxFileSize); err != nil {
		s.sendJSONResponse(w, http.StatusBadRequest, false, nil, "Failed to parse form: "+err.Error())
		return
	}

	// Получаем файл
	file, header, err := r.FormFile("package")
	if err != nil {
		s.sendJSONResponse(w, http.StatusBadRequest, false, nil, "No package file provided")
		return
	}
	defer file.Close()

	// Проверяем формат файла
	filename := header.Filename
	if !s.isAllowedFormat(filename) {
		s.sendJSONResponse(w, http.StatusBadRequest, false, nil, "Unsupported file format")
		return
	}

	// Сохраняем файл
	destPath := filepath.Join(s.config.StoragePath, filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		s.sendJSONResponse(w, http.StatusInternalServerError, false, nil, "Failed to create file")
		return
	}
	defer destFile.Close()

	// Копируем содержимое
	size, err := io.Copy(destFile, file)
	if err != nil {
		os.Remove(destPath)
		s.sendJSONResponse(w, http.StatusInternalServerError, false, nil, "Failed to save file")
		return
	}

	log.Printf("Uploaded package: %s (%d bytes)", filename, size)

	// Обновляем индекс асинхронно
	go func() {
		if err := s.indexManager.scanPackages(); err != nil {
			log.Printf("Failed to update index after upload: %v", err)
		}
	}()

	s.sendJSONResponse(w, http.StatusCreated, true, map[string]interface{}{
		"filename": filename,
		"size":     size,
	}, "Package uploaded successfully")
}

// handleRefreshIndex принудительно обновляет индекс
func (s *ApiServer) handleRefreshIndex(w http.ResponseWriter, r *http.Request) {
	// Проверяем токен
	token := r.Header.Get("Authorization")
	if token != "Bearer "+s.config.UploadToken {
		s.sendJSONResponse(w, http.StatusUnauthorized, false, nil, "Invalid token")
		return
	}

	// Сканируем пакеты
	if err := s.indexManager.scanPackages(); err != nil {
		s.sendJSONResponse(w, http.StatusInternalServerError, false, nil, "Failed to refresh index: "+err.Error())
		return
	}

	index := s.indexManager.GetIndex()
	s.sendJSONResponse(w, http.StatusOK, true, map[string]interface{}{
		"total_packages": index.TotalPackages,
		"last_updated":   index.LastUpdated,
	}, "Index refreshed successfully")
}

// isAllowedFormat проверяет, разрешен ли формат файла
func (s *ApiServer) isAllowedFormat(filename string) bool {
	for _, format := range s.config.AllowedFormats {
		if strings.HasSuffix(strings.ToLower(filename), "."+format) {
			return true
		}
	}
	return false
}

// sendJSONResponse отправляет JSON ответ
func (s *ApiServer) sendJSONResponse(w http.ResponseWriter, statusCode int, success bool, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ApiResponse{
		Success: success,
		Data:    data,
		Message: message,
	}

	if !success && message != "" {
		response.Error = message
	}

	json.NewEncoder(w).Encode(response)
}

// Start запускает HTTP сервер
func (s *ApiServer) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("Starting repository server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
