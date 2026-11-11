package usecase

import (
	"api/entity"
	"api/repository"
	"errors"
	"strings"
	"time"
)

// EventInput is the editable part of an event, as sent by the admin panel.
type EventInput struct {
	Name        string
	Description string
	Category    string
	StartDate   time.Time
	EndDate     *time.Time
	IsEstimate  bool
	Region      string
	LocationID  *string
}

type EventUsecase interface {
	GetEvents(from, to time.Time, limit, offset int) ([]entity.Event, error)
	GetAllEvents() ([]entity.Event, error)
	CreateEvent(in EventInput) (*entity.Event, error)
	UpdateEvent(id uint, in EventInput) (*entity.Event, error)
	DeleteEvent(id uint) error
}

type eventUsecase struct {
	eventRepo    repository.EventRepository
	locationRepo repository.LocationRepository
}

func NewEventUsecase(eventRepo repository.EventRepository, locationRepo repository.LocationRepository) EventUsecase {
	return &eventUsecase{eventRepo: eventRepo, locationRepo: locationRepo}
}

func (u *eventUsecase) GetEvents(from, to time.Time, limit, offset int) ([]entity.Event, error) {
	if to.Before(from) {
		return nil, errors.New("'to' must not be before 'from'")
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return u.eventRepo.GetBetween(from, to, limit, offset)
}

func (u *eventUsecase) GetAllEvents() ([]entity.Event, error) {
	return u.eventRepo.GetAll()
}

func (u *eventUsecase) validate(in *EventInput) error {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return errors.New("event name is required")
	}
	if in.StartDate.IsZero() {
		return errors.New("start date is required")
	}
	if in.EndDate != nil && in.EndDate.Before(in.StartDate) {
		return errors.New("end date must not be before start date")
	}
	if in.LocationID != nil && *in.LocationID == "" {
		in.LocationID = nil
	}
	if in.LocationID != nil {
		if _, err := u.locationRepo.GetByID(*in.LocationID); err != nil {
			return errors.New("location not found")
		}
	}
	return nil
}

func apply(e *entity.Event, in EventInput) {
	e.Name = in.Name
	e.Description = in.Description
	e.Category = in.Category
	e.StartDate = in.StartDate
	e.EndDate = in.EndDate
	e.IsEstimate = in.IsEstimate
	e.Region = in.Region
	e.LocationID = in.LocationID
	e.Location = nil
}

func (u *eventUsecase) CreateEvent(in EventInput) (*entity.Event, error) {
	if err := u.validate(&in); err != nil {
		return nil, err
	}
	event := &entity.Event{}
	apply(event, in)
	if err := u.eventRepo.Create(event); err != nil {
		return nil, errors.New("failed to create event (name and start date must be unique)")
	}
	return event, nil
}

func (u *eventUsecase) UpdateEvent(id uint, in EventInput) (*entity.Event, error) {
	if err := u.validate(&in); err != nil {
		return nil, err
	}
	event, err := u.eventRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	apply(event, in)
	if err := u.eventRepo.Update(event); err != nil {
		return nil, errors.New("failed to update event (name and start date must be unique)")
	}
	return event, nil
}

func (u *eventUsecase) DeleteEvent(id uint) error {
	if _, err := u.eventRepo.GetByID(id); err != nil {
		return err
	}
	return u.eventRepo.Delete(id)
}
