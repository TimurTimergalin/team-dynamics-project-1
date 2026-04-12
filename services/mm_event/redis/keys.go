package redis

import "fmt"

type playerKeys struct {
	id int64
}

func (ks playerKeys) rating() string {
	return fmt.Sprintf("rating:%d", ks.id)
}

func (ks playerKeys) fleet() string {
	return fmt.Sprintf("fleet:%d", ks.id)
}

func (ks playerKeys) name() string {
	return fmt.Sprintf("name:%d", ks.id)
}

func (ks playerKeys) displayedRating() string {
	return fmt.Sprintf("displayed_rating:%d", ks.id)
}

func (ks playerKeys) regId() string {
	return fmt.Sprintf("reg_id:%d", ks.id)
}

func (ks playerKeys) keys() []string {
	return []string{ks.rating(), ks.fleet(), ks.name(), ks.displayedRating(), ks.regId()}
}
