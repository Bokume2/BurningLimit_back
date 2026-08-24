package auth

import (
	"context"
	"errors"
)

var (
	ErrInvalidToken            = errors.New("invalid authentication token")
	ErrVerificationUnavailable = errors.New("authentication verification unavailable")
)

type TokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (Identity, error)
}
