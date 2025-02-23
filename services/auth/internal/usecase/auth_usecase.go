package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/repository"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/infrastructure/messagebroker"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/util"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo      repository.UserRepository
	roleRepo      repository.RoleRepository
	jwtSecret     string
	jwtDuration   time.Duration
	messageBroker *messagebroker.RabbitMQ
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	jwtSecret string,
	jwtDuration time.Duration,
	messageBroker *messagebroker.RabbitMQ,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:      userRepo,
		roleRepo:      roleRepo,
		jwtSecret:     jwtSecret,
		jwtDuration:   jwtDuration,
		messageBroker: messageBroker,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, user *model.User) error {
	existingUser, err := uc.userRepo.FindByEmail(ctx, user.Email)
	if err == nil { // User found with email
		// If registering with GoogleID and existing user has same GoogleID
		if user.GoogleID != nil && existingUser.GoogleID != nil && *user.GoogleID == *existingUser.GoogleID {
			// Return special error that can be handled for auto-login
			return &AutoLoginError{
				Message:  "User already exists with this Google ID",
				Email:    user.Email,
				GoogleID: *user.GoogleID,
			}
		}
		return errors.New("email already registered")
	}

	// Check if GoogleID already exists (if provided)
	if user.GoogleID != nil {
		existingUserByGoogleID, err := uc.userRepo.FindByGoogleID(ctx, *user.GoogleID)
		if err == nil { // User found with GoogleID
			return &AutoLoginError{
				Message:  "User already exists with this Google ID",
				Email:    existingUserByGoogleID.Email,
				GoogleID: *user.GoogleID,
			}
		}
	}

	// Validation
	if user.Password == nil && user.GoogleID == nil {
		return errors.New("either password or google_id must be provided")
	}
	if user.Password != nil && user.GoogleID != nil {
		return errors.New("cannot provide both password and google_id")
	}

	// Hash password if it exists
	if user.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		hashedStr := string(hashedPassword)
		user.Password = &hashedStr
	}

	return uc.userRepo.Create(ctx, user)
}

type AutoLoginError struct {
	Message  string
	Email    string
	GoogleID string
}

func (e *AutoLoginError) Error() string {
	return e.Message
}

// In auth service's usecase
func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if user.Password == nil {
		return "", errors.New("this account uses Google authentication")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Use generateToken instead of creating claims directly
	return uc.generateToken(user)
}

// LoginOAuth for mobile users (email + googleID)
func (uc *AuthUseCase) LoginOAuth(ctx context.Context, email, googleID string) (string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if user.GoogleID == nil {
		return "", errors.New("this account uses password authentication")
	}

	if *user.GoogleID != googleID {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := util.GenerateJWT(user.ID, user.Email, user.Role.Name, uc.jwtSecret, uc.jwtDuration)
	if err != nil {
		return "", err
	}

	return token, nil
}

// internal/usecase/auth_usecase.go
func (uc *AuthUseCase) generateToken(user *model.User) (string, error) {
	// Debug print untuk memeriksa nilai
	fmt.Printf("Debug - User data: ID=%v, Email=%v, Role=%v\n",
		user.ID, user.Email, user.Role.Name)

	token, err := util.GenerateJWT(
		user.ID,
		user.Email, // Pastikan urutan parameter sesuai
		user.Role.Name,
		uc.jwtSecret,
		uc.jwtDuration,
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (uc *AuthUseCase) GetDefaultRole(ctx context.Context) (*model.Role, error) {
	// Find role by name "user"
	defaultRole, err := uc.roleRepo.FindByName(ctx, "user")
	if err != nil {
		return nil, fmt.Errorf("default role 'user' not found: %w", err)
	}
	return defaultRole, nil
}
