package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/kevinnaserwan/crm-be/services/auth/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/auth/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// Errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidOTP         = errors.New("invalid or expired OTP")
	ErrEmailRequired      = errors.New("email is required")
	ErrEmailExists        = errors.New("email already exists")
)

// AuthUsecase mendefinisikan operasi-operasi usecase untuk Auth
type AuthUsecase interface {
	RegisterAdmin(ctx context.Context, req entity.RegisterRequest) (*entity.User, error)
	LoginAdmin(ctx context.Context, req entity.LoginRequest) (*entity.LoginResponse, error)
	ResetPasswordRequest(ctx context.Context, req entity.ResetPasswordRequest) error
	VerifyOTP(ctx context.Context, req entity.VerifyOTPRequest) error
	LoginCustomer(ctx context.Context, req entity.LoginRequest) (*entity.LoginResponse, error)
	LoginWithGoogle(ctx context.Context, req entity.GoogleLoginRequest) (*entity.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*entity.User, error)
}

type authUsecase struct {
	userRepo           repository.UserRepository
	tokenGenerator     repository.TokenGenerator
	otpGenerator       repository.OTPGenerator
	eventPublisher     repository.EventPublisher
	emailSender        repository.EmailSender
	googleAuthProvider repository.GoogleAuthProvider
}

// NewAuthUsecase membuat instance baru AuthUsecase
func NewAuthUsecase(
	userRepo repository.UserRepository,
	tokenGenerator repository.TokenGenerator,
	otpGenerator repository.OTPGenerator,
	eventPublisher repository.EventPublisher,
	emailSender repository.EmailSender,
	googleAuthProvider repository.GoogleAuthProvider,
) AuthUsecase {
	return &authUsecase{
		userRepo:           userRepo,
		tokenGenerator:     tokenGenerator,
		otpGenerator:       otpGenerator,
		eventPublisher:     eventPublisher,
		emailSender:        emailSender,
		googleAuthProvider: googleAuthProvider,
	}
}

// RegisterAdmin mendaftarkan admin baru
func (a *authUsecase) RegisterAdmin(ctx context.Context, req entity.RegisterRequest) (*entity.User, error) {
	// Check if email already exists
	existingUser, err := a.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entity.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     entity.RoleAdmin,
	}

	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Publish event
	if err := a.eventPublisher.PublishUserCreated(user.ID, user.Email, user.Role); err != nil {
		// Log error but don't fail
	}

	return user, nil
}

// LoginAdmin melakukan login admin
func (a *authUsecase) LoginAdmin(ctx context.Context, req entity.LoginRequest) (*entity.LoginResponse, error) {
	user, err := a.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is admin
	if user.Role != entity.RoleAdmin {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, _, err := a.tokenGenerator.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Publish event
	if err := a.eventPublisher.PublishUserLoggedIn(user.ID, user.Email, user.Role); err != nil {
		// Log error but don't fail
	}

	return &entity.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// ResetPasswordRequest mengirim OTP untuk reset password
func (a *authUsecase) ResetPasswordRequest(ctx context.Context, req entity.ResetPasswordRequest) error {
	user, err := a.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	// Generate OTP
	otpCode, err := a.otpGenerator.GenerateOTP()
	if err != nil {
		return err
	}

	// Create OTP record
	otp := &entity.OTP{
		UserID:    user.ID,
		Code:      otpCode,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := a.userRepo.CreateOTP(ctx, otp); err != nil {
		return err
	}

	// Send OTP email
	if err := a.emailSender.SendOTP(user.Email, user.Name, otpCode); err != nil {
		return err
	}

	// Publish event
	if err := a.eventPublisher.PublishPasswordReset(user.ID, user.Email); err != nil {
		// Log error but don't fail
	}

	return nil
}

// VerifyOTP memverifikasi OTP dan mengatur password baru
func (a *authUsecase) VerifyOTP(ctx context.Context, req entity.VerifyOTPRequest) error {
	otp, err := a.userRepo.FindOTPByCode(ctx, req.Code)
	if err != nil || otp == nil {
		return ErrInvalidOTP
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		return ErrInvalidOTP
	}

	// Get user
	user, err := a.userRepo.FindByID(ctx, otp.UserID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user password
	user.Password = string(hashedPassword)
	if err := a.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Delete OTP
	if err := a.userRepo.DeleteOTP(ctx, otp.ID); err != nil {
		// Log error but don't fail
	}

	return nil
}

// LoginCustomer melakukan login customer
func (a *authUsecase) LoginCustomer(ctx context.Context, req entity.LoginRequest) (*entity.LoginResponse, error) {
	user, err := a.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is customer
	if user.Role != entity.RoleCustomer {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, _, err := a.tokenGenerator.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Publish event
	if err := a.eventPublisher.PublishUserLoggedIn(user.ID, user.Email, user.Role); err != nil {
		// Log error but don't fail
	}

	return &entity.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// LoginWithGoogle melakukan login dengan Google OAuth
func (a *authUsecase) LoginWithGoogle(ctx context.Context, req entity.GoogleLoginRequest) (*entity.LoginResponse, error) {
	// Verify Google ID token
	googleID, email, name, err := a.googleAuthProvider.VerifyIDToken(req.Token)
	if err != nil {
		return nil, err
	}

	// Check if user exists by Google ID
	user, err := a.userRepo.FindByGoogleID(ctx, googleID)

	// If user doesn't exist, check by email
	if err != nil || user == nil {
		user, err = a.userRepo.FindByEmail(ctx, email)

		// If user doesn't exist by email either, create new user
		if err != nil || user == nil {
			user = &entity.User{
				Email:    email,
				Name:     name,
				GoogleID: googleID,
				Role:     entity.RoleCustomer,
			}

			if err := a.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}

			// Publish event
			if err := a.eventPublisher.PublishUserCreated(user.ID, user.Email, user.Role); err != nil {
				// Log error but don't fail
			}
		} else {
			// Update existing user with Google ID
			user.GoogleID = googleID
			if err := a.userRepo.Update(ctx, user); err != nil {
				return nil, err
			}
		}
	}

	// Generate token
	token, _, err := a.tokenGenerator.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Publish event
	if err := a.eventPublisher.PublishUserLoggedIn(user.ID, user.Email, user.Role); err != nil {
		// Log error but don't fail
	}

	return &entity.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// ValidateToken memvalidasi token JWT
func (a *authUsecase) ValidateToken(ctx context.Context, token string) (*entity.User, error) {
	userID, _, _, err := a.tokenGenerator.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	user, err := a.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
