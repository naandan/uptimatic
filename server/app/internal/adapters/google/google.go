package google

import (
	"uptimatic/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleCLient struct {
	*oauth2.Config
}

func NewGoogleClient(cfg *config.Config) *GoogleCLient {
	return &GoogleCLient{
		Config: &oauth2.Config{
			RedirectURL:  cfg.AppScheme + "://" + cfg.AppDomain + "/api/v1/auth/google/callback",
			ClientID:     cfg.GoogleClientId,
			ClientSecret: cfg.GoogleClientSecret,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
	}
}
