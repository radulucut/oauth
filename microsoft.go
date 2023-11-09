package oauth

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
