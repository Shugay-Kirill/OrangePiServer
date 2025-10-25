package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	commands map[string]func(update tgbotapi.Update)
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		api:      api,
		commands: make(map[string]func(update tgbotapi.Update)),
	}

	bot.registerCommands()
	return bot, nil
}

func (b *Bot) registerCommands() {
	// Регистрация команд
	b.commands["/start"] = b.handleStart
	b.commands["/help"] = b.handleHelp
	b.commands["/echo"] = b.handleEcho
	b.commands["/calc"] = b.handleCalc
}

func (b *Bot) handleStart(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"🤖 Добро пожаловать!\n\n"+
			"Я пример бота на Go.\n"+
			"Доступные команды:\n"+
			"/start - начать работу\n"+
			"/help - помощь\n"+
			"/echo [текст] - эхо\n"+
			"/calc [число] [оператор] [число] - калькулятор\n\n"+
			"Просто напиши мне что-нибудь!")
	b.sendMessage(msg)
}

func (b *Bot) handleHelp(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"📖 Помощь по боту:\n\n"+
			"• /start - начать работу\n"+
			"• /help - показать эту справку\n"+
			"• /echo [текст] - повторить текст\n"+
			"• /calc [число] [+-*/] [число] - простой калькулятор\n\n"+
			"Примеры:\n"+
			"/echo Привет мир!\n"+
			"/calc 5 + 3")
	b.sendMessage(msg)
}

func (b *Bot) handleEcho(update tgbotapi.Update) {
	text := update.Message.Text
	args := strings.TrimSpace(strings.TrimPrefix(text, "/echo"))

	if args == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📝 Использование: /echo [текст]")
		b.sendMessage(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🔊 "+args)
	b.sendMessage(msg)
}

func (b *Bot) handleCalc(update tgbotapi.Update) {
	text := update.Message.Text
	args := strings.Fields(strings.TrimPrefix(text, "/calc"))

	if len(args) != 3 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"🧮 Использование: /calc [число] [оператор] [число]\n"+
				"Пример: /calc 5 + 3")
		b.sendMessage(msg)
		return
	}

	a, err1 := strconv.ParseFloat(args[0], 64)
	operator := args[1]
	bNum, err2 := strconv.ParseFloat(args[2], 64)

	if err1 != nil || err2 != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Ошибка: неверный формат чисел")
		b.sendMessage(msg)
		return
	}

	var result float64
	var errorMsg string

	switch operator {
	case "+":
		result = a + bNum
	case "-":
		result = a - bNum
	case "*":
		result = a * bNum
	case "/":
		if bNum == 0 {
			errorMsg = "❌ Ошибка: деление на ноль"
		} else {
			result = a / bNum
		}
	default:
		errorMsg = "❌ Ошибка: неверный оператор. Используйте +, -, *, /"
	}

	if errorMsg != "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, errorMsg)
		b.sendMessage(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"🧮 Результат: "+args[0]+" "+operator+" "+args[2]+" = "+strconv.FormatFloat(result, 'f', -1, 64))
	b.sendMessage(msg)
}

func (b *Bot) handleUnknownCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"❓ Неизвестная команда. Используйте /help для списка команд.")
	b.sendMessage(msg)
}

func (b *Bot) handleTextMessage(update tgbotapi.Update) {
	response := "📨 Вы написали: " + update.Message.Text
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	b.sendMessage(msg)
}

func (b *Bot) sendMessage(msg tgbotapi.MessageConfig) {
	msg.ParseMode = "HTML"
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

func (b *Bot) Start() {
	log.Printf("Авторизован как %s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Обработка команд
		if update.Message.IsCommand() {
			command := update.Message.Command()
			if handler, exists := b.commands["/"+command]; exists {
				handler(update)
			} else {
				b.handleUnknownCommand(update)
			}
			continue
		}

		// Обработка обычных текстовых сообщений
		b.handleTextMessage(update)
	}
}

func main() {
	// Получение токена из переменной окружения
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не установлен")
	}

	bot, err := NewBot(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Start()
}
