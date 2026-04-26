package handlersTelegramBot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

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
	threadID := message.MessageThreadID

	h.SendMessage(chatID, threadID, "📁 Введите путь для создания директории (например, /photos):")
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
	threadID := message.MessageThreadID

	h.SendMessage(chatID, threadID, "📁 Введите путь для удаления директории (например, /photos):")
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
	threadID := message.MessageThreadID

	h.SendMessage(chatID, threadID, "📁 Введите путь для просмотра содержимого директории (например, /photos):")
	step1 := func(text string) (string, InputHandler, error) {
		path := text
		print("step1. text %s", path)

		files, err := yandexapi.PrintDirectoryContents(path)
		var builder strings.Builder

		fmt.Fprintf(&builder, "📁 Contents of '%s':\n", path)
		fmt.Fprintln(&builder, strings.Repeat("─", 20))
		for _, file := range files {
			name, _ := file["name"].(string)
			fileType, _ := file["type"].(string)

			if fileType == "file" {
				size, _ := file["size"].(float64)
				fmt.Fprintf(&builder, "📄 %-30s %10s \n", name, yandexapi.FormatBytes(int64(size)))
			} else {
				fmt.Fprintf(&builder, "📁 %s/\n", name)
			}
		}

		if err != nil {
			return "", nil, fmt.Errorf("не удалось просмотреть директорию: %w", err)
		}

		fmt.Fprintf(&builder, (strings.Repeat("─", 20)))
		fmt.Fprintf(&builder, "\nTotal items: %d\n", len(files))
		h.SendMessage(chatID, threadID, builder.String())

		return "", nil, nil
	}
	h.states.Store(chatID, &UserState{handler: step1})
}

func (h *MessageHandler) HandleUploadFile(update models.Update) {
	message := update.Message
	chatID := message.Chat.ID
	threadID := message.MessageThreadID

	// Проверяем, есть ли активная сессия для этого чата
	sessionI, ok := h.uploadSessions.Load(chatID)
	if !ok {
		// Нет сессии - начинаем новую (команда /upload)
		session := &UploadSession{
			Step:     1,
			ThreadID: threadID,
		}
		h.uploadSessions.Store(chatID, session)
		h.SendMessage(chatID, threadID, "📁 Введите путь директории для загрузки (например, /photos):")
		return
	}

	// Есть сессия - обрабатываем ввод
	session := sessionI.(*UploadSession)

	// Шаг 1: ожидание пути
	if session.Step == 1 {
		path := message.Text
		if path == "" {
			h.SendMessage(chatID, threadID, "Путь не может быть пустым. Попробуйте снова.")
			return
		}
		session.Path = path
		session.Step = 2
		session.LastFileTime = time.Now()
		h.uploadSessions.Store(chatID, session)
		h.SendMessage(chatID, threadID, "📤 Ожидаю файлы для загрузки. Таймер 60 секунд будет сбрасываться при каждом файле.\nОтправьте /cancel для отмены.")
		return
	}

	// Шаг 2: ожидание файлов
	// Проверка таймаута
	if time.Since(session.LastFileTime) > 60*time.Second {
		h.uploadSessions.Delete(chatID)
		h.SendMessage(chatID, threadID, "⏰ Время ожидания истекло. Загрузка завершена.")
		return
	}

	// Проверка отмены
	if message.Text == "cancel" || message.Text == "/cancel" || message.Text == "отмена" {
		h.uploadSessions.Delete(chatID)
		h.SendMessage(chatID, threadID, "❌ Загрузка отменена.")
		return
	}

	// Определяем, есть ли файл
	var fileID, fileName string
	if message.Document.FileID != "" {
		fileID = message.Document.FileID
		fileName = message.Document.FileName
	} else if len(message.Photo) > 0 {
		// Берём самое большое фото
		photo := message.Photo[len(message.Photo)-1]
		fileID = photo.FileID
		fileName = fmt.Sprintf("photo_%d.jpg", time.Now().Unix())
	} else {
		// Не файл – напоминаем
		h.SendMessage(chatID, threadID, "⏳ Ожидаю файлы. Чтобы отменить, отправьте /cancel")
		return
	}

	// Скачиваем файл
	fileData, err := h.downloadTelegramFile(fileID)
	if err != nil {
		h.SendMessage(chatID, threadID, fmt.Sprintf("❌ Ошибка скачивания файла: %v", err))
		return
	}

	// Обновляем время
	session.LastFileTime = time.Now()
	h.uploadSessions.Store(chatID, session)

	// Загружаем на Яндекс.Диск
	err = yandexapi.UploadFile(session.Path, fileName, fileData)
	if err != nil {
		h.SendMessage(chatID, threadID, fmt.Sprintf("❌ Ошибка загрузки на Яндекс.Диск: %v", err))
		return
	}

	h.SendMessage(chatID, threadID, fmt.Sprintf("✅ Файл «%s» загружен в %s", fileName, session.Path))
}

// downloadTelegramFile скачивает файл по fileID и возвращает его содержимое
func (h *MessageHandler) downloadTelegramFile(fileID string) ([]byte, error) {
	// Получаем путь к файлу на серверах Telegram
	filePath, err := h.getTelegramFilePath(fileID)
	if err != nil {
		return nil, err
	}

	// Формируем URL для скачивания
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", h.Token, filePath)

	// Скачиваем файл
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP ошибка: %s", resp.Status)
	}

	// Читаем всё содержимое
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	return data, nil
}

// getTelegramFilePath запрашивает путь к файлу по file_id
func (h *MessageHandler) getTelegramFilePath(fileID string) (string, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", h.Token, fileID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if !result.OK {
		return "", fmt.Errorf("telegram API error")
	}
	return result.Result.FilePath, nil
}
