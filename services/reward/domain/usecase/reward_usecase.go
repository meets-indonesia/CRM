package usecase

import (
	"context"
	"errors"

	"github.com/kevinnaserwan/crm-be/services/reward/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/reward/domain/repository"
)

// Errors
var (
	ErrRewardNotFound     = errors.New("reward not found")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrInsufficientPoints = errors.New("insufficient points")
	ErrClaimNotFound      = errors.New("claim not found")
	ErrInvalidClaimStatus = errors.New("invalid claim status")
)

// RewardUsecase mendefinisikan operasi-operasi usecase untuk Reward
type RewardUsecase interface {
	// Reward operations
	CreateReward(ctx context.Context, req entity.CreateRewardRequest) (*entity.Reward, error)
	GetReward(ctx context.Context, id uint) (*entity.Reward, error)
	UpdateReward(ctx context.Context, id uint, req entity.UpdateRewardRequest) (*entity.Reward, error)
	DeleteReward(ctx context.Context, id uint) error
	ListRewards(ctx context.Context, activeOnly bool, page, limit int) (*entity.RewardListResponse, error)

	// Claim operations
	ClaimReward(ctx context.Context, userID uint, req entity.ClaimRewardRequest) (*entity.RewardClaim, error)
	GetClaim(ctx context.Context, id uint) (*entity.RewardClaim, error)
	UpdateClaimStatus(ctx context.Context, id uint, req entity.UpdateClaimStatusRequest) (*entity.RewardClaim, error)
	ListUserClaims(ctx context.Context, userID uint, page, limit int) (*entity.ClaimListResponse, error)
	ListAllClaims(ctx context.Context, page, limit int) (*entity.ClaimListResponse, error)
	ListClaimsByStatus(ctx context.Context, status entity.ClaimStatus, page, limit int) (*entity.ClaimListResponse, error)
}

type rewardUsecase struct {
	rewardRepo       repository.RewardRepository
	claimRepo        repository.ClaimRepository
	eventPublisher   repository.EventPublisher
	userPointService repository.UserPointService
}

// NewRewardUsecase membuat instance baru RewardUsecase
func NewRewardUsecase(
	rewardRepo repository.RewardRepository,
	claimRepo repository.ClaimRepository,
	eventPublisher repository.EventPublisher,
	userPointService repository.UserPointService,
) RewardUsecase {
	return &rewardUsecase{
		rewardRepo:       rewardRepo,
		claimRepo:        claimRepo,
		eventPublisher:   eventPublisher,
		userPointService: userPointService,
	}
}

// CreateReward membuat reward baru
func (u *rewardUsecase) CreateReward(ctx context.Context, req entity.CreateRewardRequest) (*entity.Reward, error) {
	reward := &entity.Reward{
		Name:        req.Name,
		Description: req.Description,
		PointCost:   req.PointCost,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}

	if err := u.rewardRepo.Create(ctx, reward); err != nil {
		return nil, err
	}

	return reward, nil
}

// GetReward mendapatkan reward berdasarkan ID
func (u *rewardUsecase) GetReward(ctx context.Context, id uint) (*entity.Reward, error) {
	reward, err := u.rewardRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if reward == nil {
		return nil, ErrRewardNotFound
	}

	return reward, nil
}

// UpdateReward memperbarui reward
func (u *rewardUsecase) UpdateReward(ctx context.Context, id uint, req entity.UpdateRewardRequest) (*entity.Reward, error) {
	reward, err := u.rewardRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if reward == nil {
		return nil, ErrRewardNotFound
	}

	// Update fields
	if req.Name != "" {
		reward.Name = req.Name
	}
	if req.Description != "" {
		reward.Description = req.Description
	}
	if req.PointCost > 0 {
		reward.PointCost = req.PointCost
	}
	if req.Stock >= 0 {
		reward.Stock = req.Stock
	}
	if req.ImageURL != "" {
		reward.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		reward.IsActive = *req.IsActive
	}

	if err := u.rewardRepo.Update(ctx, reward); err != nil {
		return nil, err
	}

	return reward, nil
}

// DeleteReward menghapus reward
func (u *rewardUsecase) DeleteReward(ctx context.Context, id uint) error {
	reward, err := u.rewardRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if reward == nil {
		return ErrRewardNotFound
	}

	return u.rewardRepo.Delete(ctx, id)
}

// ListRewards mendapatkan daftar reward
func (u *rewardUsecase) ListRewards(ctx context.Context, activeOnly bool, page, limit int) (*entity.RewardListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	rewards, total, err := u.rewardRepo.List(ctx, activeOnly, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.RewardListResponse{
		Rewards: rewards,
		Total:   total,
		Page:    page,
		Limit:   limit,
	}, nil
}

// ClaimReward mengklaim reward
func (u *rewardUsecase) ClaimReward(ctx context.Context, userID uint, req entity.ClaimRewardRequest) (*entity.RewardClaim, error) {
	// Ambil reward
	reward, err := u.rewardRepo.FindByID(ctx, req.RewardID)
	if err != nil {
		return nil, err
	}
	if reward == nil || !reward.IsActive {
		return nil, ErrRewardNotFound
	}

	// Cek stok
	stock, err := u.rewardRepo.CheckStock(ctx, req.RewardID)
	if err != nil {
		return nil, err
	}
	if stock <= 0 {
		return nil, ErrInsufficientStock
	}

	// Cek poin pengguna
	userPoints, err := u.userPointService.CheckUserPoints(ctx, userID)
	if err != nil {
		return nil, err
	}
	if userPoints < reward.PointCost {
		return nil, ErrInsufficientPoints
	}

	// Tentukan level reward berdasarkan PointCost reward
	rewardLevel := getRewardLevel(reward.PointCost)

	// Ambil semua klaim user
	existingClaims, _, err := u.claimRepo.ListByUserID(ctx, userID, 1, 1000)
	if err != nil {
		return nil, err
	}

	// Cek apakah user sudah klaim reward di level tersebut
	for _, claim := range existingClaims {
		if getRewardLevel(claim.PointCost) == rewardLevel {
			return nil, errors.New("Anda sudah mengklaim reward pada level ini")
		}
	}

	// Lanjut buat klaim
	claim := &entity.RewardClaim{
		UserID:    userID,
		RewardID:  req.RewardID,
		PointCost: reward.PointCost,
		Status:    entity.ClaimStatusPending,
		Reward:    *reward,
	}

	if err := u.claimRepo.Create(ctx, claim); err != nil {
		return nil, err
	}

	// Kurangi stok
	if err := u.rewardRepo.DecreaseStock(ctx, req.RewardID, 1); err != nil {
		return nil, err
	}

	// Kirim event (tidak blocking meski error)
	_ = u.eventPublisher.PublishRewardClaimed(claim)

	return claim, nil
}

// GetClaim mendapatkan klaim berdasarkan ID
func (u *rewardUsecase) GetClaim(ctx context.Context, id uint) (*entity.RewardClaim, error) {
	claim, err := u.claimRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if claim == nil {
		return nil, ErrClaimNotFound
	}

	return claim, nil
}

// UpdateClaimStatus memperbarui status klaim
func (u *rewardUsecase) UpdateClaimStatus(ctx context.Context, id uint, req entity.UpdateClaimStatusRequest) (*entity.RewardClaim, error) {
	claim, err := u.claimRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if claim == nil {
		return nil, ErrClaimNotFound
	}

	// Validate status transition
	if !isValidStatusTransition(claim.Status, req.Status) {
		return nil, ErrInvalidClaimStatus
	}

	// Handle stock changes if status is rejected or cancelled
	if (req.Status == entity.ClaimStatusRejected || req.Status == entity.ClaimStatusCancelled) &&
		(claim.Status == entity.ClaimStatusPending || claim.Status == entity.ClaimStatusApproved) {
		if err := u.rewardRepo.IncreaseStock(ctx, claim.RewardID, 1); err != nil {
			return nil, err
		}
	}

	// Update claim
	claim.Status = req.Status
	claim.Notes = req.Notes

	if err := u.claimRepo.Update(ctx, claim); err != nil {
		return nil, err
	}

	// Publish event
	if err := u.eventPublisher.PublishClaimStatusUpdated(claim); err != nil {
		// Log error but don't fail
	}

	return claim, nil
}

// ListUserClaims mendapatkan daftar klaim user
func (u *rewardUsecase) ListUserClaims(ctx context.Context, userID uint, page, limit int) (*entity.ClaimListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	claims, total, err := u.claimRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.ClaimListResponse{
		Claims: claims,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

// ListAllClaims mendapatkan semua klaim
func (u *rewardUsecase) ListAllClaims(ctx context.Context, page, limit int) (*entity.ClaimListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	claims, total, err := u.claimRepo.ListAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.ClaimListResponse{
		Claims: claims,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

// ListClaimsByStatus mendapatkan klaim berdasarkan status
func (u *rewardUsecase) ListClaimsByStatus(ctx context.Context, status entity.ClaimStatus, page, limit int) (*entity.ClaimListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	claims, total, err := u.claimRepo.ListByStatus(ctx, status, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.ClaimListResponse{
		Claims: claims,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

// isValidStatusTransition memeriksa apakah transisi status valid
func isValidStatusTransition(from, to entity.ClaimStatus) bool {
	switch from {
	case entity.ClaimStatusPending:
		return to == entity.ClaimStatusApproved || to == entity.ClaimStatusRejected || to == entity.ClaimStatusCancelled
	case entity.ClaimStatusApproved:
		return to == entity.ClaimStatusCancelled
	default:
		return false
	}
}

func getRewardLevel(point int) string {
	switch {
	case point >= 200:
		return "PLATINUM"
	case point >= 100:
		return "GOLD"
	case point >= 50:
		return "SILVER"
	default:
		return "BRONZE"
	}
}
