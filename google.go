package oauth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type GooglePayload struct {
	Id            string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	PictureURL    string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

type GoogleError struct {
	Name        string `json:"error"`
	Description string `json:"error_description"`
}

func (e *GoogleError) Error() string {
	return e.Description
}

func (c *client) Google(ctx context.Context, token string) (*GooglePayload, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.googleURL+"?access_token="+token, nil)
	if err != nil {
		return nil, err
	}
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
