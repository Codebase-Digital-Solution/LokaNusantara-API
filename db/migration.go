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

	if err := db.AutoMigrate(&entity.AppSetting{}); err != nil {
		return fmt.Errorf("failed to migrate AppSetting: %w", err)
	}

	if err := db.AutoMigrate(&entity.SavedPlan{}); err != nil {
		return fmt.Errorf("failed to migrate SavedPlan: %w", err)
	}

	if err := db.AutoMigrate(&entity.Bookmark{}); err != nil {
		return fmt.Errorf("failed to migrate Bookmark: %w", err)
	}

	// Seed initial data if not exist
	seedCategories(db)
	seedCategoryOrder(db)
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

// defaultCategoryOrder puts traditions first, then places of worship, then
// food. Categories not listed keep their place after these.
var defaultCategoryOrder = []string{"tradisi", "pura", "masjid", "gereja", "vihara", "klenteng", "kuliner"}

// seedCategoryOrder gives categories a position once, when none has one yet
// (fresh install or right after the column was added). Later changes from
// the admin panel are left alone.
func seedCategoryOrder(db *gorm.DB) {
	var ordered int64
	db.Model(&entity.Category{}).Where("position > 0").Count(&ordered)
	if ordered > 0 {
		return
	}
	var all []entity.Category
	db.Order("id ASC").Find(&all)
	rank := map[string]int{}
	for i, name := range defaultCategoryOrder {
		rank[name] = i + 1
	}
	next := len(defaultCategoryOrder) + 1
	for _, c := range all {
		pos, ok := rank[c.Name]
		if !ok {
			pos = next
			next++
		}
		db.Model(&entity.Category{}).Where("id = ?", c.ID).Update("position", pos)
	}
}
