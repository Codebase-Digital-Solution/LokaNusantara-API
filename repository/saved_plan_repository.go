package repository

import (
	"api/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SavedPlanRepository interface {
	// GetByUser returns the user's plans, newest first, with their place.
	GetByUser(userID string) ([]entity.SavedPlan, error)
	CountByUser(userID string) (int64, error)
	Create(plan *entity.SavedPlan) error
	// Delete removes the plan only when it belongs to userID; it reports
	// whether a row was deleted.
	Delete(userID string, id uint) (bool, error)
}

type savedPlanRepository struct {
	db *gorm.DB
}

func NewSavedPlanRepository(db *gorm.DB) SavedPlanRepository {
	return &savedPlanRepository{db: db}
}

func (r *savedPlanRepository) GetByUser(userID string) ([]entity.SavedPlan, error) {
	var plans []entity.SavedPlan
	err := r.db.
		Preload("Location", func(db *gorm.DB) *gorm.DB { return db.Select(listSelect) }).
		Preload("Location.Category").
		Where("user_id = ?", userID).
		Order("created_at DESC, id DESC").
		Find(&plans).Error
	return plans, err
}

func (r *savedPlanRepository) CountByUser(userID string) (int64, error) {
	var n int64
	err := r.db.Model(&entity.SavedPlan{}).Where("user_id = ?", userID).Count(&n).Error
	return n, err
}

func (r *savedPlanRepository) Create(plan *entity.SavedPlan) error {
	return r.db.Omit(clause.Associations).Create(plan).Error
}

func (r *savedPlanRepository) Delete(userID string, id uint) (bool, error) {
	res := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entity.SavedPlan{})
	return res.RowsAffected > 0, res.Error
}
