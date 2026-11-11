package repository

import (
	"api/entity"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TrailRepository interface {
	GetAll() ([]entity.Trail, error)
	GetByID(id uint) (*entity.Trail, error)
	// Save creates or updates the trail and replaces all of its stops.
	Save(trail *entity.Trail) error
	Delete(id uint) error
}

type trailRepository struct {
	db *gorm.DB
}

func NewTrailRepository(db *gorm.DB) TrailRepository {
	return &trailRepository{db: db}
}

func (r *trailRepository) preloaded() *gorm.DB {
	return r.db.
		Preload("Stops", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Stops.Location.Category")
}

func (r *trailRepository) GetAll() ([]entity.Trail, error) {
	var trails []entity.Trail
	err := r.preloaded().Order("id ASC").Find(&trails).Error
	return trails, err
}

func (r *trailRepository) GetByID(id uint) (*entity.Trail, error) {
	var trail entity.Trail
	err := r.preloaded().First(&trail, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("trail not found")
	}
	return &trail, err
}

func (r *trailRepository) Save(trail *entity.Trail) error {
	stops := trail.Stops
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Save(trail).Error; err != nil {
			return err
		}
		if err := tx.Where("trail_id = ?", trail.ID).Delete(&entity.TrailStop{}).Error; err != nil {
			return err
		}
		for i := range stops {
			stops[i].ID = 0
			stops[i].TrailID = trail.ID
			stops[i].Position = i + 1
			stops[i].Location = nil
		}
		if len(stops) > 0 {
			if err := tx.Omit(clause.Associations).Create(&stops).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *trailRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("trail_id = ?", id).Delete(&entity.TrailStop{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.Trail{}, id).Error
	})
}
