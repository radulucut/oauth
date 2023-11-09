package oauth

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
