package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
)

type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Role, error)
	FindByName(ctx context.Context, name string) (*model.Role, error)
}
