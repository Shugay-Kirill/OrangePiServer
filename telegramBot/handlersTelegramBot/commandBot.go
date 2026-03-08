package handlersTelegramBot

import (
	"fmt"
	"log"

	"telegramBot/models"
	"telegramBot/yandexapi"
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

func (h *MessageHandler) HandleInfoDiskCommand(update models.Update) {
	message := update.Message

	info, err := yandexapi.PrintDiskUsage()

	if err != nil {
		log.Printf("ERROR PrintDiskUsage: %v", err)
		response := fmt.Sprintf("❌ Ошибка получения информации о диске: %v", err)
		h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
		return
	}

	fmt.Println(info)
	response := fmt.Sprintf(info)
	h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
}

func (h *MessageHandler) HandleCreateDirectory(update models.Update) {
	message := update.Message
	chatID := message.Chat.ID
	thredID := message.MessageThreadID

	h.SendMessage(chatID, thredID, "📁 Введите путь для создания директории (например, /photos):")
	step1 := func(text string) (string, InputHandler, error) {
		path := text
		print("step1. text %s", path)

		step2 := func(name string) (string, InputHandler, error) {
			print("step2. name %s", name)
			err := yandexapi.CreateDirectory(path, name)

			if err != nil {
				return "", nil, fmt.Errorf("не удалось создать директорию: %w", err)
			}
			return "✅ Директория успешно создана!", nil, nil
		}
		return "📁 Введите имя новой директории:", step2, nil
	}
	h.states.Store(chatID, &UserState{handler: step1})
}

func (h *MessageHandler) HandleDeleteDirectory(update models.Update) {
	message := update.Message
	chatID := message.Chat.ID
	thredID := message.MessageThreadID

	h.SendMessage(chatID, thredID, "📁 Введите путь для удаления директории (например, /photos):")
	step1 := func(text string) (string, InputHandler, error) {
		path := text
		print("step1. text %s", path)

		step2 := func(name string) (string, InputHandler, error) {
			print("step2. name %s", name)
			err := yandexapi.DeleteDirectory(path, name)

			if err != nil {
				return "", nil, fmt.Errorf("не удалось удалить директорию: %w", err)
			}
			return "✅ Директория успешно удалена!", nil, nil
		}
		return "📁 Введите имя директории для удаления:", step2, nil
	}
	h.states.Store(chatID, &UserState{handler: step1})
}

func (h *MessageHandler) HandleContentsDirectory(update models.Update) {
	message := update.Message
	chatID := message.Chat.ID
	thredID := message.MessageThreadID

	h.SendMessage(chatID, thredID, "📁 Введите путь для просмотра содержимого директории (например, /photos):")
	step1 := func(text string) (string, InputHandler, error) {
		path := text
		print("step1. text %s", path)

		err := yandexapi.PrintDirectoryContents(path)

		if err != nil {
			return "", nil, fmt.Errorf("не удалось просмотреть директорию: %w", err)
		}
		return "✅ Содержимое директории", nil, nil
	}
	h.states.Store(chatID, &UserState{handler: step1})
}
