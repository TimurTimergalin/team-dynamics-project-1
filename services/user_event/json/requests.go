package json

import (
	"encoding/json"
	"fmt"
)

type RequestType string

const (
	AddFriend        RequestType = "AddFriend"
	RemoveFriend     RequestType = "RemoveFriend"
	Subscribe        RequestType = "Subscribe"
	Unsubscribe      RequestType = "Unsubscribe"
	Challenge        RequestType = "Challenge"
	CancelChallenge  RequestType = "CancelChallenge"
	AcceptChallenge  RequestType = "AcceptChallenge"
	DeclineChallenge RequestType = "DeclineChallenge"
)

type rawRequest struct {
	Type    RequestType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Request struct {
	Type    RequestType
	Payload any
}

type AddFriendPayload struct {
	OtherUserId int64 `json:"otherUserId"`
}

type RemoveFriendPayload struct {
	OtherUserId int64 `json:"otherUserId"`
}

type SubscribePayload struct {
	Users []int64 `json:"users"`
}

type UnsubscribePayload struct {
	Users []int64 `json:"users"`
}

type ChallengePayload struct {
	UserId int64 `json:"userId"`
}

type CancelChallengePayload struct {
}

type AcceptChallengePayload struct {
}

type DeclineChallengePayload struct {
}

func ParseRequest(data []byte) (*Request, error) {
	var raw rawRequest
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse request: %w", err)
	}

	var payload any
	switch raw.Type {
	case AddFriend:
		var p AddFriendPayload
		if err := json.Unmarshal(raw.Payload, &p); err != nil {
			return nil, fmt.Errorf("failed to parse AddFriend payload: %w", err)
		}
		payload = p
	case RemoveFriend:
		var p RemoveFriendPayload
		if err := json.Unmarshal(raw.Payload, &p); err != nil {
			return nil, fmt.Errorf("failed to parse RemoveFriend payload: %w", err)
		}
		payload = p
	case Subscribe:
		var p SubscribePayload
		if err := json.Unmarshal(raw.Payload, &p); err != nil {
			return nil, fmt.Errorf("failed to parse Subscribe payload: %w", err)
		}
		payload = p
	case Unsubscribe:
		var p UnsubscribePayload
		if err := json.Unmarshal(raw.Payload, &p); err != nil {
			return nil, fmt.Errorf("failed to parse Unsubscribe payload: %w", err)
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
	case AcceptChallenge:
		payload = AcceptChallengePayload{}
	case DeclineChallenge:
		payload = DeclineChallengePayload{}
	default:
		return nil, fmt.Errorf("unknown request type: %q", raw.Type)
	}

	return &Request{Type: raw.Type, Payload: payload}, nil
}