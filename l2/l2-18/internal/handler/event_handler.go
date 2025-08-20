package handler

import (
	"encoding/json"
	"fmt"
	"l2-18/internal/models"
	"l2-18/internal/service"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// EventHandler обработчик HTTP запросов для событий
type EventHandler struct {
	service *service.EventService
}

// NewEventHandler создает новый обработчик
func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// CreateEvent обработчик создания события
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := h.parseCreateEventRequest(r)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.service.CreateEvent(req)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccess(w, "event created successfully", event)
}

// UpdateEvent обработчик обновления события
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := h.parseUpdateEventRequest(r)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.service.UpdateEvent(req)
	if err != nil {
		if isBusinessLogicError(err) {
			h.sendError(w, err.Error(), http.StatusServiceUnavailable)
		} else {
			h.sendError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	h.sendSuccess(w, "event updated successfully", event)
}

// DeleteEvent обработчик удаления события
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := h.parseDeleteEventRequest(r)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteEvent(req)
	if err != nil {
		if isBusinessLogicError(err) {
			h.sendError(w, err.Error(), http.StatusServiceUnavailable)
		} else {
			h.sendError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	h.sendSuccess(w, "event deleted successfully", nil)
}

// GetEventsForDay обработчик получения событий на день
func (h *EventHandler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, date, err := h.parseGetEventsRequest(r)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.service.GetEventsForDay(userID, date)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccess(w, "events retrieved successfully", events)
}

// GetEventsForWeek обработчик получения событий на неделю
func (h *EventHandler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, date, err := h.parseGetEventsRequest(r)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.service.GetEventsForWeek(userID, date)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccess(w, "events retrieved successfully", events)
}

// GetEventsForMonth обработчик получения событий на месяц
func (h *EventHandler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, date, err := h.parseGetEventsRequest(r)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.service.GetEventsForMonth(userID, date)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccess(w, "events retrieved successfully", events)
}

// parseCreateEventRequest парсит запрос на создание события
func (h *EventHandler) parseCreateEventRequest(r *http.Request) (*models.CreateEventRequest, error) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		var req models.CreateEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	}

	// Парсим как form data
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		return nil, err
	}

	return &models.CreateEventRequest{
		UserID:      userID,
		Date:        r.FormValue("date"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}, nil
}

// parseUpdateEventRequest парсит запрос на обновление события
func (h *EventHandler) parseUpdateEventRequest(r *http.Request) (*models.UpdateEventRequest, error) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		var req models.UpdateEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	}

	// Парсим как form data
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		return nil, err
	}

	return &models.UpdateEventRequest{
		ID:          id,
		UserID:      userID,
		Date:        r.FormValue("date"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}, nil
}

// parseDeleteEventRequest парсит запрос на удаление события
func (h *EventHandler) parseDeleteEventRequest(r *http.Request) (*models.DeleteEventRequest, error) {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		var req models.DeleteEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return &req, nil
	}

	// Парсим как form data
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		return nil, err
	}

	return &models.DeleteEventRequest{
		ID:     id,
		UserID: userID,
	}, nil
}

// parseGetEventsRequest парсит параметры для получения событий
func (h *EventHandler) parseGetEventsRequest(r *http.Request) (int, time.Time, error) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return 0, time.Time{}, err
	}

	userID, err := strconv.Atoi(values.Get("user_id"))
	if err != nil {
		return 0, time.Time{}, err
	}

	dateStr := values.Get("date")
	if dateStr == "" {
		return 0, time.Time{}, fmt.Errorf("date parameter is required")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, time.Time{}, err
	}

	return userID, date, nil
}

// sendSuccess отправляет успешный ответ
func (h *EventHandler) sendSuccess(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.APIResponse{
		Result: message,
		Data:   data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// sendError отправляет ответ с ошибкой
func (h *EventHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.APIResponse{
		Error: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// isBusinessLogicError определяет, является ли ошибка бизнес-логической
func isBusinessLogicError(err error) bool {
	errMsg := err.Error()
	return errMsg == "event with ID not found" ||
		errMsg == "event does not belong to user" ||
		errMsg == "failed to delete event: event with ID not found" ||
		errMsg == "failed to update event: event with ID not found"
}
