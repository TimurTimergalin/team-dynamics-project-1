package redis

import "fmt"

type PlayerKeySet struct {
	id int64
}

func (k PlayerKeySet) name() string {
	return fmt.Sprintf("name:%d", k.id)
}

func (k PlayerKeySet) rating() string {
	return fmt.Sprintf("rating:%d", k.id)
}

func (k PlayerKeySet) subscriptions() string {
	return fmt.Sprintf("subscriptions:%d", k.id)
}

func (k PlayerKeySet) currentChallenge() string {
	return fmt.Sprintf("currentChallenge:%d", k.id)
}

func (k PlayerKeySet) status() string {
	return fmt.Sprintf("status:%d", k.id)
}
