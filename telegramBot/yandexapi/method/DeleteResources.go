package method

import (
	"encoding/json"

	. "telegramBot/yandexapi/authenticated"
)

func DeleteResources(pathDirectory string, nameDirectory string) (map[string]interface{}, error) {

	params := map[string]string{
		"path":        pathDirectory + "/" + nameDirectory,
		"permanently": "false",
	}

	body, err := AuthenticatedRequest("DELETE", "/resources", params, nil)
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
