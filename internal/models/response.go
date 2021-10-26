package models

// BadResponse type which presents response for all bad responses.
type BadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
