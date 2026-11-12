package usecase

import (
	"api/entity"
	"api/repository"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(name, email, password string) (*entity.User, error)
	Login(email, password string) (*entity.User, error)
	// LoginWithGoogle signs in (or registers) the owner of a verified Google
	// account. An existing account with the same email is linked to it.
	LoginWithGoogle(googleID, email, name, avatarURL string) (*entity.User, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
	}
}

func (u *authUsecase) Register(name, email, password string) (*entity.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email, and password are required")
	}

	// Check if user already exists
	existingUser, _ := u.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Generate UUID
	userID := generateUUID()

	user := &entity.User{
		ID:       userID,
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     entity.RoleUser,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (u *authUsecase) Login(email, password string) (*entity.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (u *authUsecase) LoginWithGoogle(googleID, email, name, avatarURL string) (*entity.User, error) {
	if googleID == "" || email == "" {
		return nil, errors.New("google account has no id or email")
	}

	if user, err := u.userRepo.GetByGoogleID(googleID); err == nil {
		return user, nil
	}

	// Same email already registered with a password: link the Google account.
	if user, err := u.userRepo.GetByEmail(email); err == nil {
		user.GoogleID = &googleID
		if user.AvatarURL == "" {
			user.AvatarURL = avatarURL
		}
		if err := u.userRepo.Update(user); err != nil {
			return nil, errors.New("failed to link google account")
		}
		return user, nil
	}

	if name == "" {
		name = email
	}
	// Google-only accounts get a random password nobody knows, so password
	// login stays impossible until a reset flow exists.
	hashed, err := bcrypt.GenerateFromPassword([]byte(generateUUID()+generateUUID()), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to create user")
	}
	user := &entity.User{
		ID:        generateUUID(),
		Name:      name,
		Email:     email,
		Password:  string(hashed),
		Role:      entity.RoleUser,
		GoogleID:  &googleID,
		AvatarURL: avatarURL,
	}
	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}
	return user, nil
}

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
