package firebaseauth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	firebaseAuth "firebase.google.com/go/v4/auth"
	domainAuth "github.com/Bokume2/FirstHackathon2026Summer_back/internal/auth"
)

type Verifier struct {
	client *firebaseAuth.Client
}

func NewVerifier(client *firebaseAuth.Client) (*Verifier, error) {
	if client == nil {
		return nil, errors.New("firebase auth client is required")
	}

	return &Verifier{client: client}, nil
}

func (v *Verifier) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (domainAuth.Identity, error) {
	if strings.TrimSpace(idToken) == "" {
		return domainAuth.Identity{}, domainAuth.ErrInvalidToken
	}

	token, err := v.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		if firebaseAuth.IsIDTokenInvalid(err) {
			return domainAuth.Identity{}, fmt.Errorf(
				"%w: %v",
				domainAuth.ErrInvalidToken,
				err,
			)
		}

		return domainAuth.Identity{}, fmt.Errorf(
			"%w: %v",
			domainAuth.ErrVerificationUnavailable,
			err,
		)
	}

	if token.UID == "" {
		return domainAuth.Identity{}, domainAuth.ErrInvalidToken
	}

	return domainAuth.Identity{UID: token.UID}, nil
}
