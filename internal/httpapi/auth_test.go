package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bokume2/FirstHackathon2026Summer_back/internal/auth"
	"github.com/labstack/echo/v5"
)

type tokenVerifierStub struct {
	identity auth.Identity
	err      error
	token    string
}

func (v *tokenVerifierStub) VerifyIDToken(_ context.Context, token string) (auth.Identity, error) {
	v.token = token
	return v.identity, v.err
}

func TestBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
		ok     bool
	}{
		{name: "valid", header: "Bearer id-token", want: "id-token", ok: true},
		{name: "case insensitive", header: "bearer id-token", want: "id-token", ok: true},
		{name: "missing token", header: "Bearer", ok: false},
		{name: "wrong scheme", header: "Basic credentials", ok: false},
		{name: "empty", header: "", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := bearerToken(tt.header)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("bearerToken(%q) = (%q, %t), want (%q, %t)", tt.header, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestAuthMe(t *testing.T) {
	tests := []struct {
		name       string
		header     string
		identity   auth.Identity
		verifyErr  error
		wantStatus int
		wantToken  string
	}{
		{
			name:       "authenticated",
			header:     "Bearer firebase-id-token",
			identity:   auth.Identity{UID: "firebase-user"},
			wantStatus: http.StatusOK,
			wantToken:  "firebase-id-token",
		},
		{
			name:       "missing authorization",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			header:     "Bearer invalid-token",
			verifyErr:  auth.ErrInvalidToken,
			wantStatus: http.StatusUnauthorized,
			wantToken:  "invalid-token",
		},
		{
			name:       "verification unavailable",
			header:     "Bearer firebase-id-token",
			verifyErr:  auth.ErrVerificationUnavailable,
			wantStatus: http.StatusServiceUnavailable,
			wantToken:  "firebase-id-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier := &tokenVerifierStub{identity: tt.identity, err: tt.verifyErr}
			service, err := auth.NewService(verifier)
			if err != nil {
				t.Fatalf("auth.NewService() error = %v", err)
			}

			e := echo.New()
			RegisterAuthRoutes(e, service)
			request := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
			request.Header.Set(echo.HeaderAuthorization, tt.header)
			response := httptest.NewRecorder()

			e.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("GET /auth/me status = %d, want %d; body = %s", response.Code, tt.wantStatus, response.Body.String())
			}
			if verifier.token != tt.wantToken {
				t.Errorf("verified token = %q, want %q", verifier.token, tt.wantToken)
			}
		})
	}
}
