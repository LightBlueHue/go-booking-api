package responses

// ErrorResponse contains response information about an error that occured for a request.
type ErrorResponse struct {
	Error Error `json:"error"`
}

// Error contains information such as http status code and error details.
type Error struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Target  string        `json:"target"`
	Details []ErrorDetail `json:"details"`
}

// ErrorDetail conatins detailed error information.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target"`
}
