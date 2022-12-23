package model

type Connection struct {
	ConnectionID   string
	ExpirationTime int
}

type Request[T any] struct {
	Action  string `json:"action"`
	Payload T      `json:"payload"`
}

type Response[T any] struct {
	Action   string `json:"action"`
	Response T      `json:"response"`
}

type PingRequestPayload struct{}
type PongResponsePayload struct{}

type MessageRequestPayload struct {
	Message string `json:"message"`
}

type MessageResponsePayload struct {
	Message string `json:"message"`
}
