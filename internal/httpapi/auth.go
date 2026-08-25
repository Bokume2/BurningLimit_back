package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Bokume2/FirstHackathon2026Summer_back/internal/auth"
	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	service *auth.Service
}

func RegisterAuthRoutes(e *echo.Echo, service *auth.Service) {
	handler := AuthHandler{service: service}
	e.GET("/auth/me", handler.Me)
}

func (h AuthHandler) Me(c *echo.Context) error {
	idToken, ok := bearerToken(c.Request().Header.Get(echo.HeaderAuthorization))
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authentication required"})
	}

	identity, err := h.service.Authenticate(c.Request().Context(), idToken)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidToken):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authentication token"})
		case errors.Is(err, auth.ErrVerificationUnavailable):
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "authentication unavailable"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "authentication failed"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"uid": identity.UID})
}

func bearerToken(authorization string) (string, bool) {
	scheme, token, found := strings.Cut(strings.TrimSpace(authorization), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
		return "", false
	}

	return strings.TrimSpace(token), true
}
