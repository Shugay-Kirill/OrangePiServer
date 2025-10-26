package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Bot struct {
	token  string
	config *Config
}

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"` // Изменено на указатель
}

type Message struct {
	MessageID       int         `json:"message_id"`
	From            User        `json:"from"`
	Chat            Chat        `json:"chat"`
	Text            string      `json:"text"`
	MessageThreadID int         `json:"message_thread_id"`
	Photo           []PhotoSize `json:"photo"`
	Document        Document    `json:"document"`
	Caption         string      `json:"caption"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size"`
}

type Document struct {
	FileID       string    `json:"file_id"`
	FileUniqueID string    `json:"file_unique_id"`
	Thumbnail    PhotoSize `json:"thumb"`
	FileName     string    `json:"file_name"`
	MimeType     string    `json:"mime_type"`
	FileSize     int       `json:"file_size"`
}

func NewBot(config *Config) *Bot {
	return &Bot{
		token:  config.TelegramToken,
		config: config,
	}
}

func (b *Bot) startPolling() {
	log.Println("🚀 Бот запущен с прямым polling...")
	log.Printf("📏 Максимальная длина вывода API: %d символов", b.config.MaxLengthAPIOutput)
	log.Println("📝 Ожидаю сообщения в топиках...")
	offset := 0

	for {
		updates, err := b.getUpdates(offset)
		if err != nil {
			log.Printf("❌ Ошибка получения updates: %v", err)
			continue
		}

		for _, update := range updates {
			b.handleUpdate(update)
			offset = update.UpdateID + 1
		}
	}
}

func (b *Bot) getUpdates(offset int) ([]Update, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=60", b.token, offset)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Логируем сырой ответ для отладки с использованием MAX_LENGTH_MESSEGE_API
	if len(body) > 0 {
		maxLength := b.config.MaxLengthAPIOutput
		output := string(body)
		if len(output) > maxLength {
			output = output[:maxLength] + "..."
		}
		log.Printf("📨 Получен ответ от API (%d/%d символов): %s", len(body), maxLength, output)
	}

	var response struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return nil, err
	}

	if !response.OK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	return response.Result, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (b *Bot) handleUpdate(update Update) {
	// Проверяем, что Message не nil
	if update.Message == nil {
		return
	}

	message := update.Message

	log.Printf("📩 Получено сообщение:")
	log.Printf("   👤 От: %s (@%s)", message.From.FirstName, message.From.Username) // Исправлено: message.From.Username
	log.Printf("   🆔 Chat ID: %d", message.Chat.ID)
	log.Printf("   🏷️ Thread ID: %d", message.MessageThreadID)
	log.Printf("   📊 Тип чата: %s", message.Chat.Type)
	if message.Chat.Title != "" {
		log.Printf("   🏷️ Название чата: %s", message.Chat.Title)
	}

	if len(message.Photo) > 0 {
		log.Printf("   📸 Фото: %d вариантов размера", len(message.Photo))
		b.handlePhoto(update)
		return
	}

	if message.Document.FileID != "" {
		log.Printf("   📎 Документ: %s", message.Document.FileName)
		b.handleDocument(update)
		return
	}

	if message.Text == "" {
		log.Printf("   💬 Текст: (пустое сообщение или другой тип)")
		b.handleOtherMessage(update)
		return
	}

	log.Printf("   💬 Текст: %s", message.Text)

	if message.Text == "/start" {
		b.handleStart(update)
		return
	}

	if message.Text == "/help" {
		b.handleHelp(update)
		return
	}

	if message.Text == "/features" {
		b.handleFeatures(update)
		return
	}

	if message.Text == "/info" {
		b.handleInfo(update)
		return
	}

	b.handleRegularMessage(update)
}

func (b *Bot) isJPGImage(document Document) bool {
	if strings.HasPrefix(document.MimeType, "image/jpeg") {
		return true
	}

	fileName := strings.ToLower(document.FileName)
	if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg") {
		return true
	}

	if document.MimeType == "image/jpg" {
		return true
	}

	return false
}

func (b *Bot) getLargestPhoto(photos []PhotoSize) PhotoSize {
	if len(photos) == 0 {
		return PhotoSize{}
	}

	largest := photos[0]
	for _, photo := range photos {
		if photo.FileSize > largest.FileSize {
			largest = photo
		}
	}
	return largest
}

func (b *Bot) handlePhoto(update Update) {
	message := update.Message
	chatID := message.Chat.ID
	threadID := message.MessageThreadID

	largestPhoto := b.getLargestPhoto(message.Photo)

	response := fmt.Sprintf(`📸 <b>Получено фото!</b>

🖼️ <b>Информация о фото:</b>
• 📏 Размер: <b>%d×%d</b> пикселей
• 💾 Вес: <b>%.2f KB</b>
• 🆔 File ID: <code>%s</code>

📝 <b>Подпись:</b> %s

✅ <b>Статус:</b> Это JPG изображение (Telegram конвертирует все фото в JPG)

🎯 <i>Фото успешно обработано!</i>`,
		largestPhoto.Width,
		largestPhoto.Height,
		float64(largestPhoto.FileSize)/1024,
		largestPhoto.FileID[:min(20, len(largestPhoto.FileID))]+"...",
		b.getCaptionText(message.Caption),
	)

	if err := b.sendMessage(chatID, threadID, response); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (b *Bot) handleDocument(update Update) {
	message := update.Message
	chatID := message.Chat.ID
	threadID := message.MessageThreadID

	document := message.Document
	isJPG := b.isJPGImage(document)

	var status string
	if isJPG {
		status = "✅ <b>Это JPG изображение!</b>"
	} else {
		status = "❌ <b>Это не JPG изображение</b>"
	}

	response := fmt.Sprintf(`📎 <b>Получен документ!</b>

📋 <b>Информация о файле:</b>
• 📝 Имя: <code>%s</code>
• 🏷️ MIME Type: <b>%s</b>
• 💾 Размер: <b>%.2f KB</b>
• 🆔 File ID: <code>%s</code>

📝 <b>Подпись:</b> %s

%s

🎯 <i>Документ проверен на соответствие формату JPG!</i>`,
		document.FileName,
		document.MimeType,
		float64(document.FileSize)/1024,
		document.FileID[:min(20, len(document.FileID))]+"...",
		b.getCaptionText(message.Caption),
		status,
	)

	if err := b.sendMessage(chatID, threadID, response); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (b *Bot) handleOtherMessage(update Update) {
	message := update.Message
	chatID := message.Chat.ID
	threadID := message.MessageThreadID

	response := `🔮 <b>Получено сообщение другого типа!</b>

📊 <b>Информация:</b>
• Тип: Не текстовое сообщение
• Может содержать: фото, документ, стикер, голосовое и т.д.

💡 <b>Что я умею проверять:</b>
• 📸 Фотографии (автоматически определяю как JPG)
• 📎 Документы (проверяю формат JPG)
• 💬 Текстовые сообщения

🎯 <i>Используйте /help для получения списка команд</i>`

	if err := b.sendMessage(chatID, threadID, response); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (b *Bot) getCaptionText(caption string) string {
	if caption == "" {
		return "<i>нет подписи</i>"
	}
	return caption
}

func (b *Bot) handleStart(update Update) {
	message := update.Message
	chatID := message.Chat.ID
	threadID := message.MessageThreadID

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

📸 <b>Проверка фото:</b>
Отправьте мне:
• Фото - я определю его параметры
• Документ JPG - я проверю формат
• Любой файл - я скажу, JPG ли это

📊 <b>Информация о текущем сообщении:</b>
• 👤 Ваше имя: <b>%s</b>
• 🆔 Ваш ID: <code>%d</code>
• 💬 ID чата: <code>%d</code>
• 🏷️ ID топика: <code>%d</code>

💡 <b>Просто отправьте мне фото или документ для проверки!</b>`,
		message.From.FirstName,
		b.config.MaxLengthAPIOutput,
		message.From.FirstName,
		message.From.ID,
		chatID,
		threadID,
	)

	if err := b.sendMessage(chatID, threadID, response); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (b *Bot) handleHelp(update Update) {
	message := update.Message
	chatID := message.Chat.ID
	threadID := message.MessageThreadID

	response := `🆘 <b>Помощь по боту</b>

📚 <b>Доступные команды:</b>
• /start - начать работу с ботом
• /help - показать эту справку
• /features - возможности бота
• /info - информация о текущем чате

🔧 <b>Как использовать:</b>
1. Просто отправьте мне любое сообщение
2. Я отвечу в том же разделе/топике
3. Используйте команды для получения информации

📸 <b>Проверка файлов:</b>
• Отправьте фото - увидите его параметры
• Отправьте документ - проверю формат JPG
• Все файлы анализируются автоматически

⚙️ <b>Настройки:</b>
• Максимальная длина вывода API: <b>%d символов</b>
• Можно изменить через переменную MAX_LENGTH_MESSEGE_API

❓ <b>Частые вопросы:</b>
• Бот не отвечает? Проверьте, что он добавлен в группу
• Сообщения не в том топике? Обновите библиотеку бота
• Есть вопросы? Напишите разработчику

💡 <b>Совет:</b> Используйте /features чтобы узнать о всех возможностях!`

	response = fmt.Sprintf(response, b.config.MaxLengthAPIOutput)

	if err := b.sendMessage(chatID, threadID, response); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (b *Bot) handleFeatures(update Update) {
	message := update.Message
	chatID := message.Chat.ID
	threadID := message.MessageThreadID

	response := `🚀 <b>Возможности бота</b>

🎯 <b>Основные функции:</b>
• <b>Умные ответы</b> - Анализирую ваши сообщения и отвечаю соответствующим образом
• <b>Работа с топиками</b> - Отвечаю в том же разделе, где было отправлено сообщение
• <b>Поддержка групп</b> - Работаю в супергруппах, группах и личных сообщениях
• <b>Форматирование</b> - Поддерживаю HTML разметку для красивого отображения

📸 <b>Функции проверки файлов:</b>
• <b>Автоопределение фото</b> - Анализирую размер, вес и параметры изображений
• <b>Проверка JPG</b> - Определяю, является ли документ JPG изображением
• <b>MIME type анализ</b> - Проверяю тип файла по MIME и расширению

🔧 <b>Технические возможности:</b>
• <b>Детальное логирование</b> - Записываю всю информацию о входящих сообщениях
• <b>Настраиваемый вывод API</b> - Макс. длина логов: <b>%d символов</b>
• <b>Обработка ошибок</b> - Грамотно обрабатываю ошибки и уведомляю о них
• <b>Поддержка Docker</b> - Могу работать в контейнерах Docker
• <b>Конфигурация через .env</b> - Легко настраиваюсь через переменные окружения

📊 <b>Информационные функции:</b>
• Показываю ID чата и топика
• Отображаю информацию о пользователе
• Поддерживаю различные типы чатов
• Предоставляю детальную статистику

⚙️ <b>Настройки конфигурации:</b>
• TELEGRAM_BOT_TOKEN - Токен бота
• DEBUG - Режим отладки
• MAX_LENGTH_MESSEGE_API - Макс. длина вывода API (текущее значение: %d)

🛠️ <b>В разработке:</b>
• Интеграция с базами данных
• Планировщик задач
• Система плагинов
• Webhook поддержка

💡 <b>Напишите любое сообщение, чтобы протестировать мои возможности!</b>`

	response = fmt.Sprintf(response, b.config.MaxLengthAPIOutput, b.config.MaxLengthAPIOutput)

	if err := b.sendMessage(chatID, threadID, response); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (b *Bot) handleInfo(update Update) {
	message := update.Message
	chatID := message.Chat.ID
	threadID := message.MessageThreadID

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
	if threadID != 0 {
		topicStatus = fmt.Sprintf("✅ Да (ID: %d)", threadID)
	}

	response := fmt.Sprintf(`ℹ️ <b>Информация о чате</b>

📋 <b>Основная информация:</b>
• 💬 Тип чата: <b>%s</b>
• 🏷️ Название: <b>%s</b>
• 🆔 ID чата: <code>%d</code>
• 🏷️ Топик: %s

👤 <b>Информация о пользователе:</b>
• Имя: <b>%s</b>
• Username: @%s
• ID пользователя: <code>%d</code>

🔧 <b>Техническая информация:</b>
• Бот: @%s
• Поддержка топиков: ✅ Включена
• Макс. длина API логов: <b>%d символов</b>
• Режим отладки: ✅ Включен

💡 <b>Примечание:</b>
Этот бот специально разработан для работы с топиками в Telegram группах и всегда отвечает в том же разделе, откуда пришло сообщение.`,
		chatType,
		b.getChatTitle(message.Chat),
		chatID,
		topicStatus,
		message.From.FirstName,
		message.From.Username, // Исправлено: message.From.Username
		message.From.ID,
		b.getBotUsername(),
		b.config.MaxLengthAPIOutput,
	)

	if err := b.sendMessage(chatID, threadID, response); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (b *Bot) handleRegularMessage(update Update) {
	message := update.Message
	chatID := message.Chat.ID
	threadID := message.MessageThreadID

	response := fmt.Sprintf(`✅ <b>Сообщение получено!</b>

📝 <b>Ваше сообщение:</b>
<code>%s</code>

👤 <b>От:</b> <b>%s</b> (@%s)

📊 <b>Техническая информация:</b>
• 💬 Чат ID: <code>%d</code>
• 🏷️ Топик ID: <code>%d</code>
• 📏 Макс. длина API: <b>%d символов</b>

📸 <b>Попробуйте отправить:</b>
• Фотографию - я покажу её параметры
• Документ JPG - я проверю формат
• Команду /features - все возможности

🎯 <i>Этот ответ отправлен в тот же топик!</i>`,
		message.Text,
		message.From.FirstName,
		message.From.Username, // Исправлено: message.From.Username
		chatID,
		threadID,
		b.config.MaxLengthAPIOutput,
	)

	if err := b.sendMessage(chatID, threadID, response); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (b *Bot) getChatTitle(chat Chat) string {
	if chat.Title != "" {
		return chat.Title
	}
	return "Без названия"
}

func (b *Bot) getBotUsername() string {
	testURL := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", b.token)
	resp, err := http.Get(testURL)
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "unknown"
	}

	if result["ok"].(bool) {
		botInfo := result["result"].(map[string]interface{})
		return botInfo["username"].(string)
	}
	return "unknown"
}

func (b *Bot) sendMessage(chatID int64, threadID int, text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)

	params := url.Values{}
	params.Add("chat_id", strconv.FormatInt(chatID, 10))
	params.Add("text", text)
	params.Add("parse_mode", "HTML")

	if threadID != 0 {
		params.Add("message_thread_id", strconv.Itoa(threadID))
		log.Printf("📤 Отправка сообщения в топик %d", threadID)
	} else {
		log.Printf("📤 Отправка сообщения в основной чат")
	}

	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Ошибка API: %s - %s", resp.Status, string(body))
		return fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	log.Printf("✅ Сообщение успешно отправлено!")
	return nil
}

func main() {
	log.Println("🔧 Загрузка конфигурации...")
	config := LoadConfig()

	if config.TelegramToken == "" {
		log.Fatal("❌ TELEGRAM_BOT_TOKEN не установлен. Добавьте его в .env файл")
	}

	log.Println("🤖 Инициализация бота...")
	bot := NewBot(config)

	log.Println("🔌 Проверка подключения к Telegram API...")
	testURL := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", config.TelegramToken)
	resp, err := http.Get(testURL)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("❌ Ошибка парсинга ответа: %v", err)
	}

	if result["ok"].(bool) {
		botInfo := result["result"].(map[string]interface{})
		log.Printf("✅ Бот @%s готов к работе!", botInfo["username"])
	} else {
		log.Fatal("❌ Неверный токен бота")
	}

	log.Println("✨ Бот теперь имеет следующие команды:")
	log.Println("   /start - показать приветствие и возможности")
	log.Println("   /help - помощь по использованию")
	log.Println("   /features - все функции бота")
	log.Println("   /info - информация о чате")
	log.Println("📸 Функции проверки файлов:")
	log.Println("   - Автоматическое определение фотографий")
	log.Println("   - Проверка документов на формат JPG")
	log.Println("   - Анализ параметров файлов")
	log.Printf("⚙️  Настройки конфигурации:")
	log.Printf("   - MAX_LENGTH_MESSEGE_API: %d символов", config.MaxLengthAPIOutput)
	log.Printf("   - DEBUG: %v", config.Debug)

	bot.startPolling()
}
