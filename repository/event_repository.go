package repository

import (
	"api/entity"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepository interface {
	// GetBetween returns events overlapping [from, to], earliest first.
	GetBetween(from, to time.Time, limit, offset int) ([]entity.Event, error)
	// GetAll returns every event, newest start date first (admin listing).
	GetAll() ([]entity.Event, error)
	GetByID(id uint) (*entity.Event, error)
	Create(event *entity.Event) error
	Update(event *entity.Event) error
	Delete(id uint) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) GetBetween(from, to time.Time, limit, offset int) ([]entity.Event, error) {
	var events []entity.Event
	err := r.db.
		Preload("Location.Category").
		Where("start_date <= ? AND COALESCE(end_date, start_date) >= ?", to, from).
		Order("start_date ASC, name ASC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

func (r *eventRepository) GetAll() ([]entity.Event, error) {
	var events []entity.Event
	err := r.db.
		Preload("Location", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Order("start_date DESC, name ASC").
		Find(&events).Error
	return events, err
}

func (r *eventRepository) GetByID(id uint) (*entity.Event, error) {
	var event entity.Event
	err := r.db.First(&event, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("event not found")
	}
	return &event, err
}

func (r *eventRepository) Create(event *entity.Event) error {
	return r.db.Omit(clause.Associations).Create(event).Error
}

func (r *eventRepository) Update(event *entity.Event) error {
	return r.db.Omit(clause.Associations).Save(event).Error
}

func (r *eventRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Event{}, id).Error
}
