package usecase

import (
	"api/entity"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestValidatePlan(t *testing.T) {
	ok := PlanRequest{LocationID: "x", Budget: 500_000, People: 2, Days: 1}
	if err := validatePlan(&ok, 300); err != nil {
		t.Fatalf("valid request rejected: %v", err)
	}
	bad := []PlanRequest{
		{Budget: 500_000, People: 2, Days: 1},
		{LocationID: "x", Budget: 5_000, People: 2, Days: 1},
		{LocationID: "x", Budget: 500_000, People: 0, Days: 1},
		{LocationID: "x", Budget: 500_000, People: 2, Days: 9},
		{LocationID: "x", Budget: 500_000, People: 2, Days: 1, Notes: strings.Repeat("a", 301)},
	}
	for i, r := range bad {
		err := validatePlan(&r, 300)
		var inErr *PlanInputError
		if !errors.As(err, &inErr) {
			t.Errorf("case %d: want PlanInputError, got %v", i, err)
		}
	}
}

func TestPlanDisabledWithoutKey(t *testing.T) {
	u := NewPlannerUsecase(NewAISettingsUsecase(nil, "", ""), nil, nil)
	if _, err := u.Plan(t.Context(), PlanRequest{}); !errors.Is(err, ErrPlannerDisabled) {
		t.Fatalf("want ErrPlannerDisabled, got %v", err)
	}
}

func TestBuildPlanPromptMentionsOnlyTargetEvents(t *testing.T) {
	target := &entity.Location{ID: "t1", Name: "Pura A", Category: entity.Category{Name: "pura"}}
	other := "zz"
	prompt := buildPlanPrompt(
		PlanRequest{Budget: 300_000, People: 3, Days: 1},
		target,
		[]entity.Location{{ID: "n1", Name: "Warung B", Category: entity.Category{Name: "kuliner"}}},
		[]entity.Event{
			{Name: "Odalan Pura A", LocationID: &target.ID, StartDate: time.Now(), IsEstimate: true},
			{Name: "Acara Lain", LocationID: &other, StartDate: time.Now()},
		},
	)
	for _, want := range []string{"Rp300000", "id=t1", "id=n1", "Odalan Pura A", "(perkiraan)"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing %q", want)
		}
	}
	if strings.Contains(prompt, "Acara Lain") {
		t.Error("prompt includes an event from another place")
	}
}

func TestPlanLanguageName(t *testing.T) {
	cases := map[string]string{
		"":               "Bahasa Indonesia",
		"id":             "Bahasa Indonesia",
		"en":             "English",
		"en-US,en;q=0.9": "English",
		"ko":             "Korean (한국어)",
		"ja_JP":          "Japanese (日本語)",
		"zh":             "Simplified Chinese (简体中文)",
		"fr":             "Bahasa Indonesia",
	}
	for in, want := range cases {
		if got := planLanguageName(in); got != want {
			t.Errorf("planLanguageName(%q) = %q, want %q", in, got, want)
		}
	}
}
