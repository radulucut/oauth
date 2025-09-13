package oauth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type FacebookPayload struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"` // can be empty
}

type FacebookError struct {
	Err struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	} `json:"error"`
}

func (e *FacebookError) Error() string {
	return e.Err.Message
}

func (c *client) Facebook(ctx context.Context, token string) (*FacebookPayload, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.facebookURL+"?fields=email,name&access_token="+token, nil)
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
