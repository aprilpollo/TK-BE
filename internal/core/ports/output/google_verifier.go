package output

import (
	"context"

	"aprilpollo/internal/core/domain"
)

// GoogleVerifier is the output port for verifying Google ID tokens.
type GoogleVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*domain.GoogleUserInfo, error)
}
