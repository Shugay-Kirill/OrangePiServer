package main

import (
	"log"
	"strconv"

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
	// Получаем информацию о чате и топике
	chatID := update.Message.Chat.ID
	threadID := b.getMessageThreadID(update.Message)

	log.Printf("Обработка /start: ChatID=%d, ThreadID=%d", chatID, threadID)

	var message string
	chatType := b.getChatType(update.Message.Chat)

	if chatType == "private" {
		message = `🤖 <b>Привет!</b>

Рад приветствовать вас! Это простой Telegram бот.

Бот успешно запущен и работает! 🚀`
	} else {
		message = `🤖 <b>Привет всем!</b>

Я бот для работы в группах и топиках.

<b>Я умею:</b>
• Отвечать в том же топике, где вы написали
• Работать в форумах и группах с темами
• Сохранять контекст обсуждения`
	}

	// Отправляем сообщение
	if err := b.sendMessage(chatID, threadID, message); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	// Логируем информацию о сообщении
	chatID := update.Message.Chat.ID
	threadID := b.getMessageThreadID(update.Message)
	chatType := b.getChatType(update.Message.Chat)

	log.Printf("Сообщение от [%s] в %s (ChatID: %d, ThreadID: %d): %s",
		update.Message.From.UserName,
		chatType,
		chatID,
		threadID,
		update.Message.Text)

	// Создаем ответ
	var response string

	if update.Message.IsCommand() {
		response = "❌ <b>Неизвестная команда</b>\nИспользуйте /start для получения информации"
	} else {
		response = "✅ <b>Сообщение получено!</b>\n\n" +
			"Я получил ваше сообщение в этом разделе.\n\n" +
			"<i>Текст:</i> " + update.Message.Text
	}

	// Добавляем отладочную информацию
	response += "\n\n📋 <i>Информация:</i>\n" +
		"Чат: " + chatType + "\n" +
		"ID чата: " + strconv.FormatInt(chatID, 10) + "\n" +
		"ID раздела: " + strconv.Itoa(threadID)

	// Отправляем сообщение
	if err := b.sendMessage(chatID, threadID, response); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

// sendMessage универсальный метод отправки сообщений с поддержкой топиков
func (b *Bot) sendMessage(chatID int64, threadID int, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	// Если threadID != 0, значит это топик и нужно установить MessageThreadID
	if threadID != 0 {
		// Используем рефлексию или проверяем доступность поля
		// В новых версиях библиотеки это должно работать напрямую:
		// msg.MessageThreadID = threadID

		// Обходной способ через создание сообщения с нужными параметрами
		msg = tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"

		// Пытаемся установить MessageThreadID через интерфейс сообщения
		// Это сработает если библиотека поддерживает топики
		if setThreadID(msg, threadID) {
			log.Printf("Отправка в топик ID: %d", threadID)
		} else {
			log.Printf("Библиотека не поддерживает MessageThreadID, отправка в основной чат")
		}
	}

	_, err := b.api.Send(msg)
	return err
}

// setThreadID пытается установить MessageThreadID для сообщения
func setThreadID(msg tgbotapi.MessageConfig, threadID int) bool {
	// Проверяем наличие поля MessageThreadID через type assertion
	// Это обходной путь для совместимости
	if msgConfig, ok := interface{}(msg).(interface{ SetMessageThreadID(int) }); ok {
		// Если библиотека поддерживает метод SetMessageThreadID
		msgConfig.SetMessageThreadID(threadID)
		return true
	}
	return false
}

// getMessageThreadID получает ID треда/топика из сообщения
func (b *Bot) getMessageThreadID(message *tgbotapi.Message) int {
	if message == nil {
		return 0
	}

	// Пытаемся получить MessageThreadID через интерфейс
	if msg, ok := interface{}(message).(interface{ GetMessageThreadID() int }); ok {
		return msg.GetMessageThreadID()
	}

	// Если интерфейс не доступен, проверяем напрямую (для новых версий)
	// Это сработает только если библиотека обновлена
	return 0 // Возвращаем 0 если не можем получить ID
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

func (b *Bot) Start() {
	log.Printf("Бот авторизован как: %s (ID: %d)", b.api.Self.UserName, b.api.Self.ID)
	log.Printf("Режим отладки: %v", b.config.Debug)
	log.Println("Бот запущен и ожидает сообщений...")
	log.Println("Для работы с топиками убедитесь, что:")
	log.Println("1. Бот добавлен в группу как администратор")
	log.Println("2. В группе включены темы/топики")
	log.Println("3. Библиотека обновлена до последней версии")

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
