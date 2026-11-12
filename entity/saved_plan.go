package entity

import "time"

// SavedPlan is a generated visit plan a signed-in user kept. The plan itself
// is stored as JSON (the planner's TripPlan) so its shape can evolve.
type SavedPlan struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     string    `gorm:"type:char(36);index;not null" json:"-"`
	LocationID string    `gorm:"type:char(36);index;not null" json:"location_id"`
	PlanJSON   string    `gorm:"type:mediumtext;not null" json:"-"`
	CreatedAt  time.Time `json:"created_at"`

	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Location *Location `gorm:"foreignKey:LocationID;constraint:OnDelete:CASCADE" json:"location,omitempty"`
}

func (SavedPlan) TableName() string {
	return "saved_plans"
}
