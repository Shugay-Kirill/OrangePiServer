package authenticated

import (
	"fmt"
	"net/url"

	"telegramBot/yandexapi/initYD"
)

// BuildURL создает URL с параметрами
func BuildURL(pathNameURl string, parametr map[string]string) string {
	apiAuth := initYD.GetYandexDiskAPI()
	uriRequest := apiAuth.HostNameURL + pathNameURl
	if len(parametr) > 0 {
		query := url.Values{}
		for key, value := range parametr {
			fmt.Printf("\nkey = %s, value = %s\n", key, value)
			query.Add(key, value)
		}
		fmt.Printf("\nURL = %s\n", uriRequest)
		uriRequest += "?" + query.Encode()
	}
	fmt.Printf("\nURL = %s\n", uriRequest)
	return uriRequest
}
