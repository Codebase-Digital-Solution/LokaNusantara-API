package entity

// Trail is a curated themed route ("Jalur Kisah") through several locations.
type Trail struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Slug     string `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Title    string `gorm:"type:varchar(255);not null" json:"title"`
	Subtitle string `gorm:"type:varchar(255)" json:"subtitle"`
	Story    string `gorm:"type:text" json:"story"`
	Color    string `gorm:"type:varchar(9)" json:"color"` // hex, e.g. #A87618
	// ImageURL is the cover photo: "/uploads/..." or a full http(s) URL.
	ImageURL     string             `gorm:"type:varchar(500)" json:"image_url"`
	Stops        []TrailStop        `gorm:"foreignKey:TrailID;constraint:OnDelete:CASCADE" json:"stops"`
	Translations []TrailTranslation `gorm:"foreignKey:TrailID;constraint:OnDelete:CASCADE" json:"translations,omitempty"`
}

func (Trail) TableName() string {
	return "trails"
}

// TrailStop is one ordered stop on a trail.
type TrailStop struct {
	ID           uint                   `gorm:"primaryKey" json:"id"`
	TrailID      uint                   `gorm:"index;not null" json:"trail_id"`
	Position     int                    `gorm:"not null" json:"position"`
	LocationID   string                 `gorm:"type:char(36);index;not null" json:"location_id"`
	Location     *Location              `gorm:"foreignKey:LocationID;constraint:OnDelete:CASCADE" json:"location,omitempty"`
	Note         string                 `gorm:"type:text" json:"note"`
	Translations []TrailStopTranslation `gorm:"foreignKey:TrailStopID;constraint:OnDelete:CASCADE" json:"translations,omitempty"`
}

func (TrailStop) TableName() string {
	return "trail_stops"
}
