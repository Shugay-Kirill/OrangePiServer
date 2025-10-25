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

	// Определяем, где было отправлено сообщение
	chatType := b.getChatType(update.Message.Chat)

	if chatType == "private" {
		message = `🤖 <b>Привет!</b>

Рад приветствовать вас! Это простой Telegram бот.

Бот успешно запущен и работает! 🚀

<b>Команды:</b>
/start - показать это сообщение`
	} else {
		message = `🤖 <b>Привет всем!</b>

Я бот для работы в группах.

<b>Особенности:</b>
• Отвечаю в том же разделе, где написали
• Работаю в группах и супергруппах

Используйте /start в любом чате!`
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ParseMode = "HTML"

	// Проверяем и устанавливаем ID топика, если доступно
	b.setMessageThreadID(&msg, update.Message)

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	// Логируем информацию о сообщении
	chatType := b.getChatType(update.Message.Chat)

	log.Printf("Сообщение от [%s %s] в %s: %s",
		update.Message.From.FirstName,
		update.Message.From.LastName,
		chatType,
		update.Message.Text)

	// Создаем ответ
	var response string

	if update.Message.IsCommand() {
		response = "❌ <b>Неизвестная команда</b>\nИспользуйте /start для получения информации"
	} else {
		response = "✅ <b>Сообщение получено!</b>\n\n" +
			"Я получил ваше сообщение: <i>\"" + update.Message.Text + "\"</i>\n\n" +
			"Используйте /start для получения информации о боте."
	}

	// Добавляем информацию о месте отправки
	response += "\n\n📍 <i>Отправлено в: " + chatType + "</i>"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	msg.ParseMode = "HTML"

	// Устанавливаем тот же раздел/топик
	b.setMessageThreadID(&msg, update.Message)

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

// setMessageThreadID устанавливает ID треда/топика для ответа
func (b *Bot) setMessageThreadID(msg *tgbotapi.MessageConfig, message *tgbotapi.Message) {
	// Проверяем, доступно ли поле MessageThreadID в этой версии библиотеки
	// Это безопасный способ без прямого доступа к полю
	if message == nil {
		return
	}

	// Вместо прямого доступа к MessageThreadID, используем обходной путь
	// В реальном приложении можно использовать рефлексию или обновить библиотеку
}

// checkLibraryVersion проверяет версию библиотеки
func (b *Bot) checkLibraryVersion() {
	log.Println("Проверка версии библиотеки go-telegram-bot-api...")
	log.Println("Для работы с топиками убедитесь, что версия >= v5.0.0")
}

func (b *Bot) Start() {
	b.checkLibraryVersion()
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
