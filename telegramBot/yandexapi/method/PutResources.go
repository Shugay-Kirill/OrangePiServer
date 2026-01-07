package method

import (
	"encoding/json"

	. "telegramBot/yandexapi/authenticated"
	. "telegramBot/yandexapi/init"
)

func PutResources(api *YandexDiskAPI, pathDirectory string, nameDirectory string) (map[string]interface{}, error) {

	params := map[string]string{
		"path": pathDirectory + "/" + nameDirectory,
	}

	url := BuildURL(api, "/resources", params)

	body, err := AuthenticatedRequest(api, "PUT", url, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
