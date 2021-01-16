package models

// Error is used to unmarshal an error response from Microsoft Graph.
type Error struct {
	Error struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		InnerError struct {
			Date            string `json:"date"`
			RequestId       string `json:"request-id"`
			ClientRequestId string `json:"client-request-id"`
		} `json:"innerError"`
	} `json:"error"`
}
