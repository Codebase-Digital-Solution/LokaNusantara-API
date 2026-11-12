package usecase

import (
	"api/entity"
	"api/repository"
	"encoding/json"
	"errors"
	"time"
)

const maxSavedPlansPerUser = 100

// SavedPlanView is one saved plan as returned to the app.
type SavedPlanView struct {
	ID        uint             `json:"id"`
	Location  *entity.Location `json:"location"`
	Plan      TripPlan         `json:"plan"`
	CreatedAt time.Time        `json:"created_at"`
}

type SavedPlanUsecase interface {
	List(userID string) ([]SavedPlanView, error)
	Save(userID, locationID string, plan TripPlan) (*SavedPlanView, error)
	Delete(userID string, id uint) error
}

// ErrSavedPlanNotFound is returned when the plan does not exist or belongs
// to someone else.
var ErrSavedPlanNotFound = errors.New("rencana tidak ditemukan")

type savedPlanUsecase struct {
	repo         repository.SavedPlanRepository
	locationRepo repository.LocationRepository
}

func NewSavedPlanUsecase(repo repository.SavedPlanRepository, locationRepo repository.LocationRepository) SavedPlanUsecase {
	return &savedPlanUsecase{repo: repo, locationRepo: locationRepo}
}

// recomputeTotals keeps a stored plan's totals consistent with its costs,
// whatever the client sent.
func recomputeTotals(p *TripPlan) {
	p.TotalEstimate = 0
	for i := range p.Costs {
		if p.Costs[i].Amount < 0 {
			p.Costs[i].Amount = 0
		}
		p.TotalEstimate += p.Costs[i].Amount
	}
	if p.People < 1 {
		p.People = 1
	}
	p.PerPersonEstimate = (p.TotalEstimate + int64(p.People) - 1) / int64(p.People)
	p.WithinBudget = p.TotalEstimate <= p.Budget
}

func (u *savedPlanUsecase) List(userID string) ([]SavedPlanView, error) {
	rows, err := u.repo.GetByUser(userID)
	if err != nil {
		return nil, err
	}
	out := make([]SavedPlanView, 0, len(rows))
	for _, r := range rows {
		var plan TripPlan
		if json.Unmarshal([]byte(r.PlanJSON), &plan) != nil {
			continue // unreadable row; skip rather than fail the whole list
		}
		out = append(out, SavedPlanView{ID: r.ID, Location: r.Location, Plan: plan, CreatedAt: r.CreatedAt})
	}
	return out, nil
}

func (u *savedPlanUsecase) Save(userID, locationID string, plan TripPlan) (*SavedPlanView, error) {
	location, err := u.locationRepo.GetByID(locationID)
	if err != nil {
		return nil, inputErr("lokasi tidak ditemukan")
	}
	if plan.Days < 1 || plan.People < 1 || plan.Budget < 0 || len(plan.Costs) > 50 || len(plan.Itinerary) > 100 || len(plan.Tips) > 20 {
		return nil, inputErr("isi rencana tidak valid")
	}
	if n, err := u.repo.CountByUser(userID); err != nil {
		return nil, err
	} else if n >= maxSavedPlansPerUser {
		return nil, inputErr("rencana tersimpan sudah mencapai batas 100, hapus beberapa dulu")
	}

	recomputeTotals(&plan)
	data, err := json.Marshal(plan)
	if err != nil {
		return nil, err
	}
	row := &entity.SavedPlan{UserID: userID, LocationID: locationID, PlanJSON: string(data)}
	if err := u.repo.Create(row); err != nil {
		return nil, err
	}
	location.Reviews = nil
	return &SavedPlanView{ID: row.ID, Location: location, Plan: plan, CreatedAt: row.CreatedAt}, nil
}

func (u *savedPlanUsecase) Delete(userID string, id uint) error {
	ok, err := u.repo.Delete(userID, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrSavedPlanNotFound
	}
	return nil
}
