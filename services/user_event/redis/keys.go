package redis

import "fmt"

type PlayerKeySet struct {
	Id int64
}

func (k PlayerKeySet) name() string {
	return fmt.Sprintf("name:%d", k.Id)
}

func (k PlayerKeySet) rating() string {
	return fmt.Sprintf("rating:%d", k.Id)
}

func (k PlayerKeySet) subscriptions() string {
	return fmt.Sprintf("subscriptions:%d", k.Id)
}

func (k PlayerKeySet) currentChallenge() string {
	return fmt.Sprintf("currentChallenge:%d", k.Id)
}

func (k PlayerKeySet) status() string {
	return fmt.Sprintf("status:%d", k.Id)
}

func (k PlayerKeySet) mailbox() string {
	return fmt.Sprintf("mailbox:%d", k.Id)
}

func (k PlayerKeySet) keys() []string {
	return []string{k.name(), k.rating(), k.subscriptions(), k.currentChallenge(), k.status(), k.mailbox()}
}

type messageKeySet struct {
	id string
}

func (k messageKeySet) type_() string {
	return fmt.Sprintf("type:%s", k.id)
}

func (k messageKeySet) senderName() string {
	return fmt.Sprintf("senderName:%s", k.id)
}

func (k messageKeySet) senderId() string {
	return fmt.Sprintf("senderId:%s", k.id)
}

func (k messageKeySet) senderRating() string {
	return fmt.Sprintf("senderRating:%s", k.id)
}

func (k messageKeySet) receiverId() string {
	return fmt.Sprintf("receiverId:%s", k.id)
}

func (k messageKeySet) address() string {
	return fmt.Sprintf("address:%s", k.id)
}

func (k messageKeySet) originalMessageId() string {
	return fmt.Sprintf("originalMessageId:%s", k.id)
}
