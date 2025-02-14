package repository

import (
	"context"

	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
)

type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
}
