package models

type MessageType int64
const (
	Challenge MessageType = 1 + iota
	ChallengeAccepted
	ChallengeDeclined
	ChallengeCancelled
)

type Message struct {
	Type MessageType
	Id string
	SenderName string
	SenderId int64
	SenderRating int64
	ReceiverId int64
	Address string
	OriginalMessageId string
}

func NewChallenge(Id string, SenderName string, SenderId int64, SenderRating int64, ReceiverId int64) *Message {
	return &Message{Type: Challenge, Id: Id, SenderName: SenderName, SenderId: SenderId, SenderRating: SenderRating, ReceiverId: ReceiverId}
}

func NewChallengeAccepted(Id string, Address string) *Message {
	return &Message{Type: ChallengeAccepted, Id: Id, Address: Address}
}

func NewChallengeDeclined(Id string) *Message {
	return &Message{Type: ChallengeDeclined, Id: Id}
}

func NewChallengeCancelled(Id string, OriginalMessageId string) *Message {
	return &Message{Type: ChallengeCancelled, Id: Id, OriginalMessageId: OriginalMessageId}
}
