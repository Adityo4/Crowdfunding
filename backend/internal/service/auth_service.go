package service

import (
	"errors"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(req models.RegisterRequest) (*models.User, error)
	Login(req models.LoginRequest) (*models.AuthResponse, error)
	GetProfile(id uuid.UUID) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateProfile(id uuid.UUID, req models.UpdateProfileRequest) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(req models.RegisterRequest) (*models.User, error) {
	// Check if email already registered
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "user",
		PhoneNumber:  req.PhoneNumber,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	match := utils.CheckPasswordHash(req.Password, user.PasswordHash)
	if !match {
		return nil, errors.New("email atau password salah")
	}

	token, expiresAt, err := utils.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *user,
	}, nil
}

func (s *authService) GetProfile(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *authService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *authService) UpdateProfile(id uuid.UUID, req models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	user.FullName = req.FullName
	user.PhoneNumber = req.PhoneNumber
	user.UpdatedAt = time.Now()

	if req.Password != "" {
		hashed, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hashed
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
