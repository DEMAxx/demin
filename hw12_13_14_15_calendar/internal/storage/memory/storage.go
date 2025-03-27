package memorystorage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
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
	mu     sync.RWMutex //nolint:unused
	events map[uuid.UUID]Event
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]Event),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if _, exists := s.events[event.ID]; exists {
		return fmt.Errorf("event with ID %s already exists", event.ID)
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) GetEvent(id uuid.UUID) (Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[id]
	if !exists {
		return Event{}, fmt.Errorf("event with ID %s not found", id)
	}

	return event, nil
}

func (s *Storage) DeleteEvent(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return fmt.Errorf("event with ID %s not found", id)
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) UpdateEvent(event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; !exists {
		return fmt.Errorf("event with ID %s not found", event.ID)
	}

	s.events[event.ID] = event
	return nil
}
