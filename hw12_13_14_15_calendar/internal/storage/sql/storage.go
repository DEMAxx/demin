package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

type Event struct {
	ID          uuid.UUID
	Title       string
	Date        time.Time
	Duration    time.Duration
	Description string
	User        uuid.UUID
	Notify      time.Time
}

type Storage struct {
	db  *sql.DB
	dsn string
}

func New(dsn string) (*Storage, error) {
	return &Storage{dsn: dsn}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sql.Open("postgres", s.dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	s.db = db
	return nil
}

func (s *Storage) GetDB() *sql.DB {
	return s.db
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event Event) error {
	query := `INSERT INTO events (id, title, date, duration, description, user_id, notify) 
VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Title,
		event.Date,
		event.Duration,
		event.Description,
		event.User,
		event.Notify,
	)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (Event, error) {
	query := `SELECT id, title, date, duration, description, user_id, notify FROM events WHERE id = $1`
	row := s.db.QueryRowContext(ctx, query, id)

	var event Event
	if err := row.Scan(&event.ID, &event.Title, &event.Date); err != nil {
		if err == sql.ErrNoRows {
			return Event{}, fmt.Errorf("event with ID %s not found", id)
		}
		return Event{}, fmt.Errorf("failed to get event: %w", err)
	}

	return event, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event Event) error {
	query := `UPDATE events SET 
                  title = $2, date = $3, duration = $4, description = $5, user_id = $6, notify = $7 
              WHERE id = $1`
	_, err := s.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Title,
		event.Date,
		event.Duration,
		event.Description,
		event.User,
		event.Notify,
	)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	return nil
}
