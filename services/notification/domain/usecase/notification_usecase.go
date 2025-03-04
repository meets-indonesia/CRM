package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/kevinnaserwan/crm-be/services/notification/domain/entity"
	"github.com/kevinnaserwan/crm-be/services/notification/domain/repository"
)

// Errors
var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidUserID        = errors.New("invalid user ID")
	ErrSendingNotification  = errors.New("error sending notification")
)

// NotificationUsecase mendefinisikan operasi-operasi usecase untuk Notification
type NotificationUsecase interface {
	// Email operations
	SendEmail(ctx context.Context, req entity.CreateEmailRequest) (*entity.Notification, error)

	// Push notification operations
	SendPushNotification(ctx context.Context, req entity.CreatePushNotificationRequest) (*entity.Notification, error)

	// Get operations
	GetNotification(ctx context.Context, id uint) (*entity.Notification, error)
	ListUserNotifications(ctx context.Context, userID uint, page, limit int) (*entity.NotificationListResponse, error)

	// Processing
	ProcessPendingNotifications(ctx context.Context) error

	// Event handling
	HandleArticleCreated(articleID uint, title, summary string, authorID uint, publishedAt time.Time) error
	HandleFeedbackResponded(feedbackID uint, userID uint, title, response string) error
	HandleRewardClaimed(claimID uint, userID uint, rewardID uint, status entity.NotificationStatus) error
}

type notificationUsecase struct {
	notificationRepo repository.NotificationRepository
	emailSender      repository.EmailSender
	pushSender       repository.PushNotificationSender
}

// NewNotificationUsecase membuat instance baru NotificationUsecase
func NewNotificationUsecase(
	notificationRepo repository.NotificationRepository,
	emailSender repository.EmailSender,
	pushSender repository.PushNotificationSender,
) NotificationUsecase {
	return &notificationUsecase{
		notificationRepo: notificationRepo,
		emailSender:      emailSender,
		pushSender:       pushSender,
	}
}

// SendEmail mengirim email
func (u *notificationUsecase) SendEmail(ctx context.Context, req entity.CreateEmailRequest) (*entity.Notification, error) {
	if req.UserID == 0 {
		return nil, ErrInvalidUserID
	}

	// Buat notifikasi
	notification := &entity.Notification{
		UserID:  req.UserID,
		Type:    entity.TypeEmail,
		Title:   req.Subject,
		Content: req.Body,
		Status:  entity.StatusPending,
		Data:    req.EmailTo, // Simpan alamat email di field data
	}

	// Simpan notifikasi ke database
	if err := u.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	// Coba kirim email langsung
	err := u.emailSender.SendEmail(req.EmailTo, req.Subject, req.Body, req.IsHTML)
	now := time.Now()

	if err != nil {
		// Update status jika gagal
		notification.Status = entity.StatusFailed
		notification.Error = err.Error()
		notification.SentAt = &now
		u.notificationRepo.Update(ctx, notification)
		return notification, ErrSendingNotification
	}

	// Update status jika berhasil
	notification.Status = entity.StatusSent
	notification.SentAt = &now
	if err := u.notificationRepo.Update(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

// SendPushNotification mengirim push notification
func (u *notificationUsecase) SendPushNotification(ctx context.Context, req entity.CreatePushNotificationRequest) (*entity.Notification, error) {
	if req.UserID == 0 {
		return nil, ErrInvalidUserID
	}

	// Buat notifikasi
	dataJSON, _ := json.Marshal(req.Data)

	notification := &entity.Notification{
		UserID:  req.UserID,
		Type:    entity.TypePushNotification,
		Title:   req.Title,
		Content: req.Message,
		Status:  entity.StatusPending,
		Data:    string(dataJSON),
	}

	// Simpan notifikasi ke database
	if err := u.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	// Coba kirim push notification langsung
	err := u.pushSender.SendPushNotification(req.UserID, req.Title, req.Message, req.Data)
	now := time.Now()

	if err != nil {
		// Update status jika gagal
		notification.Status = entity.StatusFailed
		notification.Error = err.Error()
		notification.SentAt = &now
		u.notificationRepo.Update(ctx, notification)
		return notification, ErrSendingNotification
	}

	// Update status jika berhasil
	notification.Status = entity.StatusSent
	notification.SentAt = &now
	if err := u.notificationRepo.Update(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

// GetNotification mendapatkan notifikasi berdasarkan ID
func (u *notificationUsecase) GetNotification(ctx context.Context, id uint) (*entity.Notification, error) {
	notification, err := u.notificationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if notification == nil {
		return nil, ErrNotificationNotFound
	}

	return notification, nil
}

// ListUserNotifications mendapatkan daftar notifikasi user
func (u *notificationUsecase) ListUserNotifications(ctx context.Context, userID uint, page, limit int) (*entity.NotificationListResponse, error) {
	if userID == 0 {
		return nil, ErrInvalidUserID
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	notifications, total, err := u.notificationRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	return &entity.NotificationListResponse{
		Notifications: notifications,
		Total:         total,
		Page:          page,
		Limit:         limit,
	}, nil
}

// ProcessPendingNotifications memproses notifikasi yang belum terkirim
func (u *notificationUsecase) ProcessPendingNotifications(ctx context.Context) error {
	// Ambil notifikasi yang pending
	notifications, err := u.notificationRepo.ListPending(ctx, 50)
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		var err error
		now := time.Now()

		if notification.Type == entity.TypeEmail {
			// Kirim email
			err = u.emailSender.SendEmail(notification.Data, notification.Title, notification.Content, true)
		} else if notification.Type == entity.TypePushNotification {
			// Parse data JSON
			var data map[string]interface{}
			if notification.Data != "" {
				if err := json.Unmarshal([]byte(notification.Data), &data); err != nil {
					// Log error dan lanjutkan
					data = map[string]interface{}{}
				}
			}

			// Kirim push notification
			err = u.pushSender.SendPushNotification(notification.UserID, notification.Title, notification.Content, data)
		}

		if err != nil {
			// Update status jika gagal
			notification.Status = entity.StatusFailed
			notification.Error = err.Error()
		} else {
			// Update status jika berhasil
			notification.Status = entity.StatusSent
		}

		notification.SentAt = &now
		u.notificationRepo.Update(ctx, &notification)
	}

	return nil
}

// HandleArticleCreated menangani event artikel baru dibuat
func (u *notificationUsecase) HandleArticleCreated(articleID uint, title, summary string, authorID uint, publishedAt time.Time) error {
	// Buat permintaan push notification
	req := entity.CreatePushNotificationRequest{
		UserID:  0, // Akan diisi per user
		Title:   "Artikel Baru: " + title,
		Message: summary,
		Data: map[string]interface{}{
			"article_id": articleID,
			"type":       "article_created",
		},
	}

	// TODO: Ambil semua user yang aktif dan kirim ke masing-masing
	// Untuk demo, kita bisa hardcode beberapa user ID
	userIDs := []uint{1, 2, 3} // Contoh user ID

	for _, userID := range userIDs {
		req.UserID = userID
		_, _ = u.SendPushNotification(context.Background(), req)
	}

	return nil
}

// HandleFeedbackResponded menangani event feedback direspons
func (u *notificationUsecase) HandleFeedbackResponded(feedbackID uint, userID uint, title, response string) error {
	// Buat permintaan push notification
	pushReq := entity.CreatePushNotificationRequest{
		UserID:  userID,
		Title:   "Respons untuk Feedback Anda",
		Message: "Feedback '" + title + "' telah direspons: " + response,
		Data: map[string]interface{}{
			"feedback_id": feedbackID,
			"type":        "feedback_responded",
		},
	}

	// Kirim push notification
	_, _ = u.SendPushNotification(context.Background(), pushReq)

	// Buat permintaan email
	// TODO: Ambil email dari user service
	emailReq := entity.CreateEmailRequest{
		UserID:  userID,
		EmailTo: "user@example.com", // Hardcoded untuk demo
		Subject: "Respons untuk Feedback Anda",
		Body:    "<h1>Feedback Anda telah direspons</h1><p>Feedback: " + title + "</p><p>Respons: " + response + "</p>",
		IsHTML:  true,
	}

	// Kirim email
	_, _ = u.SendEmail(context.Background(), emailReq)

	return nil
}

// HandleRewardClaimed menangani event hadiah diklaim
func (u *notificationUsecase) HandleRewardClaimed(claimID uint, userID uint, rewardID uint, status entity.NotificationStatus) error {
	// Implementasi mirip dengan HandleFeedbackResponded
	// ...

	return nil
}
