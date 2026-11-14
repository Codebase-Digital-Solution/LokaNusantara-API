package entity

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	// Position orders categories in the app and admin panel (lowest first).
	Position int `gorm:"not null;default:0;index" json:"position"`

	Locations    []Location            `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT" json:"locations,omitempty"`
	Translations []CategoryTranslation `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"translations,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
