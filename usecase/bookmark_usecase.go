package usecase

import (
	"api/entity"
	"api/repository"
	"errors"
)

type BookmarkUsecase interface {
	GetBookmarks(userID string) ([]entity.Location, error)
	AddBookmark(userID, locationID string) error
	RemoveBookmark(userID, locationID string) error
}

type bookmarkUsecase struct {
	bookmarkRepo repository.BookmarkRepository
	locationRepo repository.LocationRepository
}

func NewBookmarkUsecase(bookmarkRepo repository.BookmarkRepository, locationRepo repository.LocationRepository) BookmarkUsecase {
	return &bookmarkUsecase{bookmarkRepo: bookmarkRepo, locationRepo: locationRepo}
}

func (u *bookmarkUsecase) GetBookmarks(userID string) ([]entity.Location, error) {
	return u.bookmarkRepo.GetLocations(userID)
}

func (u *bookmarkUsecase) AddBookmark(userID, locationID string) error {
	if _, err := u.locationRepo.GetByID(locationID); err != nil {
		return errors.New("location not found")
	}
	return u.bookmarkRepo.Add(userID, locationID)
}

func (u *bookmarkUsecase) RemoveBookmark(userID, locationID string) error {
	return u.bookmarkRepo.Remove(userID, locationID)
}
