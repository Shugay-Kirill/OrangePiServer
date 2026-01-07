package method

import (
	"fmt"
	"net/url"

	. "telegramBot/yandexapi/init"
)

// type Yandexapi struct{
// 	yaapi *yandexinit.YandexDiskAPI
// }

// BuildURL создает URL с параметрами
func BuildURL(api *YandexDiskAPI, pathNameURl string, search map[string]string) string {
	uri := api.HostNameURL + pathNameURl
	if len(search) > 0 {
		query := url.Values{}
		for key, value := range search {
			fmt.Printf("\nkey = %s, value = %s\n", key, value)
			query.Add(key, value)
		}
		fmt.Printf("\nURL = %s\n", uri)
		uri += "?" + query.Encode()
	}
	fmt.Printf("\nURL = %s\n", uri)
	return uri
}
