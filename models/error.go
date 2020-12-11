package models

type Error struct {
	Error struct {
		Code       string `json:"code,readonly"`
		Message    string `json:"message,readonly"`
		InnerError struct {
			Date            string `json:"date,readonly"`
			RequestId       string `json:"request-id,readonly"`
			ClientRequestId string `json:"client-request-id,readonly"`
		} `json:"innerError,readonly"`
	} `json:"error,readonly"`
}
