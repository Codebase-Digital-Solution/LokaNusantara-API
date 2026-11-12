package entity

import "time"

// Bookmark is a location a user saved in the app. One row per user and place.
type Bookmark struct {
	UserID     string    `gorm:"primaryKey;type:char(36)" json:"user_id"`
	LocationID string    `gorm:"primaryKey;type:char(36);index" json:"location_id"`
	CreatedAt  time.Time `json:"created_at"`

	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Location *Location `gorm:"foreignKey:LocationID;constraint:OnDelete:CASCADE" json:"location,omitempty"`
}

func (Bookmark) TableName() string {
	return "bookmarks"
}
