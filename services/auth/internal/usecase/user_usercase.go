package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/repository"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}
