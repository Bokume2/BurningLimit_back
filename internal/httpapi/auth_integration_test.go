package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	domainAuth "github.com/Bokume2/FirstHackathon2026Summer_back/internal/auth"
	"github.com/Bokume2/FirstHackathon2026Summer_back/internal/firebaseauth"
	"github.com/labstack/echo/v5"
)

type firebaseAuthResponse struct {
	IDToken string `json:"idToken"`
	LocalID string `json:"localId"`
}

func TestAuthMeFirebaseIntegration(t *testing.T) {
	apiKey := os.Getenv("FIREBASE_WEB_API_KEY")
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if apiKey == "" || projectID == "" || os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("Firebase integration test environment is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	email := fmt.Sprintf("auth-integration-%d@example.com", time.Now().UnixNano())
	password := "integration-test-password-2026"
	signup := callFirebaseAuth(t, ctx, apiKey, "signUp", map[string]any{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	})
	t.Cleanup(func() {
		callFirebaseAuth(t, context.Background(), apiKey, "delete", map[string]any{
			"idToken": signup.IDToken,
		})
	})

	login := callFirebaseAuth(t, ctx, apiKey, "signInWithPassword", map[string]any{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	})

	firebaseClient, err := firebaseauth.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("firebaseauth.NewClient() error = %v", err)
	}
	verifier, err := firebaseauth.NewVerifier(firebaseClient)
	if err != nil {
		t.Fatalf("firebaseauth.NewVerifier() error = %v", err)
	}
	service, err := domainAuth.NewService(verifier)
	if err != nil {
		t.Fatalf("auth.NewService() error = %v", err)
	}

	e := echo.New()
	RegisterAuthRoutes(e, service)
	request := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	request.Header.Set(echo.HeaderAuthorization, "Bearer "+login.IDToken)
	response := httptest.NewRecorder()
	e.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("GET /auth/me status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode GET /auth/me response: %v", err)
	}
	if body["uid"] != signup.LocalID {
		t.Errorf("GET /auth/me uid = %q, want %q", body["uid"], signup.LocalID)
	}
}

func callFirebaseAuth(
	t *testing.T,
	ctx context.Context,
	apiKey string,
	action string,
	payload map[string]any,
) firebaseAuthResponse {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("encode Firebase %s request: %v", action, err)
	}
	endpoint := "https://identitytoolkit.googleapis.com/v1/accounts:" + action + "?key=" + url.QueryEscape(apiKey)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create Firebase %s request: %v", action, err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("Firebase %s request: %v", action, err)
	}
	defer response.Body.Close()

	var result firebaseAuthResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode Firebase %s response: %v", action, err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Firebase %s status = %d", action, response.StatusCode)
	}

	return result
}
