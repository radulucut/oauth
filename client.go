package oauth

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client struct {
	googleURL    string
	facebookURL  string
	microsoftURL string
	httpClient   *http.Client
}

type Config struct {
	Timeout time.Duration
}

func NewClient(config Config) *Client {
	return &Client{
		googleURL:    "https://www.googleapis.com/oauth2/v3/userinfo",
		facebookURL:  "https://graph.facebook.com/me",
		microsoftURL: "https://graph.microsoft.com/v1.0/me",
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (c *Client) Google(token *string) (*GooglePayload, error) {
	res, err := c.httpClient.Get(c.googleURL + "?access_token=" + *token)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(io.LimitReader(res.Body, 1<<20)) // 1MB
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		rErr := &GoogleError{}
		err = json.Unmarshal(b, rErr)
		if err != nil {
			return nil, err
		}
		return nil, rErr
	}
	payload := &GooglePayload{}
	err = json.Unmarshal(b, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (c *Client) Facebook(token *string) (*FacebookPayload, error) {
	res, err := c.httpClient.Get(c.facebookURL + "?fields=email,name&access_token=" + *token)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(io.LimitReader(res.Body, 1<<20)) // 1MB
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		rErr := &FacebookError{}
		err = json.Unmarshal(b, rErr)
		if err != nil {
			return nil, err
		}
		return nil, rErr
	}
	payload := &FacebookPayload{}
	err = json.Unmarshal(b, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (c *Client) Microsoft(token *string) (*MicrosoftPayload, error) {
	req, err := http.NewRequest(http.MethodGet, c.microsoftURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+*token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(io.LimitReader(res.Body, 1<<20)) // 1MB
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		rErr := &MicrosoftError{}
		err = json.Unmarshal(b, rErr)
		if err != nil {
			return nil, err
		}
		return nil, rErr
	}
	payload := &MicrosoftPayload{}
	err = json.Unmarshal(b, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
