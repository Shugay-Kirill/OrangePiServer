package yandexapi

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	. "telegramBot/yandexapi/init"
	. "telegramBot/yandexapi/method"
)

type YandexAPI struct {
	yaapi *YandexDiskAPI
}

// func testApi(api *YandexDiskAPI){

// }

// UploadFile читает локальный файл и передает его для загрузки на Яндекс.Диск
func (api *YandexAPI) UploadFile(remoutePathDirectory string, pathFile string) error {

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
	err = PostResourcesUpload(api.yaapi, remoutePathDirectory, fileData, contentType, fileInfo.Size(), pathFile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки файла через PostResourcesUpload: %v", err)
	}

	fmt.Printf("✅ Файл успешно передан для загрузки: %s → %s\n", pathFile, remoutePathDirectory)
	return nil
}

func (api *YandexAPI) CreateDirectory(pathDirectory string, nameDirectory string) error {
	directory, err := PutResources(api.yaapi, pathDirectory, nameDirectory)

	if err != nil {
		return err
	}

	fmt.Printf("\n directory = %s\n", directory)
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

// PrintDirectoryContents выводит содержимое директории
func (api *YandexAPI) PrintDirectoryContents(path string) error {
	files, err := GetResources(api.yaapi, path)
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

// PrintDiskUsage выводит информацию о использовании диска
func (api *YandexAPI) PrintDiskUsage() error {
	info, err := GetDiskInfo(api.yaapi)
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
