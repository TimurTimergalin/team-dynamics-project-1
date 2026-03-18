package online_event

type ChallengeRequestBody struct {
	RecipientId uint64 `json:"recipient_id"`
}

type ChallengeRequest struct {
	Headers json.RequestHeaders  `json:"headers"`
	Body    ChallengeRequestBody `json:"body"`
}

type IsOnlineBody struct {
	UserId uint64 `json:"user_id"`
}

type IsOnline struct {
	Headers json.ResponseHeaders `json:"headers"`
	Body    IsOnlineBody         `json:"body"`
}
