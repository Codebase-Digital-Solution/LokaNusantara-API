package repository

import (
	"api/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookmarkRepository interface {
	// GetLocations returns the user's saved places, most recently saved first.
	GetLocations(userID string) ([]entity.Location, error)
	// Add is idempotent: saving a place twice keeps one row.
	Add(userID, locationID string) error
	Remove(userID, locationID string) error
}

type bookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) BookmarkRepository {
	return &bookmarkRepository{db: db}
}

func (r *bookmarkRepository) GetLocations(userID string) ([]entity.Location, error) {
	var locations []entity.Location
	err := r.db.
		Select(listSelect).
		Joins("JOIN bookmarks ON bookmarks.location_id = locations.id").
		Preload("Translations").Preload("Category.Translations").
		Where("bookmarks.user_id = ?", userID).
		Order("bookmarks.created_at DESC").
		Find(&locations).Error
	return locations, err
}

func (r *bookmarkRepository) Add(userID, locationID string) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).
		Omit(clause.Associations).
		Create(&entity.Bookmark{UserID: userID, LocationID: locationID}).Error
}

func (r *bookmarkRepository) Remove(userID, locationID string) error {
	return r.db.Where("user_id = ? AND location_id = ?", userID, locationID).
		Delete(&entity.Bookmark{}).Error
}
