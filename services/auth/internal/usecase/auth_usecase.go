package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/repository"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/util"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo    repository.UserRepository
	jwtSecret   string
	jwtDuration time.Duration
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtSecret string, jwtDuration time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		jwtDuration: jwtDuration,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, user *model.User) error {
	// Parse role_id from string to UUID
	if user.RoleID == uuid.Nil {
		return fmt.Errorf("role_id is required")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Set other fields if they're empty
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return uc.userRepo.Create(ctx, user)
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := util.GenerateJWT(user.ID, uc.jwtSecret, uc.jwtDuration)
	if err != nil {
		return "", err
	}

	return token, nil
}
