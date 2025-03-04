package push

import (
	"fmt"

	"firebase.google.com/go/v4/messaging"
	"github.com/kevinnaserwan/crm-be/services/notification/config"
)

// FirebasePushSender implements PushNotificationSender using Firebase Cloud Messaging
type FirebasePushSender struct {
	client *messaging.Client
}

// NewFirebasePushSender creates a new FirebasePushSender
func NewFirebasePushSender(config config.PushConfig) (*FirebasePushSender, error) {
	// Untuk tujuan demo, kita gunakan implementasi sederhana yang hanya melakukan log
	// Untuk implementasi sebenarnya, kita akan menggunakan Firebase Cloud Messaging
	return &FirebasePushSender{}, nil
}

// SendPushNotification sends a push notification
func (s *FirebasePushSender) SendPushNotification(userID uint, title, message string, data map[string]interface{}) error {
	// Untuk tujuan demo, kita hanya melakukan log
	// Pada implementasi sebenarnya, kita akan mengirim ke FCM

	fmt.Printf("Sending push notification to user %d: %s - %s\n", userID, title, message)
	fmt.Printf("Data: %v\n", data)

	return nil
}

// Implementasi FCM sebenarnya akan seperti ini:
/*
func NewFirebasePushSender(config config.PushConfig) (*FirebasePushSender, error) {
    opt := option.WithCredentialsJSON([]byte(`{
        "type": "service_account",
        "project_id": "` + config.ProjectID + `",
        "private_key_id": "...",
        "private_key": "...",
        "client_email": "...",
        "client_id": "...",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_x509_cert_url": "..."
    }`))

    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        return nil, err
    }

    client, err := app.Messaging(context.Background())
    if err != nil {
        return nil, err
    }

    return &FirebasePushSender{
        client: client,
    }, nil
}

func (s *FirebasePushSender) SendPushNotification(userID uint, title, message string, data map[string]interface{}) error {
    // Dalam implementasi sebenarnya, kita perlu mendapatkan FCM token untuk userID
    // dari database atau service lain

    // Untuk tujuan demo, kita asumsikan kita punya token
    token := "user_fcm_token_" + fmt.Sprintf("%d", userID)

    // Convert data to map[string]string
    stringData := make(map[string]string)
    for k, v := range data {
        stringData[k] = fmt.Sprintf("%v", v)
    }

    // Create message
    msg := &messaging.Message{
        Notification: &messaging.Notification{
            Title: title,
            Body:  message,
        },
        Data:  stringData,
        Token: token,
    }

    // Send message
    _, err := s.client.Send(context.Background(), msg)
    return err
}
*/
