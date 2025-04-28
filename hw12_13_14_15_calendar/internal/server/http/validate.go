package internalhttp

import (
	"errors"
	"github.com/DEMAxx/demin/hw12_13_14_15_calendar/events/pb"
	"time"
)

func validateEventCreateRequest(req *pb.EventCreateRequest) error {
	event := req.GetEvent()
	if event.GetTitle() == "" {
		return errors.New("title is required")
	}
	if _, err := time.Parse("2006-01-02", event.GetDate()); err != nil {
		return errors.New("invalid date format, expected YYYY-MM-DD")
	}
	if event.GetDuration() == nil {
		return errors.New("duration is required")
	}
	if event.GetUser() == "" {
		return errors.New("user ID is required")
	}
	return nil
}
