package oauth

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
