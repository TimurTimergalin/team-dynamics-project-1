package json

type RequestHeaders struct {
	AccessToken string `json:"access_token"`
}

type ResponseHeaders struct{}

type MatchResponseBody struct {
	Address string `json:"address"`
}

type MatchResponse struct {
	Headers ResponseHeaders  `json:"headers"`
	Body    MatchResponseBody `json:"body"`
}