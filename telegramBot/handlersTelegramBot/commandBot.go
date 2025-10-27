package handlersTelegramBot

import (
	"fmt"

	"telegramBot/models"
)

func (h *MessageHandler) HandleStartCommand(update models.Update) {
	message := update.Message
	response := fmt.Sprintf(`🎉 <b>Добро пожаловать, %s!</b>

🤖 <b>Я - умный Telegram бот с различными возможностями</b>

✨ <b>Основные команды:</b>
• /start - показать это сообщение
• /help - получить помощь
• /features - возможности бота  
• /info - информация о чате

🛠️ <b>Что я умею:</b>
✅ Отвечать в том же топике/разделе
✅ Работать в группах и личных сообщениях
✅ Проверять фотографии и документы
✅ Определять JPG изображения
✅ Показывать детальную информацию
✅ Настраиваемый вывод логов API

⚙️ <b>Конфигурация:</b>
• Макс. длина вывода API: <b>%d символов</b>

📊 <b>Информация о текущем сообщении:</b>
• 👤 Ваше имя: <b>%s</b>
• 🆔 Ваш ID: <code>%d</code>
• 💬 ID чата: <code>%d</code>
• 🏷️ ID топика: <code>%d</code>`,
		message.From.FirstName,
		h.Config.MaxLengthAPIOutput,
		message.From.FirstName,
		message.From.ID,
		message.Chat.ID,
		message.MessageThreadID,
	)

	h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
}

func (h *MessageHandler) HandleHelpCommand(update models.Update) {
	message := update.Message
	response := fmt.Sprintf(`🆘 <b>Помощь по боту</b>

📚 <b>Доступные команды:</b>
• /start - начать работу с ботом
• /help - показать эту справку
• /features - возможности бота
• /info - информация о текущем чате

⚙️ <b>Настройки:</b>
• Максимальная длина вывода API: <b>%d символов</b>`,
		h.Config.MaxLengthAPIOutput,
	)

	h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
}

func (h *MessageHandler) HandleFeaturesCommand(update models.Update) {
	message := update.Message
	response := fmt.Sprintf(`🚀 <b>Возможности бота</b>

🔧 <b>Технические возможности:</b>
• <b>Настраиваемый вывод API</b> - Макс. длина логов: <b>%d символов</b>

⚙️ <b>Настройки конфигурации:</b>
• MAX_LENGTH_MESSEGE_API - Макс. длина вывода API (текущее значение: %d)`,
		h.Config.MaxLengthAPIOutput,
		h.Config.MaxLengthAPIOutput,
	)

	h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
}

func (h *MessageHandler) HandleInfoCommand(update models.Update) {
	message := update.Message

	chatType := "Неизвестный"
	switch message.Chat.Type {
	case "private":
		chatType = "💬 Личные сообщения"
	case "group":
		chatType = "👥 Группа"
	case "supergroup":
		chatType = "🌟 Супергруппа"
	case "channel":
		chatType = "📢 Канал"
	}

	topicStatus := "❌ Нет (основной чат)"
	if message.MessageThreadID != 0 {
		topicStatus = fmt.Sprintf("✅ Да (ID: %d)", message.MessageThreadID)
	}

	response := fmt.Sprintf(`ℹ️ <b>Информация о чате</b>

📋 <b>Основная информация:</b>
• 💬 Тип чата: <b>%s</b>
• 🏷️ Название: <b>%s</b>
• 🆔 ID чата: <code>%d</code>
• 🏷️ Топик: %s

🔧 <b>Техническая информация:</b>
• Макс. длина API логов: <b>%d символов</b>`,
		chatType,
		h.getChatTitle(message.Chat),
		message.Chat.ID,
		topicStatus,
		h.Config.MaxLengthAPIOutput,
	)

	h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
}
