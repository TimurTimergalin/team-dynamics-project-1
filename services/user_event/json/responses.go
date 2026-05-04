package json

import (
	"encoding/json"
	"fmt"
)

type ResponseType string

const (
	StatusChanged      ResponseType = "StatusChanged"
	ChallengeReceived  ResponseType = "ChallengeReceived"
	ChallengeAccepted  ResponseType = "ChallengeAccepted"
	ChallengeDeclined  ResponseType = "ChallengeDeclined"
	ChallengeCancelled ResponseType = "ChallengeCancelled"
	MatchStarted       ResponseType = "MatchStarted"
	Error              ResponseType = "Error"
)

type Event struct {
	Type    ResponseType    `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Response struct {
	Events []*Event `json:"events"`
}

type UserStatus string

const (
	Offline UserStatus = "Offline"
	Online  UserStatus = "Online"
	InGame  UserStatus = "InGame"
)

type StatusChangedPayload struct {
	UserId    int64      `json:"userId"`
	NewStatus UserStatus `json:"newStatus"`
}

type ChallengeReceivedPayload struct {
	MessageId string `json:"messageId"`
	UserId    int64  `json:"userId"`
	UserName  string `json:"userName"`
}

type ChallengeAcceptedPayload struct {
	Address string `json:"address"`
}

type ChallengeDeclinedPayload struct {
}

type ChallengeCancelledPayload struct {
}

type MatchStartedPayload struct {
	Address string `json:"address"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

type ResponsePayload interface {
	StatusChangedPayload | ChallengeReceivedPayload | ChallengeAcceptedPayload | ChallengeDeclinedPayload | ChallengeCancelledPayload | MatchStartedPayload | ErrorPayload
}

func MakeEvent[T ResponsePayload](payload T) (*Event, error) {
	var responseType ResponseType
	switch any(payload).(type) {
	case StatusChangedPayload:
		responseType = StatusChanged
	case ChallengeReceivedPayload:
		responseType = ChallengeReceived
	case ChallengeAcceptedPayload:
		responseType = ChallengeAccepted
	case ChallengeDeclinedPayload:
		responseType = ChallengeDeclined
	case ChallengeCancelledPayload:
		responseType = ChallengeCancelled
	case MatchStartedPayload:
		responseType = MatchStarted
	case ErrorPayload:
		responseType = Error
	}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize event payload: %w", err)
	}
	return &Event{Type: responseType, Payload: rawPayload}, nil
}

func SerializeResponse(response *Response) ([]byte, error) {
	return json.Marshal(response)
}
