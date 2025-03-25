package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/kevinnaserwan/crm-be/services/user/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/user/domain/repository"
)

// Errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidRole  = errors.New("invalid role")
)

// UserUsecase mendefinisikan operasi-operasi usecase untuk User
type UserUsecase interface {
	// User operations
	GetUser(ctx context.Context, id uint) (*entity.UserResponse, error)
	UpdateUser(ctx context.Context, id uint, req entity.UpdateUserRequest) (*entity.UserResponse, error)
	ListAdmins(ctx context.Context, page, limit int) (*entity.UserListResponse, error)
	ListCustomers(ctx context.Context, page, limit int) (*entity.UserListResponse, error)

	// Point operations
	GetCustomerPoints(ctx context.Context, id uint) (*entity.PointBalance, error)
	GetPointTransactions(ctx context.Context, userID uint, page, limit int) (*entity.PointTransactionListResponse, error)

	// Event processing
	ProcessUserCreated(userID uint, email string, role entity.Role) error
	ProcessFeedbackCreated(userID uint, feedbackID uint) error
	ProcessRewardClaimed(userID uint, rewardID uint, points int) error
}

type userUsecase struct {
	userRepo  repository.UserRepository
	pointRepo repository.PointRepository
}

// NewUserUsecase membuat instance baru UserUsecase
func NewUserUsecase(userRepo repository.UserRepository, pointRepo repository.PointRepository) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		pointRepo: pointRepo,
	}
}

// GetUser mendapatkan informasi user berdasarkan ID
func (u *userUsecase) GetUser(ctx context.Context, id uint) (*entity.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	response := &entity.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Phone:     user.Phone,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Jika user adalah customer, tambahkan informasi poin
	if user.Role == entity.RoleCustomer {
		points, err := u.pointRepo.GetPointBalance(ctx, user.ID)
		if err == nil {
			response.Points = points
		}
	}

	return response, nil
}

// UpdateUser memperbarui informasi user
func (u *userUsecase) UpdateUser(ctx context.Context, id uint, req entity.UpdateUserRequest) (*entity.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Address != "" {
		user.Address = req.Address
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return u.GetUser(ctx, id)
}

// ListAdmins mendapatkan daftar admin
func (u *userUsecase) ListAdmins(ctx context.Context, page, limit int) (*entity.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := u.userRepo.List(ctx, entity.RoleAdmin, page, limit)
	if err != nil {
		return nil, err
	}

	response := &entity.UserListResponse{
		Users: make([]entity.UserResponse, len(users)),
		Total: total,
		Page:  page,
		Limit: limit,
	}

	for i, user := range users {
		response.Users[i] = entity.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			Phone:     user.Phone,
			Address:   user.Address,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return response, nil
}

// ListCustomers mendapatkan daftar customer
func (u *userUsecase) ListCustomers(ctx context.Context, page, limit int) (*entity.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := u.userRepo.List(ctx, entity.RoleCustomer, page, limit)
	if err != nil {
		return nil, err
	}

	response := &entity.UserListResponse{
		Users: make([]entity.UserResponse, len(users)),
		Total: total,
		Page:  page,
		Limit: limit,
	}

	for i, user := range users {
		response.Users[i] = entity.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			Phone:     user.Phone,
			Address:   user.Address,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		// Tambahkan informasi poin
		points, err := u.pointRepo.GetPointBalance(ctx, user.ID)
		if err == nil {
			response.Users[i].Points = points
		}
	}

	return response, nil
}

// GetCustomerPoints mendapatkan informasi poin customer
func (u *userUsecase) GetCustomerPoints(ctx context.Context, id uint) (*entity.PointBalance, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if user.Role != entity.RoleCustomer {
		return nil, ErrInvalidRole
	}

	points, err := u.pointRepo.GetPointBalance(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	currentLevel, err := u.pointRepo.GetUserLevel(ctx, points)
	if err != nil {
		return nil, err
	}

	pointBalance := &entity.PointBalance{
		UserID: user.ID,
		Total:  points,
		Level:  currentLevel.Name,
	}

	// Check if there's a next level
	nextLevel, err := u.pointRepo.GetNextLevel(ctx, currentLevel)
	if err == nil && nextLevel != nil {
		pointBalance.NextLevel = nextLevel.Name
		pointBalance.ToNext = nextLevel.MinPoints - points
	}

	return pointBalance, nil
}

// GetPointTransactions mendapatkan daftar transaksi poin
func (u *userUsecase) GetPointTransactions(ctx context.Context, userID uint, page, limit int) (*entity.PointTransactionListResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if user.Role != entity.RoleCustomer {
		return nil, ErrInvalidRole
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	transactions, total, err := u.pointRepo.GetPointTransactions(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.PointTransactionListResponse{
		Transactions: transactions,
		Total:        total,
		Page:         page,
		Limit:        limit,
	}, nil
}

// ProcessUserCreated memproses event pembuatan user
func (u *userUsecase) ProcessUserCreated(userID uint, email string, role entity.Role) error {
	// Check if user already exists
	user, err := u.userRepo.FindByID(context.Background(), userID)
	if err != nil {
		return err
	}

	if user != nil {
		// User already exists, no need to create
		return nil
	}

	// Create user in this service
	newUser := &entity.User{
		ID:    userID,
		Email: email,
		Name:  email, // Default name is email, can be updated later
		Role:  role,
	}

	return u.userRepo.Create(context.Background(), newUser)
}

// ProcessFeedbackCreated memproses event pembuatan feedback
func (u *userUsecase) ProcessFeedbackCreated(userID uint, feedbackID uint) error {
	// Check if user exists
	user, err := u.userRepo.FindByID(context.Background(), userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Check if user already got points today
	today := time.Now().Truncate(24 * time.Hour)
	dailyPoints, err := u.pointRepo.GetDailyPoints(context.Background(), userID, today)
	if err != nil {
		return err
	}

	// If user already got points today, don't add more
	if dailyPoints > 0 {
		return nil
	}

	// Add 1 point for feedback
	transaction := &entity.PointTransaction{
		UserID:      userID,
		Amount:      1,
		Type:        "feedback",
		Description: "Points earned for submitting feedback",
		CreatedAt:   time.Now(),
	}

	return u.pointRepo.AddPoints(context.Background(), transaction)
}

// ProcessRewardClaimed memproses event klaim hadiah
func (u *userUsecase) ProcessRewardClaimed(userID uint, rewardID uint, points int) error {
	// Check if user exists
	user, err := u.userRepo.FindByID(context.Background(), userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	transaction := &entity.PointTransaction{
		UserID:      userID,
		Amount:      0, // Negative for deduction
		Type:        "reward",
		Description: "Points used for claiming reward",
		CreatedAt:   time.Now(),
	}

	return u.pointRepo.AddPoints(context.Background(), transaction)
}
