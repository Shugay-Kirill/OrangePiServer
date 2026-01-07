package handlersTelegramBot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"telegramBot/config"
	"telegramBot/models"
)

type MessageHandler struct {
	Token  string
	Config *config.Config
}

func NewMessageHandler(token string, config *config.Config) *MessageHandler {
	return &MessageHandler{
		Token:  token,
		Config: config,
	}
}

func (h *MessageHandler) HandleUpdate(update models.Update) {
	if update.Message == nil {
		return
	}

	message := update.Message

	log.Printf("📩 Получено сообщение:")
	log.Printf("   	👤 От: %s (@%s)", message.From.FirstName, message.From.Username)
	log.Printf("   	🆔 Chat ID: %d", message.Chat.ID)
	log.Printf("   	🏷️ Thread ID: %d", message.MessageThreadID)
	log.Printf("   	📊 Тип чата: %s", message.Chat.Type)

	if message.Chat.Title != "" {
		log.Printf("   	🏷️ Название чата: %s", message.Chat.Title)
	}

	// Определяем тип сообщения и передаем соответствующему обработчику
	switch {
	case len(message.Photo) > 0:
		log.Printf("   	📸 Фото: %d вариантов размера", len(message.Photo))
		// h.HandlePhoto(update)
	case message.Document.FileID != "":
		log.Printf("   	📎 Документ: %s", message.Document.FileName)
		// h.HandleDocument(update)
	case message.Text == "":
		log.Printf("   	💬 Текст: (пустое сообщение или другой тип)")
		h.HandleOtherMessage(update)
	case message.MessageThreadID == 29:
		log.Printf("   		💬 Это чат Наши фотографии")
	default:
		log.Printf("   	💬 Текст: %s", message.Text)
		h.HandleTextMessage(update)
	}
}

func (h *MessageHandler) HandleTextMessage(update models.Update) {
	message := update.Message

	switch message.Text {
	case "/start":
		h.HandleStartCommand(update)
	case "/help":
		h.HandleHelpCommand(update)
	case "/features":
		h.HandleFeaturesCommand(update)
	case "/info":
		h.HandleInfoCommand(update)
	case "/infoMessage":
		h.HandleRegularMessage(update)
	default:
	}
}

func (h *MessageHandler) HandleRegularMessage(update models.Update) {
	message := update.Message
	response := fmt.Sprintf(`✅ <b>Сообщение получено!</b>

📝 <b>Ваше сообщение:</b>
<code>%s</code>

👤 <b>От:</b> <b>%s</b> (@%s)

📊 <b>Техническая информация:</b>
• 💬 Чат ID: <code>%d</code>
• 🏷️ Топик ID: <code>%d</code>
• 📏 Макс. длина API: <b>%d символов</b>

🎯 <i>Этот ответ отправлен в тот же топик!</i>`,
		message.Text,
		message.From.FirstName,
		message.From.Username,
		message.Chat.ID,
		message.MessageThreadID,
		h.Config.MaxLengthAPIOutput,
	)

	h.SendMessage(message.Chat.ID, message.MessageThreadID, response)
}

func (h *MessageHandler) SendMessage(chatID int64, threadID int, text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", h.Token)

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

// Вспомогательные методы
func (h *MessageHandler) getCaptionText(caption string) string {
	if caption == "" {
		return "<i>нет подписи</i>"
	}
	return caption
}

func (h *MessageHandler) getChatTitle(chat models.Chat) string {
	if chat.Title != "" {
		return chat.Title
	}
	return "Без названия"
}
