package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/user/internal/domain/model"
	"github.com/kevinnaserwan/crm-be/services/user/internal/domain/repository"
)

type PointUseCase struct {
	userPointsRepo   repository.UserPointsRepository
	pointHistoryRepo repository.PointHistoryRepository
	pointResetRepo   repository.PointResetHistoryRepository
}

func NewPointUseCase(
	userPointsRepo repository.UserPointsRepository,
	pointHistoryRepo repository.PointHistoryRepository,
	pointResetRepo repository.PointResetHistoryRepository,
) *PointUseCase {
	return &PointUseCase{
		userPointsRepo:   userPointsRepo,
		pointHistoryRepo: pointHistoryRepo,
		pointResetRepo:   pointResetRepo,
	}
}

// InitializeIfNotExists creates a user points record if it doesn't exist yet
func (uc *PointUseCase) InitializeIfNotExists(ctx context.Context, userID uuid.UUID) error {
	_, err := uc.userPointsRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Printf("Creating initial points record for user %s", userID)
		now := time.Now()
		nextReset := now.AddDate(1, 0, 0) // Add 1 year

		userPoints := &model.UserPoints{
			UserID:           userID,
			TotalPoints:      0,
			RegistrationDate: now,
			LastResetDate:    now,
			NextResetDate:    nextReset,
		}

		return uc.userPointsRepo.Create(ctx, userPoints)
	}
	return nil
}

// ProcessFeedbackPoint handles point awarding for feedback
func (uc *PointUseCase) ProcessFeedbackPoint(ctx context.Context, userID, feedbackID uuid.UUID) error {
	log.Printf("Processing points for feedback %s, user %s", feedbackID, userID)

	// Ensure user has points record
	if err := uc.InitializeIfNotExists(ctx, userID); err != nil {
		return err
	}

	// Get user points
	userPoints, err := uc.userPointsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if points need to be reset
	if time.Now().After(userPoints.NextResetDate) {
		if err := uc.resetPoints(ctx, userPoints); err != nil {
			return err
		}
		// Reload user points after reset
		userPoints, err = uc.userPointsRepo.GetByUserID(ctx, userID)
		if err != nil {
			return err
		}
	}

	// Check daily point limit
	if userPoints.LastPointEarnedDate != nil {
		lastEarnedDate := userPoints.LastPointEarnedDate.UTC().Truncate(24 * time.Hour)
		today := time.Now().UTC().Truncate(24 * time.Hour)
		if lastEarnedDate.Equal(today) {
			log.Printf("Daily point limit reached for user %s", userID)
			return errors.New("daily point limit reached")
		}
	}

	// Record point history
	pointHistory := &model.PointHistory{
		UserID:       userID,
		PointsEarned: model.PointsPerFeedback,
		DateEarned:   time.Now(),
		FeedbackID:   feedbackID,
	}
	if err := uc.pointHistoryRepo.Create(ctx, pointHistory); err != nil {
		return err
	}

	// Update user points
	newTotal := userPoints.TotalPoints + model.PointsPerFeedback
	if err := uc.userPointsRepo.UpdatePoints(ctx, userID, newTotal); err != nil {
		return err
	}

	// Update last point earned date
	now := time.Now()
	err = uc.userPointsRepo.UpdateLastPointEarned(ctx, userID, now)
	if err != nil {
		return err
	}

	log.Printf("Points processed successfully: User %s now has %d points", userID, newTotal)
	return nil
}

// resetPoints handles the yearly point reset
func (uc *PointUseCase) resetPoints(ctx context.Context, userPoints *model.UserPoints) error {
	log.Printf("Resetting points for user %s", userPoints.UserID)

	now := time.Now()
	nextReset := now.AddDate(1, 0, 0)

	// Record reset history
	resetHistory := &model.PointResetHistory{
		UserID:            userPoints.UserID,
		ResetDate:         now,
		PointsBeforeReset: userPoints.TotalPoints,
		NextResetDate:     nextReset,
	}
	if err := uc.pointResetRepo.Create(ctx, resetHistory); err != nil {
		return err
	}

	// Reset points
	if err := uc.userPointsRepo.UpdatePoints(ctx, userPoints.UserID, 0); err != nil {
		return err
	}

	// Update reset dates
	return uc.userPointsRepo.UpdateResetDates(ctx, userPoints.UserID, now, nextReset)
}

// GetUserPointsAndLevel gets current points and determines user level
func (uc *PointUseCase) GetUserPointsAndLevel(ctx context.Context, userID uuid.UUID) (*model.UserPoints, string, error) {
	userPoints, err := uc.userPointsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, "", err
	}

	// Determine level based on points
	level := "Bronze"
	points := userPoints.TotalPoints

	switch {
	case points >= model.PlatinumThreshold:
		level = "Platinum"
	case points >= model.GoldThreshold:
		level = "Gold"
	case points >= model.SilverThreshold:
		level = "Silver"
	}

	return userPoints, level, nil
}

// GetPointHistory retrieves point earning history
func (uc *PointUseCase) GetPointHistory(ctx context.Context, userID uuid.UUID) ([]model.PointHistory, error) {
	return uc.pointHistoryRepo.GetByUserID(ctx, userID)
}

// GetAllUserPoints gets all users' points (admin function)
func (uc *PointUseCase) GetAllUserPoints(ctx context.Context) ([]model.UserPoints, error) {
	return uc.userPointsRepo.List(ctx)
}
