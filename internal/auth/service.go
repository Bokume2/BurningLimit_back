package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type Service struct {
	verifier TokenVerifier
}

func NewService(verifier TokenVerifier) (*Service, error) {
	if verifier == nil {
		return nil, errors.New("token verifier is required")
	}

	return &Service{verifier: verifier}, nil
}

func (s *Service) Authenticate(ctx context.Context, idToken string) (Identity, error) {
	if strings.TrimSpace(idToken) == "" {
		return Identity{}, ErrInvalidToken
	}

	identity, err := s.verifier.VerifyIDToken(ctx, idToken)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrVerificationUnavailable) {
			return Identity{}, err
		}

		return Identity{}, fmt.Errorf("%w: %v", ErrVerificationUnavailable, err)
	}

	if strings.TrimSpace(identity.UID) == "" {
		return Identity{}, ErrInvalidToken
	}

	return identity, nil
}
