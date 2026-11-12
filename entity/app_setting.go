package entity

import "time"

// AppSetting is one key/value setting editable from the admin panel.
type AppSetting struct {
	Key       string    `gorm:"primaryKey;type:varchar(100)" json:"key"`
	Value     string    `gorm:"type:text" json:"-"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AppSetting) TableName() string {
	return "app_settings"
}
