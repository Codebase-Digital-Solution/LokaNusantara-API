package usecase

import (
	"api/entity"
	"api/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrPlannerDisabled is returned when no OpenAI API key is configured (admin
// settings or OPENAI_API_KEY).
var ErrPlannerDisabled = errors.New("fitur rencana kunjungan belum dikonfigurasi di server")

const (
	plannerNearbyKm  = 15
	plannerMaxNearby = 8
)

// PlanInputError is a problem with the visitor's input (shown to them as is),
// as opposed to a failure of the AI service.
type PlanInputError struct{ Msg string }

func (e *PlanInputError) Error() string { return e.Msg }

func inputErr(msg string) error { return &PlanInputError{Msg: msg} }

// PlanRequest is what the visitor fills in on the plan form.
type PlanRequest struct {
	LocationID string `json:"location_id"`
	Budget     int64  `json:"budget"` // total for the whole group, in rupiah
	People     int    `json:"people"`
	Days       int    `json:"days"`
	Notes      string `json:"notes"`
}

// PlanCost is one line of the estimated spending, for the whole group.
type PlanCost struct {
	Category string `json:"category"`
	Item     string `json:"item"`
	Amount   int64  `json:"amount"`
	Note     string `json:"note"`
}

// PlanStep is one activity in the itinerary. LocationID is set only when the
// step happens at a place in our database.
type PlanStep struct {
	Day        int              `json:"day"`
	Time       string           `json:"time"`
	Activity   string           `json:"activity"`
	LocationID string           `json:"location_id,omitempty"`
	Location   *entity.Location `json:"location,omitempty"`
}

// TripPlan is the generated plan. Totals are computed by the server from
// Costs, so they always add up.
type TripPlan struct {
	Summary           string     `json:"summary"`
	Budget            int64      `json:"budget"`
	People            int        `json:"people"`
	Days              int        `json:"days"`
	TotalEstimate     int64      `json:"total_estimate"`
	PerPersonEstimate int64      `json:"per_person_estimate"`
	WithinBudget      bool       `json:"within_budget"`
	BudgetNote        string     `json:"budget_note"`
	Costs             []PlanCost `json:"costs"`
	Itinerary         []PlanStep `json:"itinerary"`
	Tips              []string   `json:"tips"`
}

type PlannerUsecase interface {
	Plan(ctx context.Context, req PlanRequest) (*TripPlan, error)
}

type plannerUsecase struct {
	ai           *openAIClient
	locationRepo repository.LocationRepository
	eventRepo    repository.EventRepository
}

// NewPlannerUsecase returns a planner backed by the OpenAI Chat Completions
// API, configured through settings (disabled while no key is set).
func NewPlannerUsecase(settings AISettingsSource, locationRepo repository.LocationRepository, eventRepo repository.EventRepository) PlannerUsecase {
	return &plannerUsecase{
		ai:           newOpenAIClient(settings),
		locationRepo: locationRepo,
		eventRepo:    eventRepo,
	}
}

func validatePlan(req *PlanRequest, maxNotes int) error {
	req.Notes = strings.TrimSpace(req.Notes)
	switch {
	case req.LocationID == "":
		return inputErr("lokasi wajib dipilih")
	case req.Budget < 10_000:
		return inputErr("budget minimal Rp10.000")
	case req.Budget > 1_000_000_000:
		return inputErr("budget terlalu besar")
	case req.People < 1 || req.People > 50:
		return inputErr("jumlah orang harus 1 sampai 50")
	case req.Days < 1 || req.Days > 7:
		return inputErr("lama kunjungan harus 1 sampai 7 hari")
	case len([]rune(req.Notes)) > maxNotes:
		return inputErr(fmt.Sprintf("catatan maksimal %d karakter", maxNotes))
	}
	return nil
}

const plannerInstructions = `Kamu adalah perencana perjalanan budaya di aplikasi LokaNusantara (Bali).
Susun rencana kunjungan yang realistis, hormat pada adat, dan sesuai budget.

Aturan:
- Tulis dalam Bahasa Indonesia yang ramah dan ringkas.
- Semua biaya dalam rupiah, bilangan bulat, untuk SELURUH rombongan (bukan per orang).
- Biaya adalah perkiraan wajar di Bali. Jangan mengklaim harga resmi; tulis "perkiraan" di catatan bila tidak pasti. Untuk pura atau rumah ibadah, pakai "donasi/punia" dan sewa kamen bila relevan, bukan tiket, kecuali memang umum ada tiket.
- Usahakan total biaya tidak melebihi budget. Jika budget tidak cukup, susun versi paling hemat dan jelaskan kekurangannya di budget_note.
- Itinerary hanya boleh memakai tempat dari daftar "Tempat di database" bila menyebut tempat spesifik, dan isi location_id dengan id persis dari daftar. Untuk aktivitas umum (makan, perjalanan, istirahat) kosongkan location_id.
- Kegiatan utama harus di tempat tujuan. Tempat sekitar hanya tambahan bila waktu dan budget cukup.
- Jika ada agenda budaya di tempat itu, sebutkan dan ingatkan bahwa tanggal bertanda perkiraan bisa berubah.
- Jangan menulis angka total atau biaya per orang di summary maupun budget_note; total dihitung aplikasi dari daftar costs. Cukup jelaskan apakah budget longgar, pas, atau kurang, dan apa yang dihemat.
- Sesuaikan jam dengan jenis tempat: pasar malam hanya sore atau malam, rumah ibadah hindari jam ibadah utama, upacara mengikuti jadwal adat.
- tips: 2 sampai 5 tips praktis (etika, pakaian, waktu terbaik, uang tunai).
- category biaya pilih salah satu: transportasi, tiket, donasi, makan, sewa, pemandu, oleh-oleh, lainnya.`

// planSchema forces the model to answer with exactly this JSON shape.
var planSchema = map[string]any{
	"type":                 "object",
	"additionalProperties": false,
	"required":             []string{"summary", "budget_note", "costs", "itinerary", "tips"},
	"properties": map[string]any{
		"summary":     map[string]any{"type": "string"},
		"budget_note": map[string]any{"type": "string"},
		"costs": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"required":             []string{"category", "item", "amount", "note"},
				"properties": map[string]any{
					"category": map[string]any{"type": "string", "enum": []string{"transportasi", "tiket", "donasi", "makan", "sewa", "pemandu", "oleh-oleh", "lainnya"}},
					"item":     map[string]any{"type": "string"},
					"amount":   map[string]any{"type": "integer"},
					"note":     map[string]any{"type": "string"},
				},
			},
		},
		"itinerary": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"required":             []string{"day", "time", "activity", "location_id"},
				"properties": map[string]any{
					"day":         map[string]any{"type": "integer"},
					"time":        map[string]any{"type": "string"},
					"activity":    map[string]any{"type": "string"},
					"location_id": map[string]any{"type": "string"},
				},
			},
		},
		"tips": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
	},
}

func (u *plannerUsecase) Plan(ctx context.Context, req PlanRequest) (*TripPlan, error) {
	if !u.ai.enabled() {
		return nil, ErrPlannerDisabled
	}
	if err := validatePlan(&req, u.ai.settings.Current().MaxInputChars); err != nil {
		return nil, err
	}

	target, err := u.locationRepo.GetByID(req.LocationID)
	if err != nil {
		return nil, inputErr("lokasi tidak ditemukan")
	}
	nearby, _ := u.locationRepo.GetNearby(target.Latitude, target.Longitude, plannerNearbyKm)
	places := map[string]*entity.Location{target.ID: target}
	var others []entity.Location
	for i := range nearby {
		if nearby[i].ID == target.ID || len(others) >= plannerMaxNearby {
			continue
		}
		others = append(others, nearby[i])
		places[nearby[i].ID] = &others[len(others)-1]
	}

	now := time.Now()
	events, _ := u.eventRepo.GetBetween(now, now.AddDate(0, 4, 0), 50, 0)

	prompt := buildPlanPrompt(req, target, others, events)
	raw, err := u.ai.chat(ctx, "Planner", []openAIMessage{
		{Role: "system", Content: plannerInstructions},
		{Role: "user", Content: prompt},
	}, 0.4, "trip_plan", planSchema)
	if err != nil {
		return nil, fmt.Errorf("planner: %w", err)
	}

	var out struct {
		Summary    string     `json:"summary"`
		BudgetNote string     `json:"budget_note"`
		Costs      []PlanCost `json:"costs"`
		Itinerary  []PlanStep `json:"itinerary"`
		Tips       []string   `json:"tips"`
	}
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, fmt.Errorf("planner: invalid model output: %w", err)
	}

	plan := &TripPlan{
		Summary:    out.Summary,
		Budget:     req.Budget,
		People:     req.People,
		Days:       req.Days,
		BudgetNote: out.BudgetNote,
		Tips:       out.Tips,
	}
	for _, c := range out.Costs {
		if c.Amount < 0 {
			c.Amount = 0
		}
		plan.Costs = append(plan.Costs, c)
		plan.TotalEstimate += c.Amount
	}
	plan.PerPersonEstimate = (plan.TotalEstimate + int64(req.People) - 1) / int64(req.People)
	plan.WithinBudget = plan.TotalEstimate <= req.Budget

	for _, s := range out.Itinerary {
		// Only keep ids we actually sent; drop anything the model made up.
		if loc, ok := places[s.LocationID]; ok {
			l := *loc
			l.Reviews = nil
			s.Location = &l
		} else {
			s.LocationID = ""
		}
		plan.Itinerary = append(plan.Itinerary, s)
	}
	if plan.Costs == nil {
		plan.Costs = []PlanCost{}
	}
	if plan.Itinerary == nil {
		plan.Itinerary = []PlanStep{}
	}
	if plan.Tips == nil {
		plan.Tips = []string{}
	}
	return plan, nil
}

func describePlace(b *strings.Builder, l *entity.Location) {
	fmt.Fprintf(b, "- id=%s | %s | kategori: %s", l.ID, l.Name, l.Category.Name)
	if l.Address != "" {
		fmt.Fprintf(b, " | alamat: %s", l.Address)
	}
	b.WriteString("\n")
}

func truncateRunes(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= n {
		return string(r)
	}
	return string(r[:n]) + "..."
}

func buildPlanPrompt(req PlanRequest, target *entity.Location, others []entity.Location, events []entity.Event) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Permintaan:\n- Budget total rombongan: Rp%d\n- Jumlah orang: %d\n- Lama kunjungan: %d hari\n", req.Budget, req.People, req.Days)
	if req.Notes != "" {
		fmt.Fprintf(&b, "- Catatan pengunjung: %s\n", req.Notes)
	}

	b.WriteString("\nTempat tujuan utama:\n")
	describePlace(&b, target)
	if d := truncateRunes(target.Description, 600); d != "" {
		fmt.Fprintf(&b, "  Deskripsi: %s\n", d)
	}
	if h := truncateRunes(target.History, 600); h != "" {
		fmt.Fprintf(&b, "  Sejarah singkat: %s\n", h)
	}

	b.WriteString("\nTempat di database (boleh dipakai di itinerary):\n")
	describePlace(&b, target)
	for i := range others {
		describePlace(&b, &others[i])
	}

	var related []string
	for _, e := range events {
		if e.LocationID != nil && *e.LocationID == target.ID {
			when := e.StartDate.Format("2 Jan 2006")
			if e.IsEstimate {
				when += " (perkiraan)"
			}
			related = append(related, fmt.Sprintf("- %s, %s", e.Name, when))
		}
	}
	if len(related) > 0 {
		b.WriteString("\nAgenda budaya di tempat tujuan dalam 4 bulan ke depan:\n")
		b.WriteString(strings.Join(related, "\n"))
		b.WriteString("\n")
	}
	return b.String()
}
