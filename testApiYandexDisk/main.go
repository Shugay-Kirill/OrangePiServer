package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type YandexDiskAPI struct {
	hostNameURL string
	token       string
	client      *http.Client
}

type DiskInfo struct {
	TotalSpace int64 `json:"total_space"`
	UsedSpace  int64 `json:"used_space"`
	TrashSize  int64 `json:"trash_size"`
}

// Config хранит конфигурацию из env переменных
type Config struct {
	URL   string
	Token string
}

// LoadConfig загружает конфигурацию из env переменных
func LoadConfig() *Config {
	_ = godotenv.Load()

	config := &Config{
		URL:   os.Getenv("YANDEX_DISK_URL"),
		Token: os.Getenv("YANDEX_DISK_TOKEN"),
	}

	// Если URL не указан, используем значение по умолчанию
	if config.URL == "" {
		config.URL = "https://cloud-api.yandex.net/v1/disk"
		log.Println("Using default Yandex.Disk API URL")
	}

	// Токен
	if config.Token == "" {
		log.Fatal("YANDEX_DISK_TOKEN environment variable is required")
	}

	log.Printf("Configuration loaded - URL: %s", config.URL)

	return config
}

func NewYandexDiskAPI() *YandexDiskAPI {
	config := LoadConfig()

	return &YandexDiskAPI{
		hostNameURL: config.URL,
		token:       config.Token,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: false,
			},
		},
	}
}

// BuildURL создает URL с параметрами
func (api *YandexDiskAPI) BuildURL(pathNameURl string, search map[string]string) string {
	uri := api.hostNameURL + pathNameURl

	if len(search) > 0 {
		query := url.Values{}
		for key, value := range search {
			fmt.Printf("\nkey = %s, value = %s\n", key, value)
			query.Add(key, value)
		}
		fmt.Printf("\nURL = %s\n", uri)
		uri += "?" + query.Encode()
	}
	fmt.Printf("\nURL = %s\n", uri)
	return uri
}

// authenticatedRequest выполняет авторизованный запрос , statusCode uint8
func (api *YandexDiskAPI) authenticatedRequest(method, url string, body io.Reader) ([]byte, error) {

	requestApi, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	// Устанавливаем заголовки авторизации
	requestApi.Header.Set("Authorization", "OAuth "+api.token)
	requestApi.Header.Set("Accept", "application/json")
	if body != nil && method == "PUT" {
		requestApi.Header.Set("Content-Type", "application/octet-stream")
	}
	if body != nil && method == "POST" {
		requestApi.Header.Set("Content-Type", "application/json")
	}

	log.Printf("🔗 Making %s request to: %s", method, url)

	responseApi, err := api.client.Do(requestApi)
	if err != nil {
		return nil, err
	}
	defer responseApi.Body.Close()

	responseBody, err := io.ReadAll(responseApi.Body)
	if err != nil {
		return nil, err
	}

	if responseApi.StatusCode != http.StatusOK {
		// Пытаемся извлечь message из JSON тела ответа
		var errorResponse struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Error   string `json:"error"`
		}
		if err := json.Unmarshal(responseBody, &errorResponse); err == nil {
			// Если удалось распарсить JSON, используем message из него
			errorMsg := errorResponse.Message
			if errorMsg == "" {
				errorMsg = errorResponse.Error
			}
			if errorMsg == "" {
				errorMsg = string(responseBody)
			}
			return nil, fmt.Errorf("API error %d: %s", responseApi.StatusCode, errorMsg)
		}

		// Если не JSON, выводим тело как есть
		return nil, fmt.Errorf("API error %d: %s", responseApi.StatusCode, string(responseBody))
	}

	log.Printf("✅ Request successful (Status: %d)", responseApi.StatusCode)
	return responseBody, nil
}

// // GetDiskInfo получает информацию о диске
func (api *YandexDiskAPI) GetDiskInfo() (*DiskInfo, error) {
	url := api.BuildURL("", nil)

	body, err := api.authenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var diskInfo DiskInfo
	err = json.Unmarshal(body, &diskInfo)
	if err != nil {
		return nil, err
	}

	return &diskInfo, nil
}

// GetResources возвращает список файлов в указанной директории
func (api *YandexDiskAPI) GetResources(path string) ([]map[string]interface{}, error) {
	params := map[string]string{
		"path":  path,
		"files": "name,path,type",
		"limit": "1000",
	}

	url := api.BuildURL("/resources", params)

	body, err := api.authenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	embedded, ok := result["_embedded"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	items, ok := embedded["items"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil // Пустая директория
	}

	var files []map[string]interface{}
	for _, item := range items {
		if file, ok := item.(map[string]interface{}); ok {
			files = append(files, file)
		}
	}

	return files, nil
}

// PrintDiskUsage выводит информацию о использовании диска
func (api *YandexDiskAPI) PrintDiskUsage() error {
	info, err := api.GetDiskInfo()
	if err != nil {
		return err
	}

	fmt.Println("💾 Yandex.Disk Usage:")
	fmt.Println("├── Total Space:", formatBytes(info.TotalSpace))
	fmt.Println("├── Used Space: ", formatBytes(info.UsedSpace))
	fmt.Println("└── Free Space: ", formatBytes(info.TotalSpace-info.UsedSpace))

	usagePercent := float64(info.UsedSpace) / float64(info.TotalSpace) * 100
	fmt.Printf("Usage: %.1f%%\n", usagePercent)

	return nil
}

// PrintDirectoryContents выводит содержимое директории
func (api *YandexDiskAPI) PrintDirectoryContents(path string) error {
	files, err := api.GetResources(path)
	if err != nil {
		return err
	}

	fmt.Printf("📁 Contents of '%s':\n", path)
	fmt.Println(strings.Repeat("─", 60))

	for _, file := range files {
		name, _ := file["name"].(string)
		fileType, _ := file["type"].(string)

		if fileType == "file" {
			size, _ := file["size"].(float64)
			path, _ := file["path"].(string)
			fmt.Printf("📄 %-30s %10s %s\n", name, formatBytes(int64(size)), path)
		} else {
			fmt.Printf("📁 %s/\n", name)
		}
	}

	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("Total items: %d\n", len(files))

	return nil
}

// Вспомогательная функция для форматирования байтов
func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	sizes := []string{"B", "KB", "MB", "GB", "TB"}
	i := 0
	size := float64(bytes)

	for size >= 1024 && i < len(sizes)-1 {
		size /= 1024
		i++
	}

	return fmt.Sprintf("%.1f %s", size, sizes[i])
}

func (api *YandexDiskAPI) PutResources(pathDirectory string, nameDirectory string) (map[string]interface{}, error) {

	params := map[string]string{
		"path": pathDirectory + "/" + nameDirectory,
	}

	url := api.BuildURL("/resources", params)

	body, err := api.authenticatedRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (api *YandexDiskAPI) CreateDirectory(pathDirectory string, nameDirectory string) error {
	directory, err := api.PutResources(pathDirectory, nameDirectory)

	if err != nil {
		return err
	}

	fmt.Printf("\n directory = %s\n", directory)
	return nil
}

// GetResourcesUpload получает URL для загрузки файла на Яндекс.Диск
func (api *YandexDiskAPI) GetResourcesUpload(remotePathDirectory string, fileName string) (string, error) {
	params := map[string]string{
		"path": remotePathDirectory + "/" + fileName,
	}

	url := api.BuildURL("/resources/upload", params)

	fmt.Printf("🔗 Запрос GET upload URL для: %s\n", remotePathDirectory)

	body, err := api.authenticatedRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка получения upload URL: %v", err)
	}

	var response struct {
		HREF string `json:"href"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга ответа: %v", err)
	}

	if response.HREF == "" {
		return "", fmt.Errorf("пустой upload URL в ответе")
	}

	fmt.Printf("✅ Получен upload URL: %s\n", response.HREF)
	return response.HREF, nil
}

// PostResourcesUpload загружает файл на Яндекс.Диск
func (api *YandexDiskAPI) PostResourcesUpload(remotePathDirectory string, fileData []byte, contentType string, fileSize int64, fileName string) error {

	uploadURL, err := api.GetResourcesUpload(remotePathDirectory, fileName)
	if err != nil {
		return fmt.Errorf("ошибка получения upload URL: %v", err)
	}

	fmt.Printf("🔗 Получен POST upload URL для: %s\n", remotePathDirectory)
	fmt.Printf("📤 Загрузка файла: %d bytes, тип: %s\n", len(fileData), contentType)

	// Создаем reader для файловых данных
	reader := bytes.NewReader(fileData)

	// Используем authenticatedRequest для загрузки файла (PUT запрос)
	responseBody, err := api.authenticatedRequest("PUT", uploadURL, reader)
	if err != nil {
		return fmt.Errorf("ошибка загрузки файла через authenticatedRequest: %v", err)
	}

	fmt.Printf("✅ Файл успешно загружен на Яндекс.Диск\n")

	// Логируем ответ от сервера (если есть)
	if len(responseBody) > 0 {
		fmt.Printf("📨 Ответ от Яндекс.Диска: %s\n", string(responseBody))
	}

	return nil
}

// UploadFile читает локальный файл и передает его для загрузки на Яндекс.Диск
func (api *YandexDiskAPI) UploadFile(remoutePathDirectory string, pathFile string) error {
	// Проверяем существование файла
	fileInfo, err := os.Stat(pathFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("файл не существует: %s", pathFile)
	}

	fmt.Printf("📖 Чтение файла: %s (%d bytes)\n", pathFile, fileInfo.Size())

	// Читаем файл в память
	fileData, err := os.ReadFile(pathFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	// Получаем MIME-тип файла
	contentType := http.DetectContentType(fileData)
	fmt.Printf("📄 MIME-тип файла: %s\n", contentType)

	// Передаем данные файла в PostResourcesUpload для загрузки
	err = api.PostResourcesUpload(remoutePathDirectory, fileData, contentType, fileInfo.Size(), pathFile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки файла через PostResourcesUpload: %v", err)
	}

	fmt.Printf("✅ Файл успешно передан для загрузки: %s → %s\n", pathFile, remoutePathDirectory)
	return nil
}

func main() {
	fmt.Println("🚀 Yandex.Disk API Client")
	fmt.Println()

	// Создаем API клиент (конфигурация загружается автоматически)
	api := NewYandexDiskAPI()

	// Получаем информацию о диске
	// err := api.PrintDiskUsage()
	// if err != nil {
	// 	log.Fatalf("❌ Error getting disk info: %v", err)
	// }

	fmt.Println()

	// Получаем содержимое корневой директории
	// err := api.PrintDirectoryContents("disktest:/")
	// if err != nil {
	// 	log.Fatalf("❌ Error listing files: %v", err)
	// }

	// Создание папки
	// err := api.CreateDirectory("testApi", "testApiCreateDir")
	// if err != nil {
	// 	log.Fatalf("❌ Error create directory: %v", err)
	// }

	// Загрузка файла
	err := api.PrintDirectoryContents("/testApi")
	if err != nil {
		log.Fatalf("❌ Error upload file: %v", err)
	}

}
