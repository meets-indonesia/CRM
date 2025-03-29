package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/kevinnaserwan/crm-be/services/feedback/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/feedback/infrastructure/filestore"
	"github.com/skip2/go-qrcode"
)

var (
	ErrQRFeedbackNotFound = errors.New("qr feedback not found")
	ErrInvalidQRCode      = errors.New("invalid qr code")
)

type QRFeedbackRepository interface {
	Create(ctx context.Context, qrFeedback *entity.QRFeedback) error
	FindByID(ctx context.Context, id uint) (*entity.QRFeedback, error)
	FindByQRCode(ctx context.Context, qrCode string) (*entity.QRFeedback, error)
	ListAll(ctx context.Context, page, limit int) ([]entity.QRFeedback, int64, error)
	Delete(ctx context.Context, id uint) error
}

type QRFeedbackUsecase interface {
	CreateQRFeedback(ctx context.Context, req entity.CreateQRFeedbackRequest) (*entity.QRFeedbackResponse, error)
	GetQRFeedback(ctx context.Context, id uint) (*entity.QRFeedbackResponse, error)
	ListQRFeedback(ctx context.Context, page, limit int) (*entity.QRFeedbackListResponse, error)
	GenerateQRCodeImage(ctx context.Context, id uint) ([]byte, error)
	VerifyQRCode(ctx context.Context, code string) (*entity.QRFeedbackResponse, error)
}

type qrFeedbackUsecase struct {
	qrFeedbackRepo QRFeedbackRepository
	fileService    filestore.FileService
}

func NewQRFeedbackUsecase(qrFeedbackRepo QRFeedbackRepository, fileService filestore.FileService) QRFeedbackUsecase {
	return &qrFeedbackUsecase{
		qrFeedbackRepo: qrFeedbackRepo,
		fileService:    fileService,
	}
}

// generateRandomCode generates a random string for QR code
func generateRandomCode(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buffer)[:length], nil
}

// CreateQRFeedback creates a new QR code for feedback
func (u *qrFeedbackUsecase) CreateQRFeedback(ctx context.Context, req entity.CreateQRFeedbackRequest) (*entity.QRFeedbackResponse, error) {
	// Generate a random code
	code, err := generateRandomCode(12)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random code: %w", err)
	}

	// Create QR feedback entity
	qrFeedback := &entity.QRFeedback{
		QRCode:    code,
		Station:   req.Station,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := u.qrFeedbackRepo.Create(ctx, qrFeedback); err != nil {
		return nil, fmt.Errorf("failed to create QR feedback: %w", err)
	}

	// Prepare response
	response := &entity.QRFeedbackResponse{
		ID:        qrFeedback.ID,
		QRCode:    qrFeedback.QRCode,
		Station:   qrFeedback.Station,
		CreatedAt: qrFeedback.CreatedAt,
	}

	return response, nil
}

// GetQRFeedback gets a QR feedback by ID
func (u *qrFeedbackUsecase) GetQRFeedback(ctx context.Context, id uint) (*entity.QRFeedbackResponse, error) {
	qrFeedback, err := u.qrFeedbackRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if qrFeedback == nil {
		return nil, ErrQRFeedbackNotFound
	}

	response := &entity.QRFeedbackResponse{
		ID:        qrFeedback.ID,
		QRCode:    qrFeedback.QRCode,
		Station:   qrFeedback.Station,
		CreatedAt: qrFeedback.CreatedAt,
	}

	return response, nil
}

// ListQRFeedback lists all QR feedbacks with pagination
func (u *qrFeedbackUsecase) ListQRFeedback(ctx context.Context, page, limit int) (*entity.QRFeedbackListResponse, error) {
	qrFeedbacks, total, err := u.qrFeedbackRepo.ListAll(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list QR feedbacks: %w", err)
	}

	// Prepare response
	responses := make([]entity.QRFeedbackResponse, len(qrFeedbacks))
	for i, qrFeedback := range qrFeedbacks {
		responses[i] = entity.QRFeedbackResponse{
			ID:        qrFeedback.ID,
			QRCode:    qrFeedback.QRCode,
			CreatedAt: qrFeedback.CreatedAt,
		}
	}

	return &entity.QRFeedbackListResponse{
		Total:       total,
		Page:        page,
		Limit:       limit,
		QRFeedbacks: responses,
	}, nil
}

// GenerateQRCodeImage generates a QR code image for a feedback QR
func (u *qrFeedbackUsecase) GenerateQRCodeImage(ctx context.Context, id uint) ([]byte, error) {
	qrFeedback, err := u.qrFeedbackRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if qrFeedback == nil {
		return nil, ErrQRFeedbackNotFound
	}

	// Generate QR code with just the code value (no URL)
	qrCode, err := qrcode.Encode(qrFeedback.QRCode, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qrCode, nil
}

// VerifyQRCode verifies a QR code and returns the QR feedback information
func (u *qrFeedbackUsecase) VerifyQRCode(ctx context.Context, code string) (*entity.QRFeedbackResponse, error) {
	qrFeedback, err := u.qrFeedbackRepo.FindByQRCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if qrFeedback == nil {
		return nil, ErrInvalidQRCode
	}

	response := &entity.QRFeedbackResponse{
		ID:        qrFeedback.ID,
		QRCode:    qrFeedback.QRCode,
		Station:   qrFeedback.Station,
		CreatedAt: qrFeedback.CreatedAt,
	}

	return response, nil
}
