package handlersTelegramBot

import (
	// "fmt"

	// "telegramBot/apiYandexDisk"
	"telegramBot/models"
)

// func (h *MessageHandler) HandlePhoto(update models.Update) {
// 	message := update.Message
// 	largestPhoto := h.getLargestPhoto(message.Photo)

// 	// Используем API для анализа фото
// 	analysis := api.AnalyzePhoto(largestPhoto)

// 	response := fmt.Sprintf(`📸 <b>Получено фото!</b>

// 🖼️ <b>Информация о фото:</b>
// • 📏 Размер: <b>%s</b>
// • 💾 Вес: <b>%s</b>
// • 🏷️ Тип: <b>%s</b>

// 📝 <b>Подпись:</b> %s

// 🎯 <i>Фото успешно обработано!</i>`,
// 		analysis.Dimensions,
// 		analysis.FileSizeKB,
// 		analysis.FileType,
// 		h.getCaptionText(message.Caption),
// 	)

// 	h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
// }

// func (h *MessageHandler) HandleDocument(update models.Update) {
// 	message := update.Message

// 	// Используем API для анализа документа
// 	analysis := api.AnalyzeDocument(message.Document)

// 	var status string
// 	if analysis.IsJPG {
// 		status = "✅ <b>Это JPG изображение!</b>"
// 	} else {
// 		status = "❌ <b>Это не JPG изображение</b>"
// 	}

// 	response := fmt.Sprintf(`📎 <b>Получен документ!</b>

// 📋 <b>Информация о файле:</b>
// • 📝 Имя: <code>%s</code>
// • 🏷️ MIME Type: <b>%s</b>
// • 💾 Размер: <b>%s</b>

// 📝 <b>Подпись:</b> %s

// %s

// 🎯 <i>Документ проверен на соответствие формату JPG!</i>`,
// 		message.Document.FileName,
// 		message.Document.MimeType,
// 		analysis.FileSizeKB,
// 		h.getCaptionText(message.Caption),
// 		status,
// 	)

// 	h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
// }

func (h *MessageHandler) HandleOtherMessage(update models.Update) {
	message := update.Message
	response := `🔮 <b>Получено сообщение другого типа!</b>

💡 <b>Что я умею проверять:</b>
• 📸 Фотографии (автоматически определяю как JPG)
• 📎 Документы (проверяю формат JPG)
• 💬 Текстовые сообщения`

	h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
}

func (h *MessageHandler) getLargestPhoto(photos []models.PhotoSize) models.PhotoSize {
	if len(photos) == 0 {
		return models.PhotoSize{}
	}

	largest := photos[0]
	for _, photo := range photos {
		if photo.FileSize > largest.FileSize {
			largest = photo
		}
	}
	return largest
}

// package api

// import (
// 	"fmt"
// 	"strings"

// 	"telegramBot/models"
// )

// // AnalyzePhoto анализирует фотографию и возвращает информацию о ней
// func AnalyzePhoto(photo models.PhotoSize) models.FileAnalysisResponse {
// 	return photo.ToFileAnalysisResponse()
// }

// // AnalyzeDocument анализирует документ и проверяет, является ли он JPG
// func AnalyzeDocument(document models.Document) models.FileAnalysisResponse {
// 	request := document.ToFileAnalysisRequest()

// 	isJPG := isJPGImage(request)
// 	fileType := "document"
// 	if isJPG {
// 		fileType = "image/jpeg"
// 	}

// 	return models.FileAnalysisResponse{
// 		IsJPG:      isJPG,
// 		FileType:   fileType,
// 		FileSizeKB: formatFileSize(document.FileSize),
// 	}
// }

// // isJPGImage проверяет, является ли документ JPG изображением
// func isJPGImage(request models.FileAnalysisRequest) bool {
// 	if strings.HasPrefix(request.MimeType, "image/jpeg") {
// 		return true
// 	}

// 	fileName := strings.ToLower(request.FileName)
// 	if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg") {
// 		return true
// 	}

// 	if request.MimeType == "image/jpg" {
// 		return true
// 	}

// 	return false
// }

// func formatFileSize(size int) string {
// 	return fmt.Sprintf("%.2f KB", float64(size)/1024)
// }

// // Дополнительные API методы для внешних сервисов
// func HealthCheck() models.APIResponse {
// 	return models.APIResponse{
// 		Success: true,
// 		Data:    "Service is healthy",
// 	}
// }

// func GetServiceInfo() models.APIResponse {
// 	return models.APIResponse{
// 		Success: true,
// 		Data: map[string]interface{}{
// 			"name":    "Telegram Bot API Service",
// 			"version": "1.0.0",
// 			"features": []string{
// 				"file_analysis",
// 				"jpg_detection",
// 				"photo_processing",
// 			},
// 		},
// 	}
// }
