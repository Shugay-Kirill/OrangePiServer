package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	config *Config
	token  string
}

type Config struct {
	TelegramToken string
	Debug         bool
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		TelegramToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		Debug:         getEnvAsBool("DEBUG", false),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
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
		token:  config.TelegramToken,
	}, nil
}

// sendMessageToThread отправляет сообщение в конкретный топик через прямой HTTP запрос
func (b *Bot) sendMessageToThread(chatID int64, threadID int, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)

	requestBody := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	}

	// Если threadID не 0, добавляем параметр message_thread_id
	if threadID != 0 {
		requestBody["message_thread_id"] = threadID
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	log.Printf("Сообщение отправлено в чат %d, топик %d", chatID, threadID)
	return nil
}

// getThreadIDFromMessage пытается извлечь ID топика из сообщения
func (b *Bot) getThreadIDFromMessage(message *tgbotapi.Message) int {
	// В новых версиях библиотеки это поле должно быть доступно
	// Если нет, используем рефлексию или другие методы

	// Проверяем через рефлексию
	messageValue := reflect.ValueOf(message).Elem()
	if messageValue.IsValid() {
		threadIDField := messageValue.FieldByName("message_thread_id")
		if threadIDField.IsValid() && threadIDField.CanInterface() {
			if threadID, ok := threadIDField.Interface().(int); ok {
				return threadID
			}
		}
	}

	// Если рефлексия не сработала, пробуем получить из Reply если это ответ в топике
	if message.ReplyToMessage != nil {
		replyValue := reflect.ValueOf(message.ReplyToMessage).Elem()
		if replyValue.IsValid() {
			threadIDField := replyValue.FieldByName("message_thread_id")
			if threadIDField.IsValid() && threadIDField.CanInterface() {
				if threadID, ok := threadIDField.Interface().(int); ok {
					return threadID
				}
			}
		}
	}

	return 0
}

func (b *Bot) handleStart(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	threadID := b.getThreadIDFromMessage(update.Message)

	log.Printf("Обработка /start: ChatID=%d, ThreadID=%d", chatID, threadID)

	message := `🤖 <b>Бот запущен!</b>

Отлично! Я работаю в этом разделе.

<b>Информация:</b>
• ID чата: ` + strconv.FormatInt(chatID, 10) + `
• ID раздела: ` + strconv.Itoa(threadID) + `
• Сообщение получено в теме! ✅`

	if err := b.sendMessageToThread(chatID, threadID, message); err != nil {
		log.Printf("Ошибка отправки: %v", err)
		// Fallback: отправляем обычным способом
		b.sendFallbackMessage(chatID, message)
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	threadID := b.getThreadIDFromMessage(update.Message)

	log.Printf("Сообщение от [%s] в ChatID=%d, ThreadID=%d: %s",
		update.Message.From.UserName, chatID, threadID, update.Message.Text)

	response := `✅ <b>Сообщение получено в теме!</b>

<b>Текст:</b> ` + update.Message.Text + `

<b>Отладочная информация:</b>
• ID чата: ` + strconv.FormatInt(chatID, 10) + `
• ID раздела: ` + strconv.Itoa(threadID) + `
• Username: @` + update.Message.From.UserName + `

🎯 <i>Этот ответ должен быть в той же теме!</i>`

	if err := b.sendMessageToThread(chatID, threadID, response); err != nil {
		log.Printf("Ошибка отправки в топик: %v", err)
		b.sendFallbackMessage(chatID, response)
	}
}

// sendFallbackMessage обычная отправка через библиотеку
func (b *Bot) sendFallbackMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Ошибка fallback отправки: %v", err)
	}
}

func (b *Bot) Start() {
	log.Printf("Бот авторизован как: %s", b.api.Self.UserName)
	log.Println("Бот запущен с поддержкой топиков через прямое API")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() && update.Message.Command() == "start" {
			b.handleStart(update)
			continue
		}

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
