package storage

import (
	"fmt"
	"l2-18/internal/models"
	"sync"
	"time"
)

// EventStorage интерфейс для работы с событиями
type EventStorage interface {
	Create(event *models.Event) error
	Update(event *models.Event) error
	Delete(id, userID int) error
	GetByDateRange(userID int, start, end time.Time) ([]*models.Event, error)
	GetByID(id, userID int) (*models.Event, error)
}

// InMemoryEventStorage реализация хранилища в памяти
type InMemoryEventStorage struct {
	events   map[int]*models.Event
	nextID   int
	userToID map[int][]int // userID -> []eventIDs
	mu       sync.RWMutex
}

// NewInMemoryEventStorage создает новое хранилище в памяти
func NewInMemoryEventStorage() *InMemoryEventStorage {
	return &InMemoryEventStorage{
		events:   make(map[int]*models.Event),
		nextID:   1,
		userToID: make(map[int][]int),
	}
}

// Create создает новое событие
func (s *InMemoryEventStorage) Create(event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = s.nextID
	s.nextID++
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	s.events[event.ID] = event
	s.userToID[event.UserID] = append(s.userToID[event.UserID], event.ID)

	return nil
}

// Update обновляет существующее событие
func (s *InMemoryEventStorage) Update(event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.events[event.ID]
	if !exists {
		return fmt.Errorf("event with ID %d not found", event.ID)
	}

	if existing.UserID != event.UserID {
		return fmt.Errorf("event does not belong to user")
	}

	event.CreatedAt = existing.CreatedAt
	event.UpdatedAt = time.Now()
	s.events[event.ID] = event

	return nil
}

// Delete удаляет событие
func (s *InMemoryEventStorage) Delete(id, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, exists := s.events[id]
	if !exists {
		return fmt.Errorf("event with ID %d not found", id)
	}

	if event.UserID != userID {
		return fmt.Errorf("event does not belong to user")
	}

	delete(s.events, id)

	// Удаляем из индекса пользователя
	userEvents := s.userToID[userID]
	for i, eventID := range userEvents {
		if eventID == id {
			s.userToID[userID] = append(userEvents[:i], userEvents[i+1:]...)
			break
		}
	}

	return nil
}

// GetByDateRange возвращает события в указанном диапазоне дат
func (s *InMemoryEventStorage) GetByDateRange(userID int, start, end time.Time) ([]*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*models.Event
	eventIDs := s.userToID[userID]

	for _, eventID := range eventIDs {
		event := s.events[eventID]
		if event != nil && isDateInRange(event.Date, start, end) {
			result = append(result, event)
		}
	}

	return result, nil
}

// GetByID возвращает событие по ID
func (s *InMemoryEventStorage) GetByID(id, userID int) (*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[id]
	if !exists {
		return nil, fmt.Errorf("event with ID %d not found", id)
	}

	if event.UserID != userID {
		return nil, fmt.Errorf("event does not belong to user")
	}

	return event, nil
}

// isDateInRange проверяет, попадает ли дата в заданный диапазон
func isDateInRange(date, start, end time.Time) bool {
	// Сравниваем только даты, игнорируя время
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	startOnly := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	endOnly := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)

	return (dateOnly.Equal(startOnly) || dateOnly.After(startOnly)) &&
		(dateOnly.Equal(endOnly) || dateOnly.Before(endOnly))
}
