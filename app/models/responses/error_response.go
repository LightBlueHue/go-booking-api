package responses

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Target  string        `json:"target"`
	Details []ErrorDetail `json:"details"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target"`
}
