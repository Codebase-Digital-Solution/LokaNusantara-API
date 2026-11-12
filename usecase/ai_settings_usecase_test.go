package usecase

import (
	"strings"
	"testing"
)

type memSettings struct{ m map[string]string }

func (r *memSettings) GetAll() (map[string]string, error) {
	out := map[string]string{}
	for k, v := range r.m {
		out[k] = v
	}
	return out, nil
}

func (r *memSettings) SetMany(v map[string]string) error {
	for k, val := range v {
		r.m[k] = val
	}
	return nil
}

func (r *memSettings) Delete(key string) error {
	delete(r.m, key)
	return nil
}

func TestAISettingsFallsBackToEnv(t *testing.T) {
	u := NewAISettingsUsecase(&memSettings{m: map[string]string{}}, "sk-envkey-123456789012345", "")
	s := u.Current()
	if s.APIKey != "sk-envkey-123456789012345" || s.Model != "gpt-4o-mini" ||
		s.MaxInputChars != DefaultMaxInputChars || s.MaxOutputTokens != DefaultMaxOutputTokens {
		t.Fatalf("defaults = %+v", s)
	}
	v := u.View()
	if !v.APIKeySet || v.APIKeySource != "env" || v.APIKeyHint != "2345" {
		t.Fatalf("view = %+v", v)
	}
}

func TestAISettingsUpdateNeverExposesKey(t *testing.T) {
	repo := &memSettings{m: map[string]string{}}
	u := NewAISettingsUsecase(repo, "", "gpt-4o-mini")
	key := "sk-proj-abcdefghijklmnopqrstuvwxyzWXYZ"
	v, err := u.Update(AISettingsUpdate{APIKey: &key, Model: "gpt-4.1-mini", MaxInputChars: 800, MaxOutputTokens: 3000})
	if err != nil {
		t.Fatal(err)
	}
	if v.APIKeySource != "database" || v.APIKeyHint != "WXYZ" || v.Model != "gpt-4.1-mini" {
		t.Fatalf("view = %+v", v)
	}
	if got := u.Current(); got.APIKey != key || got.MaxInputChars != 800 || got.MaxOutputTokens != 3000 {
		t.Fatalf("current = %+v", got)
	}

	// Saving again without a key keeps the stored one.
	if _, err := u.Update(AISettingsUpdate{Model: "gpt-4o-mini", MaxInputChars: 500, MaxOutputTokens: 2000}); err != nil {
		t.Fatal(err)
	}
	if u.Current().APIKey != key {
		t.Fatal("key was lost when saving without one")
	}

	// A fresh instance (API restart) reads the stored values back.
	if NewAISettingsUsecase(repo, "", "").Current().APIKey != key {
		t.Fatal("stored key not loaded on start")
	}

	// Clearing removes it.
	if _, err := u.Update(AISettingsUpdate{ClearAPIKey: true, Model: "gpt-4o-mini", MaxInputChars: 500, MaxOutputTokens: 2000}); err != nil {
		t.Fatal(err)
	}
	if u.View().APIKeySet {
		t.Fatal("key still set after clearing")
	}
}

func TestAISettingsRejectsBadInput(t *testing.T) {
	u := NewAISettingsUsecase(&memSettings{m: map[string]string{}}, "", "")
	bad := "not-a-key"
	cases := []AISettingsUpdate{
		{Model: "gpt 4o mini", MaxInputChars: 500, MaxOutputTokens: 2000},
		{Model: "gpt-4o-mini", MaxInputChars: 10, MaxOutputTokens: 2000},
		{Model: "gpt-4o-mini", MaxInputChars: 500, MaxOutputTokens: 100},
		{APIKey: &bad, Model: "gpt-4o-mini", MaxInputChars: 500, MaxOutputTokens: 2000},
	}
	for i, c := range cases {
		if _, err := u.Update(c); err == nil {
			t.Errorf("case %d accepted", i)
		}
	}
	if strings.Contains(u.View().APIKeyHint, "key") {
		t.Error("rejected key was stored")
	}
}
