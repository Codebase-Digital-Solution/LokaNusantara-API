package db

import (
	"api/entity"
	"fmt"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Auto migrate all entities
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return fmt.Errorf("failed to migrate User: %w", err)
	}

	if err := db.AutoMigrate(&entity.Category{}); err != nil {
		return fmt.Errorf("failed to migrate Category: %w", err)
	}

	if err := db.AutoMigrate(&entity.Location{}); err != nil {
		return fmt.Errorf("failed to migrate Location: %w", err)
	}

	if err := db.AutoMigrate(&entity.Review{}); err != nil {
		return fmt.Errorf("failed to migrate Review: %w", err)
	}

	if err := db.AutoMigrate(&entity.Event{}); err != nil {
		return fmt.Errorf("failed to migrate Event: %w", err)
	}

	if err := db.AutoMigrate(&entity.Trail{}, &entity.TrailStop{}); err != nil {
		return fmt.Errorf("failed to migrate Trail: %w", err)
	}

	// Seed initial data if not exist
	seedCategories(db)
	SeedLocations(db)
	SeedEvents(db)
	SeedTrails(db)

	return nil
}

func seedCategories(db *gorm.DB) {
	categories := []entity.Category{
		{Name: "tradisi"},
		{Name: "pura"},
		{Name: "masjid"},
		{Name: "gereja"},
		{Name: "vihara"},
		{Name: "klenteng"},
		{Name: "kuliner"},
	}

	for _, cat := range categories {
		var count int64
		db.Model(&entity.Category{}).Where("name = ?", cat.Name).Count(&count)
		if count == 0 {
			db.Create(&cat)
		}
	}
}
