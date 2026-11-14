package localize

import "api/entity"

const DefaultLanguage = "id"

var SupportedLanguages = []entity.Language{
	{Code: "id", Name: "Indonesia", IsDefault: true},
	{Code: "en", Name: "English"},
	{Code: "ja", Name: "日本語"},
	{Code: "ko", Name: "한국어"},
}

func Normalize(code string) string {
	for _, language := range SupportedLanguages {
		if language.Code == code {
			return code
		}
	}
	return DefaultLanguage
}

func Language(code string) (entity.Language, bool) {
	for _, language := range SupportedLanguages {
		if language.Code == code {
			return language, true
		}
	}
	return entity.Language{}, false
}

func choose[T any](rows []T, language, fallback string, code func(T) string) (T, bool) {
	var fallbackRow T
	foundFallback := false
	for _, row := range rows {
		if code(row) == language {
			return row, true
		}
		if code(row) == fallback {
			fallbackRow, foundFallback = row, true
		}
	}
	return fallbackRow, foundFallback
}

func Location(value *entity.Location, language string, keepTranslations bool) {
	if value == nil {
		return
	}
	language = Normalize(language)
	if row, ok := choose(value.Translations, language, DefaultLanguage, func(t entity.LocationTranslation) string { return t.LanguageCode }); ok {
		if row.Name != "" {
			value.Name = row.Name
		}
		if row.Description != "" {
			value.Description = row.Description
		}
		if row.History != "" {
			value.History = row.History
		}
		if row.Address != "" {
			value.Address = row.Address
		}
	}
	Category(&value.Category, language, keepTranslations)
	if !keepTranslations {
		value.Translations = nil
	}
}

func Category(value *entity.Category, language string, keepTranslations bool) {
	if value == nil {
		return
	}
	language = Normalize(language)
	if row, ok := choose(value.Translations, language, DefaultLanguage, func(t entity.CategoryTranslation) string { return t.LanguageCode }); ok && row.Label != "" {
		value.Name = row.Label
	}
	if !keepTranslations {
		value.Translations = nil
	}
}

func Event(value *entity.Event, language string, keepTranslations bool) {
	if value == nil {
		return
	}
	language = Normalize(language)
	if row, ok := choose(value.Translations, language, DefaultLanguage, func(t entity.EventTranslation) string { return t.LanguageCode }); ok {
		if row.Name != "" {
			value.Name = row.Name
		}
		if row.Description != "" {
			value.Description = row.Description
		}
		if row.Category != "" {
			value.Category = row.Category
		}
		if row.Region != "" {
			value.Region = row.Region
		}
	}
	Location(value.Location, language, keepTranslations)
	if !keepTranslations {
		value.Translations = nil
	}
}

func Trail(value *entity.Trail, language string, keepTranslations bool) {
	if value == nil {
		return
	}
	language = Normalize(language)
	if row, ok := choose(value.Translations, language, DefaultLanguage, func(t entity.TrailTranslation) string { return t.LanguageCode }); ok {
		if row.Title != "" {
			value.Title = row.Title
		}
		if row.Subtitle != "" {
			value.Subtitle = row.Subtitle
		}
		if row.Story != "" {
			value.Story = row.Story
		}
	}
	for i := range value.Stops {
		stop := &value.Stops[i]
		if row, ok := choose(stop.Translations, language, DefaultLanguage, func(t entity.TrailStopTranslation) string { return t.LanguageCode }); ok && row.Note != "" {
			stop.Note = row.Note
		}
		Location(stop.Location, language, keepTranslations)
		if !keepTranslations {
			stop.Translations = nil
		}
	}
	if !keepTranslations {
		value.Translations = nil
	}
}
