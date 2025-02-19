// services/auth/internal/delivery/http/oauth_handler.go
package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/delivery/http/models"
	"github.com/kevinnaserwan/crm-be/services/auth/internal/usecase"
)

type OAuthHandler struct {
	oauthUseCase *usecase.OAuthUseCase
}

func NewOAuthHandler(oauthUseCase *usecase.OAuthUseCase) *OAuthHandler {
	return &OAuthHandler{
		oauthUseCase: oauthUseCase,
	}
}

// Handler for meta-based approach
func (h *OAuthHandler) Handle(c *gin.Context) {
	var baseRequest models.BaseRequest
	if err := c.ShouldBindJSON(&baseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch baseRequest.Meta.Action {
	case "google_sign_in":
		h.handleGoogleSignIn(c, baseRequest.Data)
	case "apple_sign_in":
		h.handleAppleSignIn(c, baseRequest.Data)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"})
	}
}

// Handle Google sign in
func (h *OAuthHandler) handleGoogleSignIn(c *gin.Context, data interface{}) {
	var request struct {
		Code        string `json:"code" binding:"required"`
		RedirectURI string `json:"redirect_uri" binding:"required"`
	}

	if err := convertData(data, &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.oauthUseCase.GoogleSignIn(c.Request.Context(), request.Code, request.RedirectURI)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"profile_picture": user.ProfilePicture,
		},
	})
}

// Handle Apple sign in
func (h *OAuthHandler) handleAppleSignIn(c *gin.Context, data interface{}) {
	// TODO: Implement Apple sign in
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Apple sign in not implemented yet"})
}

// GoogleCallback handles the redirect from Google OAuth
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	// Get authorization code from query parameter
	code := c.Query("code")
	if code == "" {
		// If there's an error or no code, handle it
		errorMsg := c.Query("error")
		if errorMsg == "" {
			errorMsg = "No authorization code provided"
		}
		c.HTML(http.StatusBadRequest, "oauth_error.html", gin.H{
			"error": errorMsg,
		})
		return
	}

	// The redirect URI must match exactly what was used to get the code
	redirectURI := fmt.Sprintf("%s://%s%s", c.Request.URL.Scheme, c.Request.Host, c.Request.URL.Path)
	if c.Request.URL.Scheme == "" {
		// If scheme is not in the request URL (common in some setups), construct it
		proto := "http"
		if c.Request.TLS != nil {
			proto = "https"
		}
		redirectURI = fmt.Sprintf("%s://%s%s", proto, c.Request.Host, c.Request.URL.Path)
	}

	// Exchange code for token and user info
	token, user, err := h.oauthUseCase.GoogleSignIn(c.Request.Context(), code, redirectURI)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "oauth_error.html", gin.H{
			"error": "Authentication failed: " + err.Error(),
		})
		return
	}

	// For web flow: render success page with token or redirect
	c.HTML(http.StatusOK, "oauth_success.html", gin.H{
		"token": token,
		"user":  user,
	})
}

// Mobile-friendly flow that returns JSON
func (h *OAuthHandler) GoogleCallbackJSON(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		errorMsg := c.Query("error")
		if errorMsg == "" {
			errorMsg = "No authorization code provided"
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	// The redirect URI must match exactly what was used to get the code
	redirectURI := fmt.Sprintf("%s://%s%s", c.Request.URL.Scheme, c.Request.Host, c.Request.URL.Path)
	if c.Request.URL.Scheme == "" {
		proto := "http"
		if c.Request.TLS != nil {
			proto = "https"
		}
		redirectURI = fmt.Sprintf("%s://%s%s", proto, c.Request.Host, c.Request.URL.Path)
	}

	token, user, err := h.oauthUseCase.GoogleSignIn(c.Request.Context(), code, redirectURI)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"profile_picture": user.ProfilePicture,
		},
	})
}
