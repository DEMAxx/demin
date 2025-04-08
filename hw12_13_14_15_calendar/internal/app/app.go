package app

import (
	"context"
	"fmt"
	"time"

	memorystorage "github.com/DEMAxx/demin/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/google/uuid"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event memorystorage.Event) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(
	ctx context.Context,
	id uuid.UUID,
	title string,
	date time.Time,
	description string,
	user uuid.UUID,
) error {
	event := memorystorage.Event{
		ID:          id,
		Title:       title,
		Date:        date,
		Duration:    time.Minute * 5,
		Description: description,
		User:        user,
		Notify:      date.Add(-time.Hour),
	}
	if err := a.storage.CreateEvent(ctx, event); err != nil {
		a.logger.Error("failed to create event: " + err.Error())
		return err
	}
	a.logger.Info(fmt.Sprintf("event created: %q", id))
	return nil
}
