package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/domain/events"
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

	// Simpan user ke database
	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	// Publish user registered event setelah user berhasil dibuat
	event := events.UserRegisteredEvent{
		UserID:    user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	err = uc.messageBroker.PublishMessage(ctx, "auth.events", "user.registered", event)
	if err != nil {
		log.Printf("Failed to publish user registered event: %v", err)
	}

	return nil
}

// In auth service's usecase
func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Get user's role
	role, err := uc.roleRepo.FindByID(ctx, user.RoleID)
	if err != nil {
		return "", errors.New("role not found")
	}

	claims := util.Claims{
		UserID: user.ID,
		Role:   role.Name, // Include role name here
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(uc.jwtDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
