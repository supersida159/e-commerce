package services_auth

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
)

type GoogleAuthService struct {
	oauth2Service *oauth2.Service
}

func NewGoogleAuthService() (*GoogleAuthService, error) {
	ctx := context.Background()
	oauth2Service, err := oauth2.NewService(ctx, option.WithHTTPClient(http.DefaultClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 service: %v", err)
	}

	return &GoogleAuthService{
		oauth2Service: oauth2Service,
	}, nil
}

func (s *GoogleAuthService) VerifyGoogleToken(accessToken string) (*oauth2.Tokeninfo, error) {
	tokenInfo, err := s.oauth2Service.Tokeninfo().AccessToken(accessToken).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}
	return tokenInfo, nil
}
