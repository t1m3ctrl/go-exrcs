package service

import (
	"l2-18/internal/models"
	"l2-18/internal/storage"
	"testing"
	"time"
)

func TestEventService_CreateEvent(t *testing.T) {
	strg := storage.NewInMemoryEventStorage()
	service := NewEventService(strg)

	tests := []struct {
		name    string
		req     *models.CreateEventRequest
		wantErr bool
	}{
		{
			name: "valid event",
			req: &models.CreateEventRequest{
				UserID:      1,
				Date:        "2023-12-31",
				Title:       "New Year Party",
				Description: "Celebrate new year",
			},
			wantErr: false,
		},
		{
			name: "invalid user ID",
			req: &models.CreateEventRequest{
				UserID:      0,
				Date:        "2023-12-31",
				Title:       "New Year Party",
				Description: "Celebrate new year",
			},
			wantErr: true,
		},
		{
			name: "empty title",
			req: &models.CreateEventRequest{
				UserID:      1,
				Date:        "2023-12-31",
				Title:       "",
				Description: "Celebrate new year",
			},
			wantErr: true,
		},
		{
			name: "invalid date format",
			req: &models.CreateEventRequest{
				UserID:      1,
				Date:        "31-12-2023",
				Title:       "New Year Party",
				Description: "Celebrate new year",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := service.CreateEvent(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if event == nil {
					t.Error("CreateEvent() returned nil event")
					return
				}
				if event.ID == 0 {
					t.Error("CreateEvent() event ID not set")
				}
				if event.CreatedAt.IsZero() {
					t.Error("CreateEvent() CreatedAt not set")
				}
			}
		})
	}
}

func TestEventService_UpdateEvent(t *testing.T) {
	strg := storage.NewInMemoryEventStorage()
	service := NewEventService(strg)

	// Создаем событие для обновления
	createReq := &models.CreateEventRequest{
		UserID:      1,
		Date:        "2023-12-31",
		Title:       "Original Title",
		Description: "Original Description",
	}
	event, err := service.CreateEvent(createReq)
	if err != nil {
		t.Fatal("Failed to create event for test:", err)
	}

	tests := []struct {
		name    string
		req     *models.UpdateEventRequest
		wantErr bool
	}{
		{
			name: "valid update",
			req: &models.UpdateEventRequest{
				ID:          event.ID,
				UserID:      1,
				Date:        "2024-01-01",
				Title:       "Updated Title",
				Description: "Updated Description",
			},
			wantErr: false,
		},
		{
			name: "non-existent event",
			req: &models.UpdateEventRequest{
				ID:          9999,
				UserID:      1,
				Date:        "2024-01-01",
				Title:       "Updated Title",
				Description: "Updated Description",
			},
			wantErr: true,
		},
		{
			name: "wrong user",
			req: &models.UpdateEventRequest{
				ID:          event.ID,
				UserID:      2,
				Date:        "2024-01-01",
				Title:       "Updated Title",
				Description: "Updated Description",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedEvent, err := service.UpdateEvent(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if updatedEvent == nil {
					t.Error("UpdateEvent() returned nil event")
					return
				}
				if updatedEvent.Title != tt.req.Title {
					t.Errorf("UpdateEvent() title = %v, want %v", updatedEvent.Title, tt.req.Title)
				}
			}
		})
	}
}

func TestEventService_DeleteEvent(t *testing.T) {
	strg := storage.NewInMemoryEventStorage()
	service := NewEventService(strg)

	// Создаем событие для удаления
	createReq := &models.CreateEventRequest{
		UserID:      1,
		Date:        "2023-12-31",
		Title:       "Event to Delete",
		Description: "This will be deleted",
	}
	event, err := service.CreateEvent(createReq)
	if err != nil {
		t.Fatal("Failed to create event for test:", err)
	}

	tests := []struct {
		name    string
		req     *models.DeleteEventRequest
		wantErr bool
	}{
		{
			name: "valid deletion",
			req: &models.DeleteEventRequest{
				ID:     event.ID,
				UserID: 1,
			},
			wantErr: false,
		},
		{
			name: "non-existent event",
			req: &models.DeleteEventRequest{
				ID:     9999,
				UserID: 1,
			},
			wantErr: true,
		},
		{
			name: "wrong user",
			req: &models.DeleteEventRequest{
				ID:     event.ID,
				UserID: 2,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteEvent(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventService_GetEventsForDay(t *testing.T) {
	strg := storage.NewInMemoryEventStorage()
	service := NewEventService(strg)

	// Создаем несколько событий
	events := []struct {
		userID int
		date   string
		title  string
	}{
		{1, "2023-12-31", "Event 1"},
		{1, "2023-12-31", "Event 2"},
		{1, "2024-01-01", "Event 3"},
		{2, "2023-12-31", "Event 4"},
	}

	for _, e := range events {
		_, err := service.CreateEvent(&models.CreateEventRequest{
			UserID: e.userID,
			Date:   e.date,
			Title:  e.title,
		})
		if err != nil {
			t.Fatal("Failed to create event:", err)
		}
	}

	testDate, _ := time.Parse("2006-01-02", "2023-12-31")

	// Тест получения событий пользователя 1 на 31.12.2023
	dayEvents, err := service.GetEventsForDay(1, testDate)
	if err != nil {
		t.Error("GetEventsForDay() error:", err)
	}

	if len(dayEvents) != 2 {
		t.Errorf("GetEventsForDay() got %d events, want 2", len(dayEvents))
	}

	// Тест для несуществующего пользователя
	_, err = service.GetEventsForDay(0, testDate)
	if err == nil {
		t.Error("GetEventsForDay() should return error for invalid user ID")
	}
}

func TestEventService_GetEventsForWeek(t *testing.T) {
	strg := storage.NewInMemoryEventStorage()
	service := NewEventService(strg)

	// Создаем события на разные дни недели
	// 2023-12-31 - воскресенье (конец недели)
	// 2023-12-25 - понедельник (начало той же недели)
	events := []struct {
		date  string
		title string
	}{
		{"2023-12-25", "Monday Event"},    // Понедельник
		{"2023-12-27", "Wednesday Event"}, // Среда
		{"2023-12-31", "Sunday Event"},    // Воскресенье
		{"2024-01-01", "Next Week"},       // Следующая неделя
	}

	for _, e := range events {
		_, err := service.CreateEvent(&models.CreateEventRequest{
			UserID: 1,
			Date:   e.date,
			Title:  e.title,
		})
		if err != nil {
			t.Fatal("Failed to create event:", err)
		}
	}

	// Запрашиваем события на неделю, содержащую 31.12.2023
	testDate, _ := time.Parse("2006-01-02", "2023-12-31")
	weekEvents, err := service.GetEventsForWeek(1, testDate)
	if err != nil {
		t.Error("GetEventsForWeek() error:", err)
	}

	// Должны получить 3 события (25, 27, 31 декабря)
	if len(weekEvents) != 3 {
		t.Errorf("GetEventsForWeek() got %d events, want 3", len(weekEvents))
	}
}

func TestEventService_GetEventsForMonth(t *testing.T) {
	strg := storage.NewInMemoryEventStorage()
	service := NewEventService(strg)

	// Создаем события в декабре и январе
	events := []struct {
		date  string
		title string
	}{
		{"2023-12-01", "December Start"},
		{"2023-12-15", "December Middle"},
		{"2023-12-31", "December End"},
		{"2024-01-01", "January Start"},
	}

	for _, e := range events {
		_, err := service.CreateEvent(&models.CreateEventRequest{
			UserID: 1,
			Date:   e.date,
			Title:  e.title,
		})
		if err != nil {
			t.Fatal("Failed to create event:", err)
		}
	}

	// Запрашиваем события декабря 2023
	testDate, _ := time.Parse("2006-01-02", "2023-12-15")
	monthEvents, err := service.GetEventsForMonth(1, testDate)
	if err != nil {
		t.Error("GetEventsForMonth() error:", err)
	}

	// Должны получить 3 события декабря
	if len(monthEvents) != 3 {
		t.Errorf("GetEventsForMonth() got %d events, want 3", len(monthEvents))
	}
}
