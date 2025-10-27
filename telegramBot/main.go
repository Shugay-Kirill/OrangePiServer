package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"telegramBot/config"
	"telegramBot/handlersTelegramBot"
	"telegramBot/models"
)

type Bot struct {
	token   string
	config  *config.Config
	handler *handlersTelegramBot.MessageHandler
}

func NewBot(config *config.Config) *Bot {
	handler := handlersTelegramBot.NewMessageHandler(config.TelegramToken, config)
	return &Bot{
		token:   config.TelegramToken,
		config:  config,
		handler: handler,
	}
}

func (b *Bot) startPolling() {
	log.Println("🚀 Бот запущен с прямым polling...")
	log.Printf("📏 Максимальная длина вывода API: %d символов", b.config.MaxLengthAPIOutput)

	offset := 0
	for {
		updates, err := b.getUpdates(offset)
		if err != nil {
			log.Printf("❌ Ошибка получения updates: %v", err)
			continue
		}

		for _, update := range updates {
			b.handler.HandleUpdate(update)
			offset = update.UpdateID + 1
		}
	}
}

func (b *Bot) getUpdates(offset int) ([]models.Update, error) {
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
	if len(body) > 0 {
		maxLength := b.config.MaxLengthAPIOutput
		output := string(body)
		if len(output) > maxLength {
			output = output[:maxLength] + "..."
		}
		log.Printf("📨 Получен ответ от API (%d/%d символов): %s", len(body), maxLength, output)
	}

	var response struct {
		OK     bool            `json:"ok"`
		Result []models.Update `json:"result"`
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

func main() {
	log.Println("🔧 Загрузка конфигурации...")
	config := config.LoadConfig()

	if config.TelegramToken == "" {
		log.Fatal("❌ TELEGRAM_BOT_TOKEN не установлен")
	}

	log.Println("🤖 Инициализация бота...")
	bot := NewBot(config)

	// Проверка подключения
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

	log.Println("✨ Бот запущен!")
	bot.startPolling()
}
