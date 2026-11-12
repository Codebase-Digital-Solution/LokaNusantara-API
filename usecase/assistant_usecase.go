package usecase

import (
	"api/entity"
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
	"gorm.io/gorm"
)

// ErrAssistantDisabled is returned when neither ANTHROPIC_API_KEY nor
// OPENAI_API_KEY is configured.
var ErrAssistantDisabled = errors.New("Tanya Loka belum dikonfigurasi di server")

const (
	assistantModel       = "claude-opus-5-5"
	assistantMaxTokens   = 16000
	knowledgeRefreshTime = time.Hour
	maxHistoryMessages   = 10
	maxMessageChars      = 2000
)

// AssistantMessage is one earlier turn sent by the app.
type AssistantMessage struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}

// AssistantAnswer is the reply plus the locations it mentions.
type AssistantAnswer struct {
	Answer    string            `json:"answer"`
	Locations []entity.Location `json:"locations"`
}

type AssistantUsecase interface {
	Ask(ctx context.Context, question string, history []AssistantMessage) (*AssistantAnswer, error)
}

type assistantUsecase struct {
	db     *gorm.DB
	client *anthropic.Client
	openai *openAIClient

	mu        sync.Mutex
	knowledge string
	places    map[string]entity.Location
	builtAt   time.Time
}

// NewAssistantUsecase returns a usecase that answers only from the
// LokaNusantara database. It uses Claude when client is set, otherwise
// OpenAI with the admin AI settings; with neither it is disabled.
func NewAssistantUsecase(db *gorm.DB, client *anthropic.Client, settings AISettingsSource) AssistantUsecase {
	return &assistantUsecase{db: db, client: client, openai: newOpenAIClient(settings)}
}

const refusalAnswer = "Maaf, Tanya Loka tidak bisa menjawab pertanyaan itu. Coba tanyakan tentang tradisi, rumah ibadah, kuliner, atau agenda budaya."

const assistantInstructions = `Kamu adalah Tanya Loka, pemandu budaya di aplikasi LokaNusantara yang membantu orang mengenal tradisi, rumah ibadah, kuliner, dan agenda budaya di Bali.

Cara menjawab:
- Jawab dalam bahasa yang dipakai pengguna (biasanya Bahasa Indonesia), ramah dan ringkas: beberapa kalimat atau paragraf pendek, tanpa judul markdown.
- Gunakan hanya informasi dari bagian "Data LokaNusantara" di bawah. Jika jawabannya tidak ada di data, katakan terus terang bahwa informasinya belum tersedia di LokaNusantara, lalu sarankan menanyakan kepada pengelola atau tokoh adat setempat. Jangan mengarang nama, tanggal, tahun, atau tokoh.
- Agenda yang bertanda "perkiraan" harus disebut sebagai perkiraan, dan sarankan mengecek pengumuman resmi.
- Setiap kali menyebut lokasi yang ada di data, tulis penanda id-nya tepat setelah nama lokasi, misalnya: Pura Tanah Lot [[1b2c...]]. Gunakan hanya id yang ada di data.
- Hormati semua agama, kepercayaan, dan adat. Jangan menilai atau membandingkan keyakinan.
- Jika pertanyaan di luar budaya, wisata, kuliner, atau agenda budaya Indonesia, tolak dengan sopan dan arahkan kembali ke topik LokaNusantara.`

var etiquetteNotes = map[string]string{
	"pura":     "Kenakan kamen (sarung) dan selendang; perempuan yang sedang haid tidak memasuki area suci; jangan memanjat bangunan suci; jangan berjalan di depan orang yang bersembahyang; hindari menginjak canang.",
	"masjid":   "Berpakaian sopan menutup aurat, perempuan berkerudung; lepas alas kaki; hindari berkunjung saat salat berjamaah kecuali untuk beribadah.",
	"gereja":   "Berpakaian sopan; jaga ketenangan saat misa atau ibadah; minta izin sebelum memotret saat ibadah.",
	"vihara":   "Berpakaian sopan; lepas alas kaki di ruang puja; jangan menyentuh rupang Buddha atau berfoto membelakanginya.",
	"klenteng": "Berpakaian sopan; jangan menyentuh altar dan perlengkapan sembahyang; minta izin sebelum memotret umat.",
	"tradisi":  "Minta izin sebelum memotret warga atau upacara; ikuti aturan desa adat dan pemandu setempat; jangan menyentuh benda sakral.",
	"kuliner":  "Banyak warung tradisional tutup saat hidangan habis; beberapa hidangan Bali mengandung babi, tanyakan bila perlu makanan halal; siapkan uang tunai.",
}

// refresh rebuilds the knowledge text from the database at most once per
// knowledgeRefreshTime. Keeping it stable for an hour keeps the prompt cache
// warm.
func (u *assistantUsecase) refresh() (string, map[string]entity.Location, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.knowledge != "" && time.Since(u.builtAt) < knowledgeRefreshTime {
		return u.knowledge, u.places, nil
	}

	var locations []entity.Location
	if err := u.db.Preload("Category").Order("name ASC").Find(&locations).Error; err != nil {
		return "", nil, err
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	var events []entity.Event
	if err := u.db.
		Where("COALESCE(end_date, start_date) >= ? AND start_date <= ?", today, today.AddDate(1, 0, 0)).
		Order("start_date ASC").Find(&events).Error; err != nil {
		return "", nil, err
	}
	var trails []entity.Trail
	if err := u.db.
		Preload("Stops", func(db *gorm.DB) *gorm.DB { return db.Order("position ASC") }).
		Order("id ASC").Find(&trails).Error; err != nil {
		return "", nil, err
	}

	places := make(map[string]entity.Location, len(locations))
	var b strings.Builder
	b.WriteString("# Data LokaNusantara\n\n## Lokasi\n")
	for _, l := range locations {
		places[l.ID] = l
		fmt.Fprintf(&b, "\n### %s [[%s]]\nKategori: %s\nAlamat: %s\n", l.Name, l.ID, l.Category.Name, l.Address)
		if l.Description != "" {
			fmt.Fprintf(&b, "Deskripsi: %s\n", l.Description)
		}
		if l.History != "" {
			fmt.Fprintf(&b, "Sejarah: %s\n", l.History)
		}
	}

	b.WriteString("\n## Agenda Budaya (12 bulan ke depan)\n")
	for _, e := range events {
		date := e.StartDate.Format("2006-01-02")
		if e.EndDate != nil {
			date += " s.d. " + e.EndDate.Format("2006-01-02")
		}
		estimate := ""
		if e.IsEstimate {
			estimate = " (tanggal perkiraan)"
		}
		place := e.Region
		if e.LocationID != nil {
			if l, ok := places[*e.LocationID]; ok {
				place = fmt.Sprintf("%s [[%s]], %s", l.Name, l.ID, e.Region)
			}
		}
		fmt.Fprintf(&b, "- %s: %s%s. Tempat: %s. %s\n", date, e.Name, estimate, place, e.Description)
	}

	b.WriteString("\n## Jalur Kisah\n")
	for _, t := range trails {
		fmt.Fprintf(&b, "\n### %s\n%s\nUrutan:\n", t.Title, t.Story)
		for _, s := range t.Stops {
			if l, ok := places[s.LocationID]; ok {
				fmt.Fprintf(&b, "%d. %s [[%s]]: %s\n", s.Position, l.Name, l.ID, s.Note)
			}
		}
	}

	b.WriteString("\n## Etika berkunjung per kategori\n")
	for _, key := range []string{"tradisi", "pura", "masjid", "gereja", "vihara", "klenteng", "kuliner"} {
		fmt.Fprintf(&b, "- %s: %s\n", key, etiquetteNotes[key])
	}

	u.knowledge = b.String()
	u.places = places
	u.builtAt = now
	return u.knowledge, u.places, nil
}

var locationMarker = regexp.MustCompile(`\s*\[\[([0-9a-fA-F-]{36})\]\]`)

func (u *assistantUsecase) Ask(ctx context.Context, question string, history []AssistantMessage) (*AssistantAnswer, error) {
	if u.client == nil && !u.openai.enabled() {
		return nil, ErrAssistantDisabled
	}
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, errors.New("pertanyaan tidak boleh kosong")
	}
	maxChars := DefaultMaxInputChars
	if u.openai.enabled() {
		maxChars = u.openai.settings.Current().MaxInputChars
	}
	if len([]rune(question)) > maxChars {
		return nil, fmt.Errorf("pertanyaan maksimal %d karakter", maxChars)
	}

	knowledge, places, err := u.refresh()
	if err != nil {
		return nil, err
	}

	weekday := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	now := time.Now()
	dateNote := fmt.Sprintf("Tanggal hari ini: %s, %s.", weekday[now.Weekday()], now.Format("2006-01-02"))

	if u.client == nil {
		return u.askOpenAI(ctx, knowledge, dateNote, places, question, history)
	}

	params := anthropic.BetaMessageNewParams{
		Model:     assistantModel,
		MaxTokens: assistantMaxTokens,
		// Instructions and data first (cached), today's date after the
		// breakpoint so the daily change does not invalidate the cache.
		System: []anthropic.BetaTextBlockParam{
			{Text: assistantInstructions},
			{Text: knowledge, CacheControl: anthropic.NewBetaCacheControlEphemeralParam()},
			{Text: dateNote},
		},
		Messages:     buildMessages(history, question),
		OutputConfig: anthropic.BetaOutputConfigParam{Effort: anthropic.BetaOutputConfigEffortLow},
		// If a safety classifier declines, retry on the default fallback model.
		Fallbacks: anthropic.BetaFallbacksParamUnion{OfDefault: constant.ValueOf[constant.Default]()},
		Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaServerSideFallback2026_07_01},
	}

	resp, err := u.client.Beta.Messages.New(ctx, params)
	if err != nil {
		return nil, err
	}
	log.Printf("[TanyaLoka] tokens in=%d cache_read=%d out=%d",
		resp.Usage.InputTokens, resp.Usage.CacheReadInputTokens, resp.Usage.OutputTokens)

	if resp.StopReason == anthropic.BetaStopReasonRefusal {
		return &AssistantAnswer{Answer: refusalAnswer, Locations: []entity.Location{}}, nil
	}

	var text strings.Builder
	for _, block := range resp.Content {
		if t, ok := block.AsAny().(anthropic.BetaTextBlock); ok {
			text.WriteString(t.Text)
		}
	}
	return parseAnswer(text.String(), places), nil
}

// askOpenAI answers with OpenAI. Instructions and data come first and stay
// the same for an hour, so OpenAI's automatic prompt caching can reuse them;
// today's date comes after.
func (u *assistantUsecase) askOpenAI(ctx context.Context, knowledge, dateNote string, places map[string]entity.Location, question string, history []AssistantMessage) (*AssistantAnswer, error) {
	msgs := []openAIMessage{
		{Role: "system", Content: assistantInstructions + "\n\n" + knowledge + "\n\n" + dateNote},
	}
	for _, m := range trimHistory(history) {
		msgs = append(msgs, openAIMessage{Role: m.Role, Content: m.Content})
	}
	msgs = append(msgs, openAIMessage{Role: "user", Content: question})

	text, err := u.openai.chat(ctx, "TanyaLoka", msgs, 0.3, "", nil)
	if errors.Is(err, errOpenAIRefused) {
		return &AssistantAnswer{Answer: refusalAnswer, Locations: []entity.Location{}}, nil
	}
	if err != nil {
		return nil, err
	}
	return parseAnswer(text, places), nil
}

// trimHistory keeps the last turns, starting with a user turn and
// alternating roles, with long messages cut. Both providers use it.
func trimHistory(history []AssistantMessage) []AssistantMessage {
	if len(history) > maxHistoryMessages {
		history = history[len(history)-maxHistoryMessages:]
	}
	var out []AssistantMessage
	expectUser := true
	for _, m := range history {
		content := strings.TrimSpace(m.Content)
		if content == "" {
			continue
		}
		if r := []rune(content); len(r) > maxMessageChars {
			content = string(r[:maxMessageChars])
		}
		switch {
		case m.Role == "user" && expectUser:
			out = append(out, AssistantMessage{Role: "user", Content: content})
			expectUser = false
		case m.Role == "assistant" && !expectUser:
			out = append(out, AssistantMessage{Role: "assistant", Content: content})
			expectUser = true
		}
	}
	// History must end with an assistant turn before the new question.
	if !expectUser {
		out = out[:len(out)-1]
	}
	return out
}

// buildMessages keeps the last turns, starting with a user turn and
// alternating roles as the API requires, then appends the new question.
func buildMessages(history []AssistantMessage, question string) []anthropic.BetaMessageParam {
	var msgs []anthropic.BetaMessageParam
	for _, m := range trimHistory(history) {
		if m.Role == "user" {
			msgs = append(msgs, anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(m.Content)))
		} else {
			msgs = append(msgs, anthropic.BetaMessageParam{OfAssistant: &anthropic.BetaAssistantMessageParam{
				Content: []anthropic.BetaContentBlockParamUnion{anthropic.NewBetaTextBlock(m.Content)},
			}})
		}
	}
	return append(msgs, anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(question)))
}

// parseAnswer strips [[id]] markers and returns the referenced locations in
// order of first mention. Unknown ids are dropped.
func parseAnswer(text string, places map[string]entity.Location) *AssistantAnswer {
	seen := map[string]bool{}
	locations := []entity.Location{}
	for _, m := range locationMarker.FindAllStringSubmatch(text, -1) {
		id := strings.ToLower(m[1])
		if l, ok := places[id]; ok && !seen[id] {
			seen[id] = true
			locations = append(locations, l)
		}
	}
	answer := strings.TrimSpace(locationMarker.ReplaceAllString(text, ""))
	return &AssistantAnswer{Answer: answer, Locations: locations}
}
