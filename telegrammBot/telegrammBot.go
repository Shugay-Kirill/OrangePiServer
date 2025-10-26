package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Bot struct {
	token string
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
	// Прямое извлечение message_thread_id из JSON
	MessageThreadID int `json:"message_thread_id"`
	// Поля для работы с фото
	Photo    []PhotoSize `json:"photo"`
	Document Document    `json:"document"`
	Caption  string      `json:"caption"`
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
		token: config.TelegramToken,
	}
}

func (b *Bot) startPolling() {
	log.Println("🚀 Бот запущен с прямым polling...")
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

	// Логируем сырой ответ для отладки
	maxLengthMessegeAPI := 5000
	if len(body) > 0 {
		log.Printf("📨 Получен ответ от API: %s", string(body)[:min(maxLengthMessegeAPI, len(body))])
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
	if update.Message.Text == "" {
		return
	}

	// Логируем всю информацию о сообщении
	log.Printf("📩 Получено сообщение:")
	log.Printf("   👤 От: %s (@%s)", update.Message.From.FirstName, update.Message.From.Username)
	log.Printf("   💬 Текст: %s", update.Message.Text)
	log.Printf("   🆔 Chat ID: %d", update.Message.Chat.ID)
	log.Printf("   🏷️ Thread ID: %d", update.Message.MessageThreadID)
	log.Printf("   📊 Тип чата: %s", update.Message.Chat.Type)
	if update.Message.Chat.Title != "" {
		log.Printf("   🏷️ Название чата: %s", update.Message.Chat.Title)
	}

	// Обрабатываем команды
	if update.Message.Text == "/start" {
		b.handleStart(update)
		return
	}

	// Обрабатываем обычные сообщения
	if update.Message.Text == "/infoMessege" {
		b.handleRegularMessage(update)
		return
	}

}

func (b *Bot) handleStart(update Update) {
	chatID := update.Message.Chat.ID
	threadID := update.Message.MessageThreadID

	message := fmt.Sprintf(`🤖 <b>Бот запущен!</b>

Привет, <b>%s</b>! 🎉

<b>Информация о сообщении:</b>
• 💬 Чат: <code>%d</code>
• 🏷️ Топик: <code>%d</code>
• 👤 Ваш ID: <code>%d</code>

✅ <i>Этот ответ отправлен в тот же топик!</i>`,
		update.Message.From.FirstName,
		chatID,
		threadID,
		update.Message.From.ID,
	)

	if err := b.sendMessage(chatID, threadID, message); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (b *Bot) handleRegularMessage(update Update) {
	chatID := update.Message.Chat.ID
	threadID := update.Message.MessageThreadID

	message := fmt.Sprintf(`✅ <b>Сообщение получено!</b>

<b>Ваше сообщение:</b>
<code>%s</code>

<b>Детали:</b>
• 👤 От: <b>%s</b> (@%s)
• 💬 Чат ID: <code>%d</code>
• 🏷️ Топик ID: <code>%d</code>
• 📊 Тип чата: %s

🎯 <i>Этот ответ отправлен в тот же топик!</i>`,
		update.Message.Text,
		update.Message.From.FirstName,
		update.Message.From.Username,
		chatID,
		threadID,
		update.Message.Chat.Type,
	)

	if err := b.sendMessage(chatID, threadID, message); err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
	}
}

func (b *Bot) sendMessage(chatID int64, threadID int, text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)

	params := url.Values{}
	params.Add("chat_id", strconv.FormatInt(chatID, 10))
	params.Add("text", text)
	params.Add("parse_mode", "HTML")

	// Ключевой момент: передаем message_thread_id если он не 0
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

	// Читаем и логируем ответ от API
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Ошибка API: %s - %s", resp.Status, string(body))
		return fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	log.Printf("✅ Сообщение успешно отправлено!")
	log.Printf("   💬 Чат: %d", chatID)
	log.Printf("   🏷️ Топик: %d", threadID)

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

	// Тестируем подключение
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

	bot.startPolling()
}
