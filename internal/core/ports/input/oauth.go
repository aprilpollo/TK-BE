package input

import (
	"context"

	"aprilpollo/internal/core/domain"
	// "aprilpollo/internal/pkg/query"
)

// OauthService is the input port — defines what the application can do with oauth.
type OauthService interface {
	BasicLogin(ctx context.Context, oauth *domain.BasicLogin) (string, error)
	SocialLogin(ctx context.Context, req *domain.SocialLogin) (string, error)
	Register(ctx context.Context)
	RefreshToken(ctx context.Context)
}
