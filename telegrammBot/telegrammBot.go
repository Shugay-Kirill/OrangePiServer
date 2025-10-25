package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	config *Config
}

func NewBot(config *Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		return nil, err
	}

	api.Debug = config.Debug

	return &Bot{
		api:    api,
		config: config,
	}, nil
}

func (b *Bot) handleStart(update tgbotapi.Update) {
	var message string
	user := update.Message.From

	// Определяем, где было отправлено сообщение
	chatType := b.getChatType(update.Message.Chat)
	topicInfo := b.getTopicInfo(update.Message)

	if chatType == "private" {
		message = `🤖 <b>Привет, %s!</b>

Рад приветствовать вас! Это простой Telegram бот.

Бот успешно запущен и работает! 🚀

<b>Команды:</b>
/start - показать это сообщение`
	} else {
		message = `🤖 <b>Привет всем!</b>

Я бот для работы в группах и темах.

<b>Особенности:</b>
• Отвечаю в том же топике, где написали
• Понимаю контекст обсуждения
• Работаю в группах и супергруппах

Используйте /start в любом топике!`
	}

	// Добавляем информацию о топике если есть
	if topicInfo != "" {
		message += "\n\n" + topicInfo
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ParseMode = "HTML"

	// Устанавливаем ID топика, если сообщение из топика
	if update.Message.MessageThreadID != 0 {
		msg.MessageThreadID = update.Message.MessageThreadID
	}

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	// Логируем информацию о сообщении
	chatType := b.getChatType(update.Message.Chat)
	topicInfo := b.getTopicInfo(update.Message)

	log.Printf("Сообщение от [%s %s] в %s%s: %s",
		update.Message.From.FirstName,
		update.Message.From.LastName,
		chatType,
		topicInfo,
		update.Message.Text)

	// Создаем ответ с учетом топика
	var response string

	if update.Message.IsCommand() {
		response = "❌ <b>Неизвестная команда</b>\nИспользуйте /start для получения информации"
	} else {
		response = "✅ <b>Сообщение получено!</b>\n\n" +
			"Я получил ваше сообщение в этом топике: <i>\"" + update.Message.Text + "\"</i>\n\n" +
			"Используйте /start для получения информации о боте."
	}

	// Добавляем информацию о месте отправки
	response += "\n\n📍 <i>Отправлено в: " + chatType
	if topicInfo != "" {
		response += " • " + topicInfo
	}
	response += "</i>"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	msg.ParseMode = "HTML"

	// Ключевой момент: указываем тот же MessageThreadID
	if update.Message.MessageThreadID != 0 {
		msg.MessageThreadID = update.Message.MessageThreadID
		log.Printf("Отвечаю в топик ID: %d", update.Message.MessageThreadID)
	}

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

// getChatType определяет тип чата
func (b *Bot) getChatType(chat *tgbotapi.Chat) string {
	switch {
	case chat.IsPrivate():
		return "личные сообщения"
	case chat.IsGroup():
		return "группа"
	case chat.IsSuperGroup():
		return "супергруппа"
	case chat.IsChannel():
		return "канал"
	default:
		return "неизвестный чат"
	}
}

// getTopicInfo возвращает информацию о топике
func (b *Bot) getTopicInfo(message *tgbotapi.Message) string {
	if message.MessageThreadID == 0 {
		return "" // Не топик
	}

	// Если это корневое сообщение топика (TopicCreated)
	if message.ForumTopicCreated != nil {
		return "топик: " + message.ForumTopicCreated.Name
	}

	// Если это обычное сообщение в топике
	if message.MessageThreadID != 0 {
		// Пытаемся получить информацию о топике
		// В реальном приложении здесь можно кешировать названия топиков
		return "топик ID: " + string(rune(message.MessageThreadID))
	}

	return ""
}

// handleTopicCreated обрабатывает создание нового топика
func (b *Bot) handleTopicCreated(update tgbotapi.Update) {
	if update.Message == nil || update.Message.ForumTopicCreated == nil {
		return
	}

	topicName := update.Message.ForumTopicCreated.Name
	log.Printf("Создан новый топик: %s", topicName)

	response := "🎉 <b>Новый топик создан!</b>\n\n" +
		"Название: <i>" + topicName + "</i>\n\n" +
		"Приветствую всех участников! Я буду отвечать в этом топике."

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	msg.ParseMode = "HTML"
	msg.MessageThreadID = update.Message.MessageThreadID

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Ошибка отправки приветствия в топик: %v", err)
	}
}

func (b *Bot) Start() {
	log.Printf("Бот авторизован как: %s (ID: %d)", b.api.Self.UserName, b.api.Self.ID)
	log.Printf("Режим отладки: %v", b.config.Debug)
	log.Println("Бот запущен и ожидает сообщений...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Обрабатываем создание топика
		if update.Message.ForumTopicCreated != nil {
			b.handleTopicCreated(update)
			continue
		}

		// Обработка команды /start
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			b.handleStart(update)
			continue
		}

		// Обработка всех остальных сообщений
		b.handleMessage(update)
	}
}

func main() {
	log.Println("Загрузка конфигурации...")
	config := LoadConfig()

	log.Println("Инициализация бота...")
	bot, err := NewBot(config)
	if err != nil {
		log.Fatal("Ошибка инициализации бота:", err)
	}

	bot.Start()
}
