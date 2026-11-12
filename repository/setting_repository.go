package repository

import (
	"api/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingRepository interface {
	// GetAll returns every stored setting keyed by name.
	GetAll() (map[string]string, error)
	// SetMany upserts the given settings in one transaction.
	SetMany(values map[string]string) error
	Delete(key string) error
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) GetAll() (map[string]string, error) {
	if r.db == nil {
		return map[string]string{}, nil
	}
	var rows []entity.AppSetting
	if err := r.db.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		out[row.Key] = row.Value
	}
	return out, nil
}

func (r *settingRepository) SetMany(values map[string]string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for k, v := range values {
			err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "key"}},
				DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
			}).Create(&entity.AppSetting{Key: k, Value: v}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *settingRepository) Delete(key string) error {
	return r.db.Delete(&entity.AppSetting{}, "`key` = ?", key).Error
}
