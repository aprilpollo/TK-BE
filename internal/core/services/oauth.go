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
	orgRepo        output.OrganizationRepository
	googleVerifier output.GoogleVerifier
	jwtCfg         utils.JWTConfig
}

func NewOauthService(repo output.OauthRepository, orgRepo output.OrganizationRepository, googleVerifier output.GoogleVerifier, jwtCfg utils.JWTConfig) input.OauthService {
	return &oauthService{repo: repo, orgRepo: orgRepo, googleVerifier: googleVerifier, jwtCfg: jwtCfg}
}

func (s *oauthService) buildLoginResult(ctx context.Context, userID int64, email string) (*domain.LoginResult, error) {
	token, err := utils.GenerateToken(userID, email, s.jwtCfg)
	if err != nil {
		return nil, err
	}

	orgs, err := s.orgRepo.FindByUserIDWithPrimaryDetails(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResult{Token: token, Organizations: orgs}, nil
}

func (s *oauthService) BasicLogin(ctx context.Context, oauth *domain.BasicLogin) (*domain.LoginResult, error) {
	user, err := s.repo.BasicLogin(ctx, oauth)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Password == nil || !utils.ComparePassword(oauth.Password, *user.Password) {
		return nil, errors.New("invalid email or password")
	}

	return s.buildLoginResult(ctx, user.ID, user.Email)
}

func (s *oauthService) SocialLogin(ctx context.Context, req *domain.SocialLogin) (*domain.LoginResult, error) {
	switch req.Provider {
	case domain.OauthProviderGoogle:
		info, err := s.googleVerifier.VerifyIDToken(ctx, req.Token)
		if err != nil {
			return nil, errors.New("invalid google token")
		}

		user, err := s.repo.UpsertSocialUser(ctx, info, domain.OauthProviderGoogle)
		if err != nil {
			return nil, err
		}

		return s.buildLoginResult(ctx, user.ID, user.Email)

	default:
		return nil, errors.New("unsupported provider")
	}
}

func (s *oauthService) Register(ctx context.Context) {}

func (s *oauthService) RefreshToken(ctx context.Context) {}
