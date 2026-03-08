package initYD

import (
	"fmt"
	"net/http"
	"time"

	"telegramBot/config"
)

// YandexDiskAPI - структура клиента
type YandexDiskAuth struct {
	HostNameURL string
	Token       string
	Client      *http.Client
}

// Глобальная переменная (неэкспортируемая)
var YandexDiskAPI *YandexDiskAuth

// GetYandexDiskAPI - возвращает глобальный экземпляр
func GetYandexDiskAPI() *YandexDiskAuth {
	if YandexDiskAPI == nil {
		panic("YandexDiskAPI не инициализирован. Сначала вызовите InitYandexDisk()")
	}
	return YandexDiskAPI
}

// NewYandexDiskAPI - создает новый экземпляр (приватный)
func newYandexDiskAPI() *YandexDiskAuth {
	config := config.LoadConfig()

	return &YandexDiskAuth{
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

// InitYandexDisk - инициализация глобального экземпляра (вызывается один раз)
func InitYandexDisk() {
	fmt.Println("🚀 Инициализация Yandex.Disk API...")

	// Создаем и сохраняем в глобальную переменную
	YandexDiskAPI = newYandexDiskAPI()

	fmt.Printf("✅ Yandex.Disk API инициализирован. Token: %v\n",
		YandexDiskAPI.Token != "")
}
