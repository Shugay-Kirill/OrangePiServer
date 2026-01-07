package yandexinit

import (
	"fmt"
	"net/http"
	"time"

	"telegramBot/config"
)

type YandexDiskAPI struct {
	HostNameURL string
	Token       string
	Client      *http.Client
}

func NewYandexDiskAPI() *YandexDiskAPI {
	config := config.LoadConfig()

	return &YandexDiskAPI{
		HostNameURL: config.UrlYandexDisk,
		Token:       config.YandexDiskToken,
		Client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: false,
			},
		},
	}
}

func InitYandexDisk() {

	fmt.Println("🚀 Yandex.Disk API Client \n")
	// fmt.Println()

	// Создаем API клиент (конфигурация загружается автоматически)
	api := NewYandexDiskAPI()

	fmt.Printf("api", api)

}
