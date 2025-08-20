package service

import (
	"fmt"
	"l2-18/internal/models"
	"l2-18/internal/storage"
	"strings"
	"time"
)

// EventService содержит бизнес-логику для работы с событиями
type EventService struct {
	storage storage.EventStorage
}

// NewEventService создает новый сервис событий
func NewEventService(storage storage.EventStorage) *EventService {
	return &EventService{storage: storage}
}

// CreateEvent создает новое событие
func (s *EventService) CreateEvent(req *models.CreateEventRequest) (*models.Event, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	event := &models.Event{
		UserID:      req.UserID,
		Date:        date,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
	}

	if err := s.storage.Create(event); err != nil {
		return nil, fmt.Errorf("failed to create event: %v", err)
	}

	return event, nil
}

// UpdateEvent обновляет существующее событие
func (s *EventService) UpdateEvent(req *models.UpdateEventRequest) (*models.Event, error) {
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	event := &models.Event{
		ID:          req.ID,
		UserID:      req.UserID,
		Date:        date,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
	}

	if err := s.storage.Update(event); err != nil {
		return nil, fmt.Errorf("failed to update event: %v", err)
	}

	return event, nil
}

// DeleteEvent удаляет событие
func (s *EventService) DeleteEvent(req *models.DeleteEventRequest) error {
	if req.ID <= 0 {
		return fmt.Errorf("invalid event ID")
	}
	if req.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	if err := s.storage.Delete(req.ID, req.UserID); err != nil {
		return fmt.Errorf("failed to delete event: %v", err)
	}

	return nil
}

// GetEventsForDay возвращает события на день
func (s *EventService) GetEventsForDay(userID int, date time.Time) ([]*models.Event, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start

	return s.storage.GetByDateRange(userID, start, end)
}

// GetEventsForWeek возвращает события на неделю
func (s *EventService) GetEventsForWeek(userID int, date time.Time) ([]*models.Event, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Находим начало недели (понедельник)
	weekday := date.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	start := date.AddDate(0, 0, -int(weekday-1))
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	// Конец недели (воскресенье)
	end := start.AddDate(0, 0, 6)

	return s.storage.GetByDateRange(userID, start, end)
}

// GetEventsForMonth возвращает события на месяц
func (s *EventService) GetEventsForMonth(userID int, date time.Time) ([]*models.Event, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Начало месяца
	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Конец месяца
	end := start.AddDate(0, 1, -1)

	return s.storage.GetByDateRange(userID, start, end)
}

// validateCreateRequest валидирует запрос на создание события
func (s *EventService) validateCreateRequest(req *models.CreateEventRequest) error {
	if req.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if req.Date == "" {
		return fmt.Errorf("date is required")
	}
	return nil
}

// validateUpdateRequest валидирует запрос на обновление события
func (s *EventService) validateUpdateRequest(req *models.UpdateEventRequest) error {
	if req.ID <= 0 {
		return fmt.Errorf("invalid event ID")
	}
	if req.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if req.Date == "" {
		return fmt.Errorf("date is required")
	}
	return nil
}
