package json

import (
	"encoding/json"
	"fmt"
)

type RequestType string

const (
	Subscribe        RequestType = "Subscribe"
	Challenge        RequestType = "Challenge"
	CancelChallenge  RequestType = "CancelChallenge"
	AcceptChallenge  RequestType = "AcceptChallenge"
	DeclineChallenge RequestType = "DeclineChallenge"
	NotifyBusy       RequestType = "NotifyBusy"
	NotifyFree       RequestType = "NotifyFree"
)

type rawRequest struct {
	Type    RequestType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Request struct {
	Type    RequestType
	Payload any
}

type SubscribePayload struct {
	Users []int64 `json:"users"`
}

type ChallengePayload struct {
	UserId int64 `json:"userId"`
}

type CancelChallengePayload struct {
}

type AcceptChallengePayload struct {
	MessageId string `json:"messageId"`
	UserId    int64  `json:"userId"`
}

type DeclineChallengePayload struct {
	MessageId string `json:"messageId"`
	UserId    int64  `json:"userId"`
}

type NotifyBusyPayload struct {
}

type NotifyFreePayload struct {
}

func ParseRequest(data []byte) (*Request, error) {
	var raw rawRequest
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse request: %w", err)
	}

	var payload any
	switch raw.Type {
	case Subscribe:
		var p SubscribePayload
		if err := json.Unmarshal(raw.Payload, &p); err != nil {
			return nil, fmt.Errorf("failed to parse Subscribe payload: %w", err)
		}
		payload = p
	case Challenge:
		var p ChallengePayload
		if err := json.Unmarshal(raw.Payload, &p); err != nil {
			return nil, fmt.Errorf("failed to parse Challenge payload: %w", err)
		}
		payload = p
	case CancelChallenge:
		payload = CancelChallengePayload{}
	case NotifyBusy:
		payload = NotifyBusyPayload{}
	case NotifyFree:
		payload = NotifyFreePayload{}
	case AcceptChallenge:
		var p AcceptChallengePayload
		if err := json.Unmarshal(raw.Payload, &p); err != nil {
			return nil, fmt.Errorf("failed to parse Challenge payload: %w", err)
		}
		payload = p
	case DeclineChallenge:
		var p DeclineChallengePayload
		if err := json.Unmarshal(raw.Payload, &p); err != nil {
			return nil, fmt.Errorf("failed to parse DeclineChallenge payload: %w", err)
		}
		payload = p
	default:
		return nil, fmt.Errorf("unknown request type: %q", raw.Type)
	}

	return &Request{Type: raw.Type, Payload: payload}, nil
}
