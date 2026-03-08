package main

import (
	"telegramBot/handlersTelegramBot"
	"telegramBot/yandexapi/initYD"
)

func main() {
	initYD.InitYandexDisk()
	handlersTelegramBot.StartTelegramBot()
}
