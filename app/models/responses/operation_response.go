package responses

// OperationResponse contains response information if no error.
type OperationResponse struct {
	Context string      `json:"@odata.context"`
	Value   interface{} `json:"value"`
}
