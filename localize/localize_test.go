package localize

import (
	"api/entity"
	"testing"
)

func TestNormalizeUsesSupportedCodesAndDefaultsToIndonesian(t *testing.T) {
	for _, code := range []string{"id", "en", "ja", "ko"} {
		if got := Normalize(code); got != code {
			t.Errorf("Normalize(%q) = %q", code, got)
		}
	}
	if got := Normalize("fr"); got != DefaultLanguage {
		t.Errorf("unsupported language should default to id; got %q", got)
	}
	if got := Normalize(""); got != DefaultLanguage {
		t.Errorf("empty language should default to id; got %q", got)
	}
}

func TestLocationUsesRequestedTranslationAndFallsBackPerFieldToIndonesian(t *testing.T) {
	value := entity.Location{
		Name: "Lokasi", Description: "Deskripsi Indonesia", Address: "Alamat Indonesia",
		Translations: []entity.LocationTranslation{
			{LanguageCode: "id", Name: "Lokasi", Description: "Deskripsi Indonesia", Address: "Alamat Indonesia"},
			{LanguageCode: "en", Name: "Place", Address: ""},
		},
		Category: entity.Category{Name: "tradisi", Translations: []entity.CategoryTranslation{{LanguageCode: "en", Label: "Traditions"}}},
	}
	Location(&value, "en", false)
	if value.Name != "Place" {
		t.Errorf("name = %q, want Place", value.Name)
	}
	if value.Description != "Deskripsi Indonesia" {
		t.Errorf("description should fall back to Indonesian, got %q", value.Description)
	}
	if value.Address != "Alamat Indonesia" {
		t.Errorf("address should fall back to Indonesian, got %q", value.Address)
	}
	if value.Category.Name != "Traditions" {
		t.Errorf("category = %q, want Traditions", value.Category.Name)
	}
	if value.Translations != nil || value.Category.Translations != nil {
		t.Error("public output should omit translation rows")
	}
}

func TestTrailLocalizesStopNotesAndNestedLocations(t *testing.T) {
	value := entity.Trail{
		Title:        "Jalur",
		Translations: []entity.TrailTranslation{{LanguageCode: "ja", Title: "道"}},
		Stops: []entity.TrailStop{{
			Note: "Catatan", Translations: []entity.TrailStopTranslation{{LanguageCode: "ja", Note: "案内"}},
			Location: &entity.Location{Name: "Tempat", Translations: []entity.LocationTranslation{{LanguageCode: "ja", Name: "場所"}}},
		}},
	}
	Trail(&value, "ja", false)
	if value.Title != "道" || value.Stops[0].Note != "案内" || value.Stops[0].Location.Name != "場所" {
		t.Fatalf("Japanese trail fields were not selected: %#v", value)
	}
	if value.Translations != nil || value.Stops[0].Translations != nil {
		t.Error("public output should omit translation rows")
	}
}
