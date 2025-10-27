package handlersTelegramBot

import (
	"fmt"

	"telegramBot/api"
	"telegramBot/models"
)

func (h *MessageHandler) HandlePhoto(update models.Update) {
	message := update.Message
	largestPhoto := h.getLargestPhoto(message.Photo)

	// Используем API для анализа фото
	analysis := api.AnalyzePhoto(largestPhoto)

	response := fmt.Sprintf(`📸 <b>Получено фото!</b>

🖼️ <b>Информация о фото:</b>
• 📏 Размер: <b>%s</b>
• 💾 Вес: <b>%s</b>
• 🏷️ Тип: <b>%s</b>

📝 <b>Подпись:</b> %s

🎯 <i>Фото успешно обработано!</i>`,
		analysis.Dimensions,
		analysis.FileSizeKB,
		analysis.FileType,
		h.getCaptionText(message.Caption),
	)

	h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
}

func (h *MessageHandler) HandleDocument(update models.Update) {
	message := update.Message

	// Используем API для анализа документа
	analysis := api.AnalyzeDocument(message.Document)

	var status string
	if analysis.IsJPG {
		status = "✅ <b>Это JPG изображение!</b>"
	} else {
		status = "❌ <b>Это не JPG изображение</b>"
	}

	response := fmt.Sprintf(`📎 <b>Получен документ!</b>

📋 <b>Информация о файле:</b>
• 📝 Имя: <code>%s</code>
• 🏷️ MIME Type: <b>%s</b>
• 💾 Размер: <b>%s</b>

📝 <b>Подпись:</b> %s

%s

🎯 <i>Документ проверен на соответствие формату JPG!</i>`,
		message.Document.FileName,
		message.Document.MimeType,
		analysis.FileSizeKB,
		h.getCaptionText(message.Caption),
		status,
	)

	h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
}

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
