package models

import "time"

// Event представляет событие в календаре
type Event struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateEventRequest структура для создания события
type CreateEventRequest struct {
	UserID      int    `json:"user_id"`
	Date        string `json:"date"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateEventRequest структура для обновления события
type UpdateEventRequest struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Date        string `json:"date"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// DeleteEventRequest структура для удаления события
type DeleteEventRequest struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
}

// APIResponse стандартный ответ API
type APIResponse struct {
	Result string      `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}
