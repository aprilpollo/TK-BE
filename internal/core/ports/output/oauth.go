package output

import (
	"context"

	"aprilpollo/internal/core/domain"

	// "aprilpollo/internal/pkg/query"
)

// OauthRepository is the output port — defines how the core communicates with storage.
type OauthRepository interface {
	BasicLogin(ctx context.Context, oauth *domain.BasicLogin) (*domain.OauthUser, error)
	UpsertSocialUser(ctx context.Context, info *domain.GoogleUserInfo, provider domain.OauthProvider) (*domain.OauthUser, error)
	Register(ctx context.Context)
	RefreshToken(ctx context.Context)
}