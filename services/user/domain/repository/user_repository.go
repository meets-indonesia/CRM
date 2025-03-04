package repository

import (
	"context"

	"github.com/kevinnaserwan/crm-be/services/user/domain/entity"
)

// UserRepository mendefinisikan operasi-operasi repository untuk User
type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, role entity.Role, page, limit int) ([]entity.User, int64, error)
}

// EventSubscriber mendefinisikan operasi-operasi untuk berlangganan event
type EventSubscriber interface {
	SubscribeToUserEvents() error
	Close() error
}

// UserEventProcessor mendefinisikan operasi-operasi untuk memproses event user
type UserEventProcessor interface {
	ProcessUserCreated(userID uint, email string, role entity.Role) error
}
