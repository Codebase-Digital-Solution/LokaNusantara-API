package entity

// Translation rows have a composite unique index on parent ID and language.
type LocationTranslation struct {
	ID           uint     `gorm:"primaryKey" json:"id,omitempty"`
	LocationID   string   `gorm:"type:char(36);not null;uniqueIndex:idx_location_translation_language" json:"location_id,omitempty"`
	LanguageCode string   `gorm:"type:varchar(10);not null;uniqueIndex:idx_location_translation_language" json:"language_code"`
	Language     Language `gorm:"foreignKey:LanguageCode;references:Code" json:"-"`
	Name         string   `gorm:"type:varchar(255)" json:"name"`
	Description  string   `gorm:"type:longtext" json:"description"`
	History      string   `gorm:"type:longtext" json:"history"`
	Address      string   `gorm:"type:varchar(500)" json:"address"`
}

func (LocationTranslation) TableName() string { return "location_translations" }

type CategoryTranslation struct {
	ID           uint     `gorm:"primaryKey" json:"id,omitempty"`
	CategoryID   uint     `gorm:"not null;uniqueIndex:idx_category_translation_language" json:"category_id,omitempty"`
	LanguageCode string   `gorm:"type:varchar(10);not null;uniqueIndex:idx_category_translation_language" json:"language_code"`
	Language     Language `gorm:"foreignKey:LanguageCode;references:Code" json:"-"`
	Label        string   `gorm:"type:varchar(100);not null" json:"label"`
}

func (CategoryTranslation) TableName() string { return "category_translations" }

type EventTranslation struct {
	ID           uint     `gorm:"primaryKey" json:"id,omitempty"`
	EventID      uint     `gorm:"not null;uniqueIndex:idx_event_translation_language" json:"event_id,omitempty"`
	LanguageCode string   `gorm:"type:varchar(10);not null;uniqueIndex:idx_event_translation_language" json:"language_code"`
	Language     Language `gorm:"foreignKey:LanguageCode;references:Code" json:"-"`
	Name         string   `gorm:"type:varchar(191)" json:"name"`
	Description  string   `gorm:"type:text" json:"description"`
	Category     string   `gorm:"type:varchar(50)" json:"category"`
	Region       string   `gorm:"type:varchar(100)" json:"region"`
}

func (EventTranslation) TableName() string { return "event_translations" }

type TrailTranslation struct {
	ID           uint     `gorm:"primaryKey" json:"id,omitempty"`
	TrailID      uint     `gorm:"not null;uniqueIndex:idx_trail_translation_language" json:"trail_id,omitempty"`
	LanguageCode string   `gorm:"type:varchar(10);not null;uniqueIndex:idx_trail_translation_language" json:"language_code"`
	Language     Language `gorm:"foreignKey:LanguageCode;references:Code" json:"-"`
	Title        string   `gorm:"type:varchar(255)" json:"title"`
	Subtitle     string   `gorm:"type:varchar(255)" json:"subtitle"`
	Story        string   `gorm:"type:text" json:"story"`
}

func (TrailTranslation) TableName() string { return "trail_translations" }

type TrailStopTranslation struct {
	ID           uint     `gorm:"primaryKey" json:"id,omitempty"`
	TrailStopID  uint     `gorm:"not null;uniqueIndex:idx_trail_stop_translation_language" json:"trail_stop_id,omitempty"`
	LanguageCode string   `gorm:"type:varchar(10);not null;uniqueIndex:idx_trail_stop_translation_language" json:"language_code"`
	Language     Language `gorm:"foreignKey:LanguageCode;references:Code" json:"-"`
	Note         string   `gorm:"type:text" json:"note"`
}

func (TrailStopTranslation) TableName() string { return "trail_stop_translations" }
