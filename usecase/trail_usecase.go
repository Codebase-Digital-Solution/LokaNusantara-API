package usecase

import (
	"api/entity"
	"api/repository"
	"errors"
	"regexp"
	"strings"
)

// TrailInput is the editable part of a trail, as sent by the admin panel.
type TrailInput struct {
	Slug     string
	Title    string
	Subtitle string
	Story    string
	Color    string
	ImageURL string
	Stops    []TrailStopInput
}

type TrailStopInput struct {
	LocationID string
	Note       string
}

type TrailUsecase interface {
	GetTrails() ([]entity.Trail, error)
	GetTrail(id uint) (*entity.Trail, error)
	CreateTrail(in TrailInput) (*entity.Trail, error)
	UpdateTrail(id uint, in TrailInput) (*entity.Trail, error)
	DeleteTrail(id uint) error
}

type trailUsecase struct {
	trailRepo    repository.TrailRepository
	locationRepo repository.LocationRepository
}

func NewTrailUsecase(trailRepo repository.TrailRepository, locationRepo repository.LocationRepository) TrailUsecase {
	return &trailUsecase{trailRepo: trailRepo, locationRepo: locationRepo}
}

// GetTrails returns every trail, dropping stops whose location no longer
// exists so clients never receive a stop without a place.
func (u *trailUsecase) GetTrails() ([]entity.Trail, error) {
	trails, err := u.trailRepo.GetAll()
	if err != nil {
		return nil, err
	}
	for i := range trails {
		stops := trails[i].Stops[:0]
		for _, s := range trails[i].Stops {
			if s.Location != nil {
				stops = append(stops, s)
			}
		}
		trails[i].Stops = stops
	}
	return trails, nil
}

func (u *trailUsecase) GetTrail(id uint) (*entity.Trail, error) { return u.trailRepo.GetByID(id) }

var (
	slugPattern  = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	colorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
)

func (u *trailUsecase) validate(in *TrailInput) error {
	in.Slug = strings.TrimSpace(in.Slug)
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		return errors.New("trail title is required")
	}
	if !slugPattern.MatchString(in.Slug) {
		return errors.New("slug must be lowercase letters, numbers and dashes")
	}
	if !validImageURL(in.ImageURL) {
		return errImageURL
	}
	if in.Color != "" && !colorPattern.MatchString(in.Color) {
		return errors.New("color must be a hex value like #A87618")
	}
	if len(in.Stops) == 0 {
		return errors.New("a trail needs at least one stop")
	}
	for _, s := range in.Stops {
		if _, err := u.locationRepo.GetByID(s.LocationID); err != nil {
			return errors.New("stop location not found: " + s.LocationID)
		}
	}
	return nil
}

func (u *trailUsecase) save(trail *entity.Trail, in TrailInput) (*entity.Trail, error) {
	trail.Slug = in.Slug
	trail.Title = in.Title
	trail.Subtitle = in.Subtitle
	trail.Story = in.Story
	trail.Color = in.Color
	trail.ImageURL = in.ImageURL
	trail.Stops = make([]entity.TrailStop, 0, len(in.Stops))
	for _, s := range in.Stops {
		trail.Stops = append(trail.Stops, entity.TrailStop{LocationID: s.LocationID, Note: s.Note})
	}
	if err := u.trailRepo.Save(trail); err != nil {
		return nil, errors.New("failed to save trail (slug must be unique)")
	}
	return u.trailRepo.GetByID(trail.ID)
}

func (u *trailUsecase) CreateTrail(in TrailInput) (*entity.Trail, error) {
	if err := u.validate(&in); err != nil {
		return nil, err
	}
	return u.save(&entity.Trail{}, in)
}

func (u *trailUsecase) UpdateTrail(id uint, in TrailInput) (*entity.Trail, error) {
	if err := u.validate(&in); err != nil {
		return nil, err
	}
	trail, err := u.trailRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return u.save(trail, in)
}

func (u *trailUsecase) DeleteTrail(id uint) error {
	if _, err := u.trailRepo.GetByID(id); err != nil {
		return err
	}
	return u.trailRepo.Delete(id)
}
