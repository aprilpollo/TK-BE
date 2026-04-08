package input

import (
	"context"

	"aprilpollo/internal/core/domain"
	// "aprilpollo/internal/pkg/query"
)

// OauthService is the input port — defines what the application can do with oauth.
type OauthService interface {
	BasicLogin(ctx context.Context, oauth *domain.BasicLogin) (*domain.LoginResult, error)
	SocialLogin(ctx context.Context, req *domain.SocialLogin) (*domain.LoginResult, error)
	Register(ctx context.Context)
	RefreshToken(ctx context.Context)
}
