package services

import (
	"context"
	"errors"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/utils"
)

type oauthService struct {
	repo           output.OauthRepository
	googleVerifier output.GoogleVerifier
	jwtCfg         utils.JWTConfig
}

func NewOauthService(repo output.OauthRepository, googleVerifier output.GoogleVerifier, jwtCfg utils.JWTConfig) input.OauthService {
	return &oauthService{repo: repo, googleVerifier: googleVerifier, jwtCfg: jwtCfg}
}

func (s *oauthService) BasicLogin(ctx context.Context, oauth *domain.BasicLogin) (string, error) {
	user, err := s.repo.BasicLogin(ctx, oauth)
	if err != nil || user == nil {
		return "", errors.New("invalid email or password")
	}

	if user.Password == nil || !utils.ComparePassword(oauth.Password, *user.Password) {
		return "", errors.New("invalid email or password")
	}

	return utils.GenerateToken(user.ID, user.Email, s.jwtCfg)
}

func (s *oauthService) SocialLogin(ctx context.Context, req *domain.SocialLogin) (string, error) {
	switch req.Provider {
	case domain.OauthProviderGoogle:
		info, err := s.googleVerifier.VerifyIDToken(ctx, req.Token)
		if err != nil {
			return "", errors.New("invalid google token")
		}

		user, err := s.repo.UpsertSocialUser(ctx, info, domain.OauthProviderGoogle)
		if err != nil {
			return "", err
		}

		return utils.GenerateToken(user.ID, user.Email, s.jwtCfg)

	default:
		return "", errors.New("unsupported provider")
	}
}

func (s *oauthService) Register(ctx context.Context) {}

func (s *oauthService) RefreshToken(ctx context.Context) {}
