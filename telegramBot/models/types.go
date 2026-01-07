package models

import (
// "fmt"
// "net/http"
)

// Структуры для Telegram API
type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID       int         `json:"message_id"`
	From            User        `json:"from"`
	Chat            Chat        `json:"chat"`
	Text            string      `json:"text"`
	MessageThreadID int         `json:"message_thread_id"`
	Photo           []PhotoSize `json:"photo"`
	Document        Document    `json:"document"`
	Caption         string      `json:"caption"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size"`
}

type Document struct {
	FileID       string    `json:"file_id"`
	FileUniqueID string    `json:"file_unique_id"`
	Thumbnail    PhotoSize `json:"thumb"`
	FileName     string    `json:"file_name"`
	MimeType     string    `json:"mime_type"`
	FileSize     int       `json:"file_size"`
}

// Структура для APIYandexDisk
// type YandexDiskAPI struct {
// 	hostNameURL string
// 	token       string
// 	client      *http.Client
// }

// type DiskInfo struct {
// 	TotalSpace int64 `json:"total_space"`
// 	UsedSpace  int64 `json:"used_space"`
// 	TrashSize  int64 `json:"trash_size"`
// }
