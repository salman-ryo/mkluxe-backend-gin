package handler

import (
	"errors"
	"net/http"

	"mkluxe-backend/internal/config"
	"mkluxe-backend/internal/dto"
	"mkluxe-backend/internal/response"
	"mkluxe-backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(svc *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: svc,
		cfg:         cfg,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	authRes, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	// Fetch unified configuration definitions
	accessCfg := h.cfg.GetAccessCookieConfig()
	refreshCfg := h.cfg.GetRefreshCookieConfig()

	// Write token values securely down to HttpOnly storage
	setSecureCookie(c, accessCfg.Name, authRes.AccessToken, accessCfg.MaxAge, accessCfg.Path, accessCfg.HttpOnly)
	setSecureCookie(c, refreshCfg.Name, authRes.RefreshToken, refreshCfg.MaxAge, refreshCfg.Path, refreshCfg.HttpOnly)

	// Return response (tokens are safely ignored in JSON due to json:"-" tags)
	response.OK(c, "Login successful", authRes)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	// Read directly out of secure incoming cookies instead of public payload inputs
	refreshCfg := h.cfg.GetRefreshCookieConfig()
	tokenStr, err := c.Cookie(refreshCfg.Name)
	if err != nil {
		response.Unauthorized(c, "Refresh token missing")
		return
	}

	authRes, err := h.authService.RefreshToken(c.Request.Context(), tokenStr)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	accessCfg := h.cfg.GetAccessCookieConfig()
	newRefreshCfg := h.cfg.GetRefreshCookieConfig()

	// Re-apply refreshed keys directly to browser context
	setSecureCookie(c, accessCfg.Name, authRes.AccessToken, accessCfg.MaxAge, accessCfg.Path, accessCfg.HttpOnly)
	setSecureCookie(c, newRefreshCfg.Name, authRes.RefreshToken, newRefreshCfg.MaxAge, newRefreshCfg.Path, newRefreshCfg.HttpOnly)

	response.OK(c, "Token refreshed successfully", nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	accessCfg := h.cfg.GetAccessCookieConfig()
	refreshCfg := h.cfg.GetRefreshCookieConfig()

	// Expire existing tracking keys inside client browsers completely by resetting MaxAge to -1
	setSecureCookie(c, accessCfg.Name, "", -1, accessCfg.Path, accessCfg.HttpOnly)
	setSecureCookie(c, refreshCfg.Name, "", -1, refreshCfg.Path, refreshCfg.HttpOnly)

	response.OK(c, "Logged out successfully", nil)
}

func (h *AuthHandler) CurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized request")
		return
	}

	userRes, err := h.authService.GetCurrentUser(c.Request.Context(), userID.(string))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			response.NotFound(c, "User profile not found")
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, "Current user fetched successfully", userRes)
}

func setSecureCookie(c *gin.Context, name, value string, maxAge int, path string, httpOnly bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   "", // 💡 Stop setting the domain to allow browser/cross-site default binding
		Secure:   true,
		HttpOnly: httpOnly,
		SameSite: http.SameSiteNoneMode,
	})
}
