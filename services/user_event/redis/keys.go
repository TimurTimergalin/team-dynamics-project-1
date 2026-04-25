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

func (k PlayerKeySet) keys() []string {
	return []string{k.name(), k.rating(), k.subscriptions(), k.currentChallenge(), k.status()}
}
