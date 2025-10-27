package api

import (
	"fmt"
	"strings"

	"telegramBot/models"
)

// AnalyzePhoto анализирует фотографию и возвращает информацию о ней
func AnalyzePhoto(photo models.PhotoSize) models.FileAnalysisResponse {
	return photo.ToFileAnalysisResponse()
}

// AnalyzeDocument анализирует документ и проверяет, является ли он JPG
func AnalyzeDocument(document models.Document) models.FileAnalysisResponse {
	request := document.ToFileAnalysisRequest()

	isJPG := isJPGImage(request)
	fileType := "document"
	if isJPG {
		fileType = "image/jpeg"
	}

	return models.FileAnalysisResponse{
		IsJPG:      isJPG,
		FileType:   fileType,
		FileSizeKB: formatFileSize(document.FileSize),
	}
}

// isJPGImage проверяет, является ли документ JPG изображением
func isJPGImage(request models.FileAnalysisRequest) bool {
	if strings.HasPrefix(request.MimeType, "image/jpeg") {
		return true
	}

	fileName := strings.ToLower(request.FileName)
	if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg") {
		return true
	}

	if request.MimeType == "image/jpg" {
		return true
	}

	return false
}

func formatFileSize(size int) string {
	return fmt.Sprintf("%.2f KB", float64(size)/1024)
}

// Дополнительные API методы для внешних сервисов
func HealthCheck() models.APIResponse {
	return models.APIResponse{
		Success: true,
		Data:    "Service is healthy",
	}
}

func GetServiceInfo() models.APIResponse {
	return models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"name":    "Telegram Bot API Service",
			"version": "1.0.0",
			"features": []string{
				"file_analysis",
				"jpg_detection",
				"photo_processing",
			},
		},
	}
}
