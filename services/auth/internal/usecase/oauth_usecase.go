// services/auth/internal/usecase/oauth_usecase.go
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/repository"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/service"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/util"
)

type OAuthUseCase struct {
	userRepo     repository.UserRepository
	oauthRepo    repository.OAuthRepository
	roleRepo     repository.RoleRepository
	oauthService *service.OAuthService
	jwtSecret    string
	jwtDuration  time.Duration
}

func NewOAuthUseCase(
	userRepo repository.UserRepository,
	oauthRepo repository.OAuthRepository,
	roleRepo repository.RoleRepository,
	oauthService *service.OAuthService,
	jwtSecret string,
	jwtDuration time.Duration,
) *OAuthUseCase {
	return &OAuthUseCase{
		userRepo:     userRepo,
		oauthRepo:    oauthRepo,
		roleRepo:     roleRepo,
		oauthService: oauthService,
		jwtSecret:    jwtSecret,
		jwtDuration:  jwtDuration,
	}
}

func (uc *OAuthUseCase) GoogleSignIn(ctx context.Context, code, redirectURI string) (string, *model.User, error) {
	// Exchange code for tokens and user info
	userInfo, accessToken, refreshToken, expiresAt, err := uc.oauthService.ExchangeGoogleToken(ctx, code, redirectURI)
	if err != nil {
		return "", nil, err
	}

	// Check if OAuth account exists
	existingOAuth, err := uc.oauthRepo.FindByProviderID(ctx, model.GoogleProvider, userInfo.ID)
	if err == nil {
		// OAuth account exists, get user
		user, err := uc.userRepo.FindByID(ctx, existingOAuth.UserID)
		if err != nil {
			return "", nil, errors.New("user not found")
		}

		// Update OAuth tokens
		existingOAuth.AccessToken = accessToken
		existingOAuth.RefreshToken = refreshToken
		existingOAuth.ExpiresAt = expiresAt
		if err := uc.oauthRepo.Update(ctx, existingOAuth); err != nil {
			return "", nil, err
		}

		// Generate JWT
		token, err := uc.generateToken(user)
		if err != nil {
			return "", nil, err
		}

		return token, user, nil
	}

	// OAuth account doesn't exist, check if user with same email exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, userInfo.Email)
	if err == nil {
		// User exists, link OAuth account
		oauthAccount := &model.OAuthAccount{
			UserID:       existingUser.ID,
			Provider:     userInfo.Provider,
			ProviderID:   userInfo.ID,
			Email:        userInfo.Email,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    expiresAt,
		}

		if err := uc.oauthRepo.Create(ctx, oauthAccount); err != nil {
			return "", nil, err
		}

		// Generate JWT
		token, err := uc.generateToken(existingUser)
		if err != nil {
			return "", nil, err
		}

		return token, existingUser, nil
	}

	// User doesn't exist, create new user and OAuth account
	// Get default user role
	defaultRole, err := uc.getDefaultRole(ctx)
	if err != nil {
		return "", nil, err
	}

	newUser := &model.User{
		ID:             uuid.New(),
		FirstName:      userInfo.FirstName,
		LastName:       userInfo.LastName,
		Email:          userInfo.Email,
		ProfilePicture: userInfo.Picture,
		RoleID:         defaultRole.ID,
		Password:       "", // OAuth users don't have a password
	}

	// Create user
	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return "", nil, err
	}

	// Create OAuth account
	oauthAccount := &model.OAuthAccount{
		UserID:       newUser.ID,
		Provider:     userInfo.Provider,
		ProviderID:   userInfo.ID,
		Email:        userInfo.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	if err := uc.oauthRepo.Create(ctx, oauthAccount); err != nil {
		return "", nil, err
	}

	// Generate JWT
	token, err := uc.generateToken(newUser)
	if err != nil {
		return "", nil, err
	}

	return token, newUser, nil
}

func (uc *OAuthUseCase) getDefaultRole(ctx context.Context) (*model.Role, error) {
	// Find role by name "user"
	defaultRoleName := "user"

	// We need to add a FindByName method to the roleRepo interface
	role, err := uc.roleRepo.FindByName(ctx, defaultRoleName)
	if err != nil {
		return nil, fmt.Errorf("default role 'user' not found: %w", err)
	}

	return role, nil
}

func (uc *OAuthUseCase) generateToken(user *model.User) (string, error) {
	// Get user's role
	role, err := uc.roleRepo.FindByID(context.Background(), user.RoleID)
	if err != nil {
		return "", errors.New("role not found")
	}

	claims := util.Claims{
		UserID: user.ID,
		Role:   role.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(uc.jwtDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
