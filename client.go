package oauth

import (
	"context"
	"net/http"
	"time"
)

type Client interface {
	Google(ctx context.Context, token string) (*GooglePayload, error)
	Facebook(ctx context.Context, token string) (*FacebookPayload, error)
	Microsoft(ctx context.Context, token string) (*MicrosoftPayload, error)
}

type client struct {
	googleURL    string
	facebookURL  string
	microsoftURL string
	httpClient   *http.Client
}

type Config struct {
	Timeout time.Duration
}

func NewClient(config Config) Client {
	return &client{
		googleURL:    "https://www.googleapis.com/oauth2/v3/userinfo",
		facebookURL:  "https://graph.facebook.com/me",
		microsoftURL: "https://graph.microsoft.com/v1.0/me",
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}
