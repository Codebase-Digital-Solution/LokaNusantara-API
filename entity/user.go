package entity

import (
	"time"
)

type User struct {
	ID       string `gorm:"primaryKey;type:char(36)" json:"id"`
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Role     string `gorm:"type:varchar(20);not null;default:user;index" json:"role"`
	// GoogleID is the Google account "sub" for users who signed in with Google.
	GoogleID  *string   `gorm:"type:varchar(64);uniqueIndex" json:"-"`
	AvatarURL string    `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Reviews []Review `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"reviews,omitempty"`
}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

func (User) TableName() string {
	return "users"
}
