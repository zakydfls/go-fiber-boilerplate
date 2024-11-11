package responses

type APIResponse struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
