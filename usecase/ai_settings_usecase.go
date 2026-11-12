package usecase

import (
	"api/repository"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Setting keys stored in app_settings.
const (
	settingOpenAIKey       = "openai.api_key"
	settingOpenAIModel     = "openai.model"
	settingMaxInputChars   = "openai.max_input_chars"
	settingMaxOutputTokens = "openai.max_output_tokens"
)

// Defaults and limits for the admin-editable AI settings.
const (
	DefaultMaxInputChars   = 500
	DefaultMaxOutputTokens = 2000
	MinMaxInputChars       = 50
	MaxMaxInputChars       = 4000
	MinMaxOutputTokens     = 300
	MaxMaxOutputTokens     = 16000
)

// AISettings is what the OpenAI-backed features use on each request.
type AISettings struct {
	APIKey string
	Model  string
	// MaxInputChars caps a visitor's question (Tanya Loka) and notes
	// (planner).
	MaxInputChars int
	// MaxOutputTokens caps the length of each AI answer.
	MaxOutputTokens int
}

// AISettingsView is what the admin panel may see: never the key itself.
type AISettingsView struct {
	APIKeySet       bool   `json:"api_key_set"`
	APIKeyHint      string `json:"api_key_hint"`   // last 4 characters
	APIKeySource    string `json:"api_key_source"` // "database", "env" or ""
	Model           string `json:"model"`
	MaxInputChars   int    `json:"max_input_chars"`
	MaxOutputTokens int    `json:"max_output_tokens"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

// AISettingsUpdate is the admin form. A nil APIKey keeps the stored key;
// ClearAPIKey removes the stored key so the env value (if any) applies.
type AISettingsUpdate struct {
	APIKey          *string `json:"api_key"`
	ClearAPIKey     bool    `json:"clear_api_key"`
	Model           string  `json:"model"`
	MaxInputChars   int     `json:"max_input_chars"`
	MaxOutputTokens int     `json:"max_output_tokens"`
}

// AISettingsSource gives the current settings; it is cheap to call per
// request.
type AISettingsSource interface {
	Current() AISettings
}

type AISettingsUsecase interface {
	AISettingsSource
	View() AISettingsView
	Update(in AISettingsUpdate) (AISettingsView, error)
	// Test checks that the key works and can use the model.
	Test(ctx context.Context) error
}

type aiSettingsUsecase struct {
	repo     repository.SettingRepository
	envKey   string
	envModel string

	mu        sync.RWMutex
	stored    map[string]string
	updatedAt time.Time
}

// NewAISettingsUsecase loads stored settings once and keeps them in memory;
// env values are the fallback when the admin has saved nothing.
func NewAISettingsUsecase(repo repository.SettingRepository, envKey, envModel string) AISettingsUsecase {
	u := &aiSettingsUsecase{repo: repo, envKey: envKey, envModel: envModel, stored: map[string]string{}}
	if repo != nil {
		if all, err := repo.GetAll(); err == nil {
			u.stored = all
		}
	}
	return u
}

func intSetting(m map[string]string, key string, def int) int {
	if v, err := strconv.Atoi(m[key]); err == nil && v > 0 {
		return v
	}
	return def
}

func (u *aiSettingsUsecase) Current() AISettings {
	u.mu.RLock()
	defer u.mu.RUnlock()
	s := AISettings{
		APIKey:          u.envKey,
		Model:           u.envModel,
		MaxInputChars:   intSetting(u.stored, settingMaxInputChars, DefaultMaxInputChars),
		MaxOutputTokens: intSetting(u.stored, settingMaxOutputTokens, DefaultMaxOutputTokens),
	}
	if k := u.stored[settingOpenAIKey]; k != "" {
		s.APIKey = k
	}
	if m := u.stored[settingOpenAIModel]; m != "" {
		s.Model = m
	}
	if s.Model == "" {
		s.Model = "gpt-4o-mini"
	}
	return s
}

func (u *aiSettingsUsecase) View() AISettingsView {
	s := u.Current()
	u.mu.RLock()
	source := ""
	switch {
	case u.stored[settingOpenAIKey] != "":
		source = "database"
	case u.envKey != "":
		source = "env"
	}
	updated := u.updatedAt
	u.mu.RUnlock()

	v := AISettingsView{
		APIKeySet:       s.APIKey != "",
		APIKeySource:    source,
		Model:           s.Model,
		MaxInputChars:   s.MaxInputChars,
		MaxOutputTokens: s.MaxOutputTokens,
	}
	if len(s.APIKey) >= 4 {
		v.APIKeyHint = s.APIKey[len(s.APIKey)-4:]
	}
	if !updated.IsZero() {
		v.UpdatedAt = updated.Format(time.RFC3339)
	}
	return v
}

func validModel(m string) bool {
	if m == "" || len(m) > 100 {
		return false
	}
	for _, r := range m {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '.' || r == '_' || r == ':') {
			return false
		}
	}
	return true
}

func (u *aiSettingsUsecase) Update(in AISettingsUpdate) (AISettingsView, error) {
	in.Model = strings.TrimSpace(in.Model)
	switch {
	case !validModel(in.Model):
		return AISettingsView{}, errors.New("nama model tidak valid")
	case in.MaxInputChars < MinMaxInputChars || in.MaxInputChars > MaxMaxInputChars:
		return AISettingsView{}, fmt.Errorf("max input harus %d sampai %d karakter", MinMaxInputChars, MaxMaxInputChars)
	case in.MaxOutputTokens < MinMaxOutputTokens || in.MaxOutputTokens > MaxMaxOutputTokens:
		return AISettingsView{}, fmt.Errorf("max output harus %d sampai %d token", MinMaxOutputTokens, MaxMaxOutputTokens)
	}

	values := map[string]string{
		settingOpenAIModel:     in.Model,
		settingMaxInputChars:   strconv.Itoa(in.MaxInputChars),
		settingMaxOutputTokens: strconv.Itoa(in.MaxOutputTokens),
	}
	if in.APIKey != nil {
		key := strings.TrimSpace(*in.APIKey)
		if key != "" {
			if !strings.HasPrefix(key, "sk-") || len(key) < 20 {
				return AISettingsView{}, errors.New("API key OpenAI tidak valid (harus diawali sk-)")
			}
			values[settingOpenAIKey] = key
		}
	}

	if err := u.repo.SetMany(values); err != nil {
		return AISettingsView{}, err
	}
	if in.ClearAPIKey {
		if err := u.repo.Delete(settingOpenAIKey); err != nil {
			return AISettingsView{}, err
		}
	}

	u.mu.Lock()
	for k, v := range values {
		u.stored[k] = v
	}
	if in.ClearAPIKey {
		delete(u.stored, settingOpenAIKey)
	}
	u.updatedAt = time.Now()
	u.mu.Unlock()
	return u.View(), nil
}

// Test asks OpenAI for the configured model's metadata: it fails on a wrong
// key, an unknown model, or a model this key cannot use.
func (u *aiSettingsUsecase) Test(ctx context.Context) error {
	s := u.Current()
	if s.APIKey == "" {
		return errors.New("API key belum diisi")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.openai.com/v1/models/"+url.PathEscape(s.Model), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	res, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		return errors.New("tidak bisa menghubungi OpenAI")
	}
	defer res.Body.Close()
	_, _ = io.Copy(io.Discard, res.Body)
	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return errors.New("API key ditolak OpenAI")
	case http.StatusNotFound:
		return fmt.Errorf("model %q tidak tersedia untuk key ini", s.Model)
	default:
		return fmt.Errorf("OpenAI membalas status %d", res.StatusCode)
	}
}
