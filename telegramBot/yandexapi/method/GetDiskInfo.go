package method

import (
	"fmt"

	"encoding/json"
	"telegramBot/yandexapi/authenticated"
)

type DiskInfo struct {
	TotalSpace int64 `json:"total_space"`
	UsedSpace  int64 `json:"used_space"`
	TrashSize  int64 `json:"trash_size"`
}

// // GetDiskInfo получает информацию о диске
func GetDiskInfo() (*DiskInfo, error) {
	fmt.Println("start GetDiskInfo")
	body, err := authenticated.AuthenticatedRequest("GET", "", nil, nil)
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
