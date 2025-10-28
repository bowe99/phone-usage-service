package service

import (
	"context"
	"errors"
	"fmt"

	dto "github.com/bowe99/phone-usage-service/internal/application/dtos"
	"github.com/bowe99/phone-usage-service/internal/domain/model"
	"github.com/bowe99/phone-usage-service/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserService struct {
	userRepo repository.UserRepository
}

func SetupUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*model.UserResponse, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to ahsh password: %w", err)
	}

	user := &model.User{
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ToResponse(), nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, req dto.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
        user.FirstName = *req.FirstName
    }

    if req.LastName != nil {
        user.LastName = *req.LastName
    }

    if req.Email != nil {
        user.Email = *req.Email
    }

	//Not in the requirements from the project but assumed it was a typo as an update for the password would be a necessary case
    if req.Password != nil {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to ahsh password: %w", err)
		}
		user.Password = string(hashedPassword)
    }

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user.ToResponse(), nil
}
