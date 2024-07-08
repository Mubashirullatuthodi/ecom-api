package models

type ErrorHandler struct {
	Error      string `json:"error"`
	StatusCode int    `json:"code"`
	Status     string `json:"status"`
}
