package usecase

import (
	"api/entity"
	"testing"
)

func TestParseAnswerStripsMarkersAndCollectsPlaces(t *testing.T) {
	id := "11111111-2222-3333-4444-555555555555"
	places := map[string]entity.Location{id: {ID: id, Name: "Pura Tanah Lot"}}
	got := parseAnswer("Kunjungi Pura Tanah Lot [["+id+"]] saat senja. Lalu [[99999999-2222-3333-4444-555555555555]] tidak dikenal. Lagi [["+id+"]].", places)

	want := "Kunjungi Pura Tanah Lot saat senja. Lalu tidak dikenal. Lagi."
	if got.Answer != want {
		t.Errorf("answer = %q, want %q", got.Answer, want)
	}
	if len(got.Locations) != 1 || got.Locations[0].ID != id {
		t.Errorf("locations = %+v", got.Locations)
	}
}

func TestBuildMessagesAlternatesAndEndsWithQuestion(t *testing.T) {
	history := []AssistantMessage{
		{Role: "assistant", Content: "sapaan tanpa pertanyaan"}, // dropped: must start with user
		{Role: "user", Content: "Apa itu Galungan?"},
		{Role: "assistant", Content: "Hari kemenangan dharma."},
		{Role: "user", Content: "pertanyaan tanpa jawaban"}, // dropped: no reply follows
	}
	msgs := buildMessages(history, "Kapan Kuningan?")
	if len(msgs) != 3 {
		t.Fatalf("got %d messages, want 3", len(msgs))
	}
	if msgs[0].OfUser == nil || msgs[1].OfAssistant == nil || msgs[2].OfUser == nil {
		t.Errorf("roles do not alternate user/assistant/user")
	}
}

func TestTrimHistoryKeepsAlternatingTurns(t *testing.T) {
	got := trimHistory([]AssistantMessage{
		{Role: "assistant", Content: "sapaan"},
		{Role: "user", Content: "  Apa itu Galungan?  "},
		{Role: "assistant", Content: "Hari kemenangan dharma."},
		{Role: "user", Content: "menggantung"},
	})
	if len(got) != 2 || got[0].Role != "user" || got[0].Content != "Apa itu Galungan?" || got[1].Role != "assistant" {
		t.Fatalf("trimHistory = %+v", got)
	}
}

func TestAskDisabledWithoutAnyKey(t *testing.T) {
	u := NewAssistantUsecase(nil, nil, NewAISettingsUsecase(nil, "", ""))
	if _, err := u.Ask(t.Context(), "Halo", nil); err != ErrAssistantDisabled {
		t.Fatalf("want ErrAssistantDisabled, got %v", err)
	}
}
