package yandexapi

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"telegramBot/yandexapi/method"
)

// UploadFile читает локальный файл и передает его для загрузки на Яндекс.Диск
func UploadFile(remoutePathDirectory string, pathFile string) error {

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
	err = method.PostResourcesUpload(remoutePathDirectory, fileData, contentType, fileInfo.Size(), pathFile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки файла через PostResourcesUpload: %v", err)
	}

	fmt.Printf("✅ Файл успешно передан для загрузки: %s → %s\n", pathFile, remoutePathDirectory)
	return nil
}

// CreateDirectory создание директории
func CreateDirectory(pathDirectory string, nameDirectory string) error {
	print("start createDir. pathDirectory = %s, nameDirectory = %s", pathDirectory, nameDirectory)
	directory, err := method.PutResources(pathDirectory, nameDirectory)

	if err == nil {
		return err
	}

	fmt.Printf("\n directory = %s\n", directory)
	return nil
}

// CreateDirectory создание директории
func DeleteDirectory(pathDirectory string, nameDirectory string) error {
	print("start DelteDirectory. pathDirectory = %s, nameDirectory = %s", pathDirectory, nameDirectory)
	directory, err := method.DeleteResources(pathDirectory, nameDirectory)

	if err == nil {
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
func PrintDirectoryContents(pathDirectory string) error {
	files, err := method.GetResources(pathDirectory)
	if err != nil {
		return err
	}

	fmt.Printf("📁 Contents of '%s':\n", pathDirectory)
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
func PrintDiskUsage() (string, error) {
	fmt.Println("start PrintDiskUsage")

	info, err := method.GetDiskInfo()
	if err != nil {
		return "", err
	}

	usagePercent := float64(info.UsedSpace) / float64(info.TotalSpace) * 100
	// fmt.Printf("Usage: %.1f%%\n", usagePercent)

	result := fmt.Sprintf(`
💾 Yandex.Disk Usage:
├── Total Space: %s
├── Used Space: %s
└── Free Space: %s
Usage: %.1f%%`,
		formatBytes(info.TotalSpace),
		formatBytes(info.UsedSpace),
		formatBytes(info.TotalSpace-info.UsedSpace),
		usagePercent,
	)
	fmt.Println("end PrintDiskUsage")
	return result, nil
}
