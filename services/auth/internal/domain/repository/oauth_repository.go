// services/auth/internal/domain/repository/oauth_repository.go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
)

type OAuthRepository interface {
	// Find OAuth account by provider and provider ID
	FindByProviderID(ctx context.Context, provider model.OAuthProvider, providerID string) (*model.OAuthAccount, error)

	// Create new OAuth account
	Create(ctx context.Context, account *model.OAuthAccount) error

	// Update OAuth account (tokens, expiry)
	Update(ctx context.Context, account *model.OAuthAccount) error

	// Find all OAuth accounts for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.OAuthAccount, error)
}
