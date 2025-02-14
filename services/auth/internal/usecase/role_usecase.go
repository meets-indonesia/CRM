package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/repository"
)

type RoleUseCase struct {
	roleRepo repository.RoleRepository
}

func NewRoleUseCase(roleRepo repository.RoleRepository) *RoleUseCase {
	return &RoleUseCase{
		roleRepo: roleRepo,
	}
}

func (uc *RoleUseCase) Create(ctx context.Context, role *model.Role) error {
	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}
	return uc.roleRepo.Create(ctx, role)
}
