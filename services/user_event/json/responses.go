package json

import (
	"encoding/json"
	"fmt"
)

type ResponseType string

const (
	StatusChanged     ResponseType = "StatusChanged"
	ChallengeReceived ResponseType = "ChallengeReceived"
	ChallengeAccepted ResponseType = "ChallengeAccepted"
	ChallengeDeclined ResponseType = "ChallengeDeclined"
)

type Response struct {
	Type    ResponseType    `json:"type"`
	Payload json.RawMessage `json:"payload"`
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
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
}

type ChallengeAcceptedPayload struct {
}

type ChallengeDeclinedPayload struct {
}

type ResponsePayload interface {
	StatusChangedPayload | ChallengeReceivedPayload | ChallengeAcceptedPayload | ChallengeDeclinedPayload
}

func SerializeResponse[T ResponsePayload](payload T) ([]byte, error) {
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
	}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize response payload: %w", err)
	}
	return json.Marshal(Response{
		Type:    responseType,
		Payload: rawPayload,
	})
}
