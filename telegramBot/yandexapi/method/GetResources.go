package method

import (
	"encoding/json"
	"fmt"

	. "telegramBot/yandexapi/authenticated"
	. "telegramBot/yandexapi/init"
)

// GetResources возвращает список файлов в указанной директории
func GetResources(api *YandexDiskAPI, path string) ([]map[string]interface{}, error) {
	params := map[string]string{
		"path":  "testApi",
		"files": "name,path,type",
		"limit": "1000",
	}

	url := BuildURL(api, "/resources", params)

	body, err := AuthenticatedRequest(api, "GET", url, nil)
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
