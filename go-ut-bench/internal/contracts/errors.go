package contracts

type ErrorInfo struct {
	Kind       string `json:"kind"`
	Message    string `json:"message"`
	Retryable  bool   `json:"retryable"`
	StatusCode *int   `json:"status_code,omitempty"`
}
