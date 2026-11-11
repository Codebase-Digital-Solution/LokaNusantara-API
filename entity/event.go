package entity

import "time"

// Event is a cultural or religious occasion shown in the Agenda Budaya.
type Event struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(191);not null;uniqueIndex:idx_event_name_date" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Category    string     `gorm:"type:varchar(50);index" json:"category"`
	StartDate   time.Time  `gorm:"type:date;not null;uniqueIndex:idx_event_name_date;index" json:"start_date"`
	EndDate     *time.Time `gorm:"type:date" json:"end_date,omitempty"`
	// IsEstimate marks dates that follow a lunar calendar or an organiser's
	// decision and may shift; the app labels them as "Perkiraan".
	IsEstimate bool      `gorm:"not null;default:false" json:"is_estimate"`
	Region     string    `gorm:"type:varchar(100)" json:"region"`
	LocationID *string   `gorm:"type:char(36);index" json:"location_id,omitempty"`
	Location   *Location `gorm:"foreignKey:LocationID;constraint:OnDelete:SET NULL" json:"location,omitempty"`
}

func (Event) TableName() string {
	return "events"
}
