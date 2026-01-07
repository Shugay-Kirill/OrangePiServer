package method

import (
	"encoding/json"
	. "telegramBot/yandexapi/authenticated"
	. "telegramBot/yandexapi/init"
)

type DiskInfo struct {
	TotalSpace int64 `json:"total_space"`
	UsedSpace  int64 `json:"used_space"`
	TrashSize  int64 `json:"trash_size"`
}

// // GetDiskInfo получает информацию о диске
func GetDiskInfo(api *YandexDiskAPI) (*DiskInfo, error) {
	url := BuildURL(api, "", nil)
	body, err := AuthenticatedRequest(api, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var diskInfo DiskInfo
	err = json.Unmarshal(body, &diskInfo)
	if err != nil {
		return nil, err
	}

	return &diskInfo, nil
}
