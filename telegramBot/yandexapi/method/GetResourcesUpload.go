package method

import (
	"encoding/json"
	"fmt"

	. "telegramBot/yandexapi/authenticated"
	. "telegramBot/yandexapi/init"
)

// GetResourcesUpload получает URL для загрузки файла на Яндекс.Диск
func GetResourcesUpload(api *YandexDiskAPI, remotePathDirectory string, fileName string) (string, error) {
	params := map[string]string{
		"path": remotePathDirectory + "/" + fileName,
	}

	url := BuildURL(api, "/resources/upload", params)

	fmt.Printf("🔗 Запрос GET upload URL для: %s\n", remotePathDirectory)

	body, err := AuthenticatedRequest(api, "GET", url, nil)
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
