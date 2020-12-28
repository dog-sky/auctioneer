package v1

type ResponseV1 struct {
	Success bool   `json:"success,omitempty"`
	Result  string `json:"result,omitempty"`
}
