package mm_event

import "json"

type MatchmakingRequestBody struct{}

type MatchmakingRequest struct {
	Headers json.RequestHeaders         `json:"headers"`
	Body    MatchmakingRequestBody `json:"body"`
}

type ErrorResponseBody struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

type ErrorResponse struct {
	Headers json.ResponseHeaders   `json:"headers"`
	Body    ErrorResponseBody `json:"body"`
}
