package models

type ConnectionStatusResponse struct {
	Connected bool        `json:"connected"`
	Error     string      `json:"error,omitempty"`
	Config    interface{} `json:"config"`
}
