package firebaseauth

import (
	"context"
	"fmt"
	"strings"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
)

func NewClient(ctx context.Context, projectID string) (*firebaseAuth.Client, error) {
	if strings.TrimSpace(projectID) == "" {
		return nil, fmt.Errorf("firebase project ID is required")
	}
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: projectID,
	})
	if err != nil {
		return nil, fmt.Errorf("initialize Firebase app: %w", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialize Firebase Auth client: %w", err)
	}
	return client, nil
}
