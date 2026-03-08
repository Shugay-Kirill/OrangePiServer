package authenticated

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"telegramBot/yandexapi/initYD"
)

// authenticatedRequest выполняет авторизованный запрос , statusCode uint8
func AuthenticatedRequest(method, pathUrl string, parametr map[string]string, body io.Reader) ([]byte, error) {

	apiAuth := initYD.GetYandexDiskAPI()
	url := BuildURL(pathUrl, parametr)

	requestApi, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	// Устанавливаем заголовки авторизации
	requestApi.Header.Set("Authorization", "OAuth "+apiAuth.Token)
	requestApi.Header.Set("Accept", "application/json")
	if body != nil && method == "PUT" {
		requestApi.Header.Set("Content-Type", "application/octet-stream")
	}
	if body != nil && method == "POST" {
		requestApi.Header.Set("Content-Type", "application/json")
	}

	log.Printf("🔗 Making %s request to: %s", method, url)

	responseApi, err := apiAuth.Client.Do(requestApi)
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
