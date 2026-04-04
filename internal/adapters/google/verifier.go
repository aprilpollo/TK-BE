package google

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/output"

	"google.golang.org/api/idtoken"
)

type googleVerifier struct {
	clientID string
}

func NewGoogleVerifier(clientID string) output.GoogleVerifier {
	return &googleVerifier{clientID: clientID}
}

func (g *googleVerifier) VerifyIDToken(ctx context.Context, idToken string) (*domain.GoogleUserInfo, error) {
	payload, err := idtoken.Validate(ctx, idToken, g.clientID)
	if err != nil {
		return nil, err
	}

	firstName, _ := payload.Claims["given_name"].(string)
	lastName, _ := payload.Claims["family_name"].(string)
	displayName, _ := payload.Claims["name"].(string)
	email, _ := payload.Claims["email"].(string)
	picture, _ := payload.Claims["picture"].(string)

	var avatar *string
	if picture != "" {
		avatar = &picture
	}

	return &domain.GoogleUserInfo{
		ProviderID:  payload.Subject,
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		DisplayName: displayName,
		Avatar:      avatar,
	}, nil
}
