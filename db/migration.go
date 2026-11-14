package db

import (
	"api/entity"
	"api/localize"
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

	if err := db.AutoMigrate(&entity.Language{}, &entity.LocationTranslation{}, &entity.CategoryTranslation{}, &entity.EventTranslation{}, &entity.TrailTranslation{}, &entity.TrailStopTranslation{}); err != nil {
		return fmt.Errorf("failed to migrate translations: %w", err)
	}
	if err := seedLanguages(db); err != nil {
		return err
	}

	// Language rows are reference data required by translation foreign keys.
	// Content/demo seeders are intentionally not run during API startup.
	if err := backfillDefaultTranslations(db); err != nil {
		return fmt.Errorf("failed to backfill Indonesian translations: %w", err)
	}

	return nil
}

// SeedInitialData inserts the optional starter content. It is called by the
// dedicated seed command, never by normal API startup.
func SeedInitialData(db *gorm.DB) error {
	seedCategories(db)
	seedCategoryOrder(db)
	SeedLocations(db)
	SeedEvents(db)
	SeedTrails(db)
	if err := backfillDefaultTranslations(db); err != nil {
		return fmt.Errorf("failed to backfill Indonesian translations for seeded content: %w", err)
	}
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

func seedLanguages(db *gorm.DB) error {
	for _, language := range localize.SupportedLanguages {
		if err := db.Where("code = ?", language.Code).Assign(language).FirstOrCreate(&entity.Language{}).Error; err != nil {
			return fmt.Errorf("failed to seed language %s: %w", language.Code, err)
		}
	}
	return nil
}

// backfillDefaultTranslations copies the existing Indonesian content into the
// translation tables. It is safe to run at every startup and never overwrites
// translations that an administrator has edited.
func backfillDefaultTranslations(db *gorm.DB) error {
	var locations []entity.Location
	if err := db.Find(&locations).Error; err != nil {
		return err
	}
	for _, value := range locations {
		row := entity.LocationTranslation{LocationID: value.ID, LanguageCode: localize.DefaultLanguage}
		if err := db.Where("location_id = ? AND language_code = ?", value.ID, localize.DefaultLanguage).FirstOrCreate(&row, entity.LocationTranslation{
			Name: value.Name, Description: value.Description, History: value.History, Address: value.Address,
		}).Error; err != nil {
			return err
		}
	}

	var categories []entity.Category
	if err := db.Find(&categories).Error; err != nil {
		return err
	}
	for _, value := range categories {
		row := entity.CategoryTranslation{CategoryID: value.ID, LanguageCode: localize.DefaultLanguage}
		if err := db.Where("category_id = ? AND language_code = ?", value.ID, localize.DefaultLanguage).FirstOrCreate(&row, entity.CategoryTranslation{Label: value.Name}).Error; err != nil {
			return err
		}
	}

	var events []entity.Event
	if err := db.Find(&events).Error; err != nil {
		return err
	}
	for _, value := range events {
		row := entity.EventTranslation{EventID: value.ID, LanguageCode: localize.DefaultLanguage}
		if err := db.Where("event_id = ? AND language_code = ?", value.ID, localize.DefaultLanguage).FirstOrCreate(&row, entity.EventTranslation{
			Name: value.Name, Description: value.Description, Category: value.Category, Region: value.Region,
		}).Error; err != nil {
			return err
		}
	}

	var trails []entity.Trail
	if err := db.Preload("Stops").Find(&trails).Error; err != nil {
		return err
	}
	for _, value := range trails {
		row := entity.TrailTranslation{TrailID: value.ID, LanguageCode: localize.DefaultLanguage}
		if err := db.Where("trail_id = ? AND language_code = ?", value.ID, localize.DefaultLanguage).FirstOrCreate(&row, entity.TrailTranslation{
			Title: value.Title, Subtitle: value.Subtitle, Story: value.Story,
		}).Error; err != nil {
			return err
		}
		for _, stop := range value.Stops {
			stopRow := entity.TrailStopTranslation{TrailStopID: stop.ID, LanguageCode: localize.DefaultLanguage}
			if err := db.Where("trail_stop_id = ? AND language_code = ?", stop.ID, localize.DefaultLanguage).FirstOrCreate(&stopRow, entity.TrailStopTranslation{Note: stop.Note}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
