package oauth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type MicrosoftPayload struct {
	Id                string `json:"id"`
	Email             string `json:"mail"`
	DisplayName       string `json:"displayName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	PreferredLanguage string `json:"preferredLanguage"`
}

type MicrosoftError struct {
	Err struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

func (e *MicrosoftError) Error() string {
	return "Microsoft error"
}

func (c *client) Microsoft(ctx context.Context, token string) (*MicrosoftPayload, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.microsoftURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
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
