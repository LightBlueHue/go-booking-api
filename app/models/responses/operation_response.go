package responses

type OperationResponse struct {
	Context string      `json:"@odata.context"`
	Value   interface{} `json:"value"`
}
