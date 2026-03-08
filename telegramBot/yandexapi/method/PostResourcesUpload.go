package method

import (
	"bytes"
	"fmt"

	"telegramBot/yandexapi/authenticated"
)

// PostResourcesUpload загружает файл на Яндекс.Диск
func PostResourcesUpload(remotePathDirectory string, fileData []byte, contentType string, fileSize int64, fileName string) error {

	uploadURL, err := GetResourcesUpload(remotePathDirectory, fileName)
	if err != nil {
		return fmt.Errorf("ошибка получения upload URL: %v", err)
	}

	fmt.Printf("🔗 Получен POST upload URL для: %s\n", remotePathDirectory)
	fmt.Printf("📤 Загрузка файла: %d bytes, тип: %s\n", len(fileData), contentType)

	// Создаем reader для файловых данных
	reader := bytes.NewReader(fileData)

	// Используем authenticatedRequest для загрузки файла (PUT запрос)
	responseBody, err := authenticated.AuthenticatedRequest("PUT", uploadURL, nil, reader)
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
