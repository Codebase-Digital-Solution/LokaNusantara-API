package entity

// Language is a language supported by localized content.
type Language struct {
	Code      string `gorm:"primaryKey;type:varchar(10)" json:"code"`
	Name      string `gorm:"type:varchar(80);not null" json:"name"`
	IsDefault bool   `gorm:"not null;default:false" json:"is_default"`
}

func (Language) TableName() string { return "languages" }
