package db

import (
	"api/entity"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAdmin makes sure the account named by ADMIN_EMAIL exists and has the
// admin role. An existing user keeps their password and is only promoted.
// Nothing happens when email or password is empty.
func SeedAdmin(db *gorm.DB, name, email, password string) error {
	if email == "" || password == "" {
		log.Println("ADMIN_EMAIL/ADMIN_PASSWORD not set, skipping admin seed")
		return nil
	}

	var user entity.User
	err := db.Where("email = ?", email).First(&user).Error
	if err == nil {
		if user.Role != entity.RoleAdmin {
			if err := db.Model(&user).Update("role", entity.RoleAdmin).Error; err != nil {
				return err
			}
			log.Printf("Promoted %s to admin", email)
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	b := make([]byte, 16)
	rand.Read(b)

	admin := entity.User{
		ID:       hex.EncodeToString(b),
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     entity.RoleAdmin,
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	log.Printf("Created admin account %s", email)
	return nil
}
