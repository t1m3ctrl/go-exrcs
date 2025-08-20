package main

import (
	"fmt"
	"l2-18/config"
	"l2-18/internal/handler"
	"l2-18/internal/middleware"
	"l2-18/internal/service"
	"l2-18/internal/storage"
	"log"
	"net/http"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем слои приложения
	eventStorage := storage.NewInMemoryEventStorage()
	eventService := service.NewEventService(eventStorage)
	eventHandler := handler.NewEventHandler(eventService)

	// Настраиваем роуты
	mux := http.NewServeMux()

	// CRUD операции
	mux.HandleFunc("/create_event", eventHandler.CreateEvent)
	mux.HandleFunc("/update_event", eventHandler.UpdateEvent)
	mux.HandleFunc("/delete_event", eventHandler.DeleteEvent)
	mux.HandleFunc("/events_for_day", eventHandler.GetEventsForDay)
	mux.HandleFunc("/events_for_week", eventHandler.GetEventsForWeek)
	mux.HandleFunc("/events_for_month", eventHandler.GetEventsForMonth)

	// Применяем middleware
	handler := middleware.Logging(mux)

	// Запускаем сервер
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
