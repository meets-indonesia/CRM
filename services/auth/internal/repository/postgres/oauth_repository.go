// services/auth/internal/repository/postgres/oauth_repository.go
package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"gorm.io/gorm"
)

type oauthRepository struct {
	db *gorm.DB
}

func NewOAuthRepository(db *gorm.DB) *oauthRepository {
	return &oauthRepository{db: db}
}

func (r *oauthRepository) FindByProviderID(ctx context.Context, provider model.OAuthProvider, providerID string) (*model.OAuthAccount, error) {
	var account model.OAuthAccount
	err := r.db.WithContext(ctx).
		Where("provider = ? AND provider_id = ?", provider, providerID).
		First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *oauthRepository) Create(ctx context.Context, account *model.OAuthAccount) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *oauthRepository) Update(ctx context.Context, account *model.OAuthAccount) error {
	return r.db.WithContext(ctx).Save(account).Error
}

func (r *oauthRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.OAuthAccount, error) {
	var accounts []model.OAuthAccount
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
