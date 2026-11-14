package localize

import (
	"api/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveLocationTranslations(db *gorm.DB, id string, values map[string]entity.LocationTranslation) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for code, value := range values {
			if _, ok := Language(code); !ok {
				continue
			}
			value.LocationID, value.LanguageCode = id, code
			if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "location_id"}, {Name: "language_code"}}, DoUpdates: clause.AssignmentColumns([]string{"name", "description", "history", "address"})}).Create(&value).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func SaveCategoryTranslations(db *gorm.DB, id uint, values map[string]entity.CategoryTranslation) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for code, value := range values {
			if _, ok := Language(code); !ok {
				continue
			}
			value.CategoryID, value.LanguageCode = id, code
			if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "category_id"}, {Name: "language_code"}}, DoUpdates: clause.AssignmentColumns([]string{"label"})}).Create(&value).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func SaveEventTranslations(db *gorm.DB, id uint, values map[string]entity.EventTranslation) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for code, value := range values {
			if _, ok := Language(code); !ok {
				continue
			}
			value.EventID, value.LanguageCode = id, code
			if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "event_id"}, {Name: "language_code"}}, DoUpdates: clause.AssignmentColumns([]string{"name", "description", "category", "region"})}).Create(&value).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func SaveTrailTranslations(db *gorm.DB, id uint, values map[string]entity.TrailTranslation, stopValues []map[string]string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for code, value := range values {
			if _, ok := Language(code); !ok {
				continue
			}
			value.TrailID, value.LanguageCode = id, code
			if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "trail_id"}, {Name: "language_code"}}, DoUpdates: clause.AssignmentColumns([]string{"title", "subtitle", "story"})}).Create(&value).Error; err != nil {
				return err
			}
		}
		if len(stopValues) == 0 {
			return nil
		}
		var stops []entity.TrailStop
		if err := tx.Where("trail_id = ?", id).Order("position ASC").Find(&stops).Error; err != nil {
			return err
		}
		for i, stop := range stops {
			if i >= len(stopValues) {
				break
			}
			for code, note := range stopValues[i] {
				if _, ok := Language(code); !ok {
					continue
				}
				value := entity.TrailStopTranslation{TrailStopID: stop.ID, LanguageCode: code, Note: note}
				if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "trail_stop_id"}, {Name: "language_code"}}, DoUpdates: clause.AssignmentColumns([]string{"note"})}).Create(&value).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
