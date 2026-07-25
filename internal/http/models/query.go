package models

type QueryPayload struct {
	Query string `json:"query"`
	Force bool   `json:"force"`
}

type DestructiveWarningResponse struct {
	Warning string `json:"warning"`
	Code    string `json:"code"`
}
