package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

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
	jwtSecret     string
	jwtDuration   time.Duration
	messageBroker *messagebroker.RabbitMQ
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	jwtSecret string,
	jwtDuration time.Duration,
	messageBroker *messagebroker.RabbitMQ,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:      userRepo,
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

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := util.GenerateJWT(user.ID, uc.jwtSecret, uc.jwtDuration)
	if err != nil {
		return "", err
	}

	return token, nil
}
