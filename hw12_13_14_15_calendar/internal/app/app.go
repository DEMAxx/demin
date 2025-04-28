package app

import (
	"context"
	"fmt"
	"time"

	memorystorage "github.com/DEMAxx/demin/hw12_13_14_15_calendar/internal/storage/memory"
	databasestorage "github.com/DEMAxx/demin/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/google/uuid"
)

type App struct {
	logger    Logger
	storage   Storage
	dbStorage DBStorage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event memorystorage.Event) (uuid.UUID, error)
	CreateUser(ctx context.Context, user memorystorage.User) (uuid.UUID, error)
}

type DBStorage interface {
	CreateEvent(ctx context.Context, event databasestorage.Event) error
	GetEvent(ctx context.Context, id uuid.UUID) (databasestorage.Event, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	UpdateEvent(ctx context.Context, event databasestorage.Event) error
}

func New(logger Logger, storage Storage, dbStorage DBStorage) *App {
	return &App{
		logger:    logger,
		storage:   storage,
		dbStorage: dbStorage,
	}
}

func (a *App) CreateUser(
	ctx context.Context,
	name string,
) (uuid.UUID, error) {
	user := memorystorage.User{
		Name: name,
	}
	if uid, err := a.storage.CreateUser(ctx, user); err != nil {
		a.logger.Error("failed to create user: " + err.Error())
		return uuid.Nil, err
	} else {
		a.logger.Info(fmt.Sprintf("user created: %s", uid))
		return uid, nil
	}

}

func (a *App) CreateEvent(
	ctx context.Context,
	id uuid.UUID,
	title string,
	date time.Time,
	description string,
	user uuid.UUID,
) (uuid.UUID, error) {
	event := memorystorage.Event{
		ID:          id,
		Title:       title,
		Date:        date,
		Duration:    time.Minute * 5,
		Description: description,
		User:        user,
		Notify:      date.Add(-time.Hour),
	}
	if uid, err := a.storage.CreateEvent(ctx, event); err != nil {
		a.logger.Error("failed to create event: " + err.Error())
		return uuid.Nil, err
	} else {
		a.logger.Info(fmt.Sprintf("event created: %q", id))
		return uid, nil
	}
}
