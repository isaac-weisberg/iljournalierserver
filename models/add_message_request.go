package models

type AddMessageRequest struct {
	UnixSeconds int64
	Message     string
}

func NewAddMessageRequest(unixSeconds int64, message string) AddMessageRequest {
	return AddMessageRequest{UnixSeconds: unixSeconds, Message: message}
}
