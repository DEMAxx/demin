package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	storage := New()

	t.Run("CreateEvent", func(t *testing.T) {
		event := Event{
			ID:          uuid.New(),
			Title:       "Test Event",
			Date:        time.Now(),
			Duration:    time.Hour,
			Description: "This is a test event",
			User:        uuid.New(),
			Notify:      time.Now().Add(-time.Minute),
		}

		ctx := context.Background()

		err := storage.CreateEvent(ctx, event)
		require.NoError(t, err)

		storedEvent, err := storage.GetEvent(event.ID)
		require.NoError(t, err)
		require.Equal(t, event, storedEvent)
	})

	t.Run("GetEvent", func(t *testing.T) {
		eventID := uuid.New()
		_, err := storage.GetEvent(eventID)
		require.Error(t, err)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		event := Event{
			ID:          uuid.New(),
			Title:       "Test Event",
			Date:        time.Now(),
			Duration:    time.Hour,
			Description: "This is a test event",
			User:        uuid.New(),
			Notify:      time.Now().Add(-time.Minute),
		}

		ctx := context.Background()

		err := storage.CreateEvent(ctx, event)
		require.NoError(t, err)

		err = storage.DeleteEvent(event.ID)
		require.NoError(t, err)

		_, err = storage.GetEvent(event.ID)
		require.Error(t, err)
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		event := Event{
			ID:          uuid.New(),
			Title:       "Test Event",
			Date:        time.Now(),
			Duration:    time.Hour,
			Description: "This is a test event",
			User:        uuid.New(),
			Notify:      time.Now().Add(-time.Minute),
		}

		ctx := context.Background()

		err := storage.CreateEvent(ctx, event)
		require.NoError(t, err)

		updatedEvent := Event{
			ID:          event.ID,
			Title:       "Updated Event",
			Date:        event.Date,
			Duration:    event.Duration,
			Description: "This is an updated test event",
			User:        event.User,
			Notify:      event.Notify,
		}

		err = storage.UpdateEvent(updatedEvent)
		require.NoError(t, err)

		storedEvent, err := storage.GetEvent(updatedEvent.ID)
		require.NoError(t, err)
		require.Equal(t, updatedEvent, storedEvent)
	})
}
