package models

import (
	"fmt"
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

// Структуры для внешнего API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type FileAnalysisRequest struct {
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
	FileSize int    `json:"file_size"`
}

type FileAnalysisResponse struct {
	IsJPG      bool   `json:"is_jpg"`
	FileType   string `json:"file_type"`
	Dimensions string `json:"dimensions,omitempty"`
	FileSizeKB string `json:"file_size_kb"`
}

// Конвертер для API
func (d *Document) ToFileAnalysisRequest() FileAnalysisRequest {
	return FileAnalysisRequest{
		FileName: d.FileName,
		MimeType: d.MimeType,
		FileSize: d.FileSize,
	}
}

func (p *PhotoSize) ToFileAnalysisResponse() FileAnalysisResponse {
	return FileAnalysisResponse{
		IsJPG:      true,
		FileType:   "image/jpeg",
		Dimensions: formatDimensions(p.Width, p.Height),
		FileSizeKB: formatFileSize(p.FileSize),
	}
}

func formatDimensions(width, height int) string {
	return fmt.Sprintf("%dx%d", width, height)
}

func formatFileSize(size int) string {
	return fmt.Sprintf("%.2f KB", float64(size)/1024)
}
