package proxy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

const (
	AuthTimestampHeader = "X-Auth-Timestamp"
	AuthSignatureHeader = "X-Auth-Signature"
	SecretKey           = "CRMSUMSEL2025@MEETSIDN" // Your secret key
)

// AddAuthHeaders adds the required authentication headers to the request
func AddAuthHeaders(req *http.Request) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := GenerateSignature(timestamp, SecretKey)

	req.Header.Set(AuthTimestampHeader, timestamp)
	req.Header.Set(AuthSignatureHeader, signature)
}

// GenerateSignature creates HMAC-SHA256 signature
func GenerateSignature(timestamp, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(timestamp))
	return hex.EncodeToString(h.Sum(nil))
}
