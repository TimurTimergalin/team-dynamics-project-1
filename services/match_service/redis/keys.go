package redis

import "fmt"

type playerKeys struct {
	id int64
}

func (ks playerKeys) matchId() string {
	return fmt.Sprintf("match_id:%d", ks.id)
}

func (ks playerKeys) regId() string {
	return fmt.Sprintf("reg_id:%d", ks.id)
}

func (ks playerKeys) rating() string {
	return fmt.Sprintf("rating:%d", ks.id)
}

func (ks playerKeys) name() string {
	return fmt.Sprintf("match_id:%d", ks.id)
}

func (ks playerKeys) keys() []string {
	return []string{ks.matchId(), ks.regId(), ks.rating(), ks.name()}
}

type matchKeys interface {
	status() string
	fleet() string
	player1() string
	player2() string
}

type realMatchKeys struct {
	id string
}

func (ks realMatchKeys) status() string {
	return fmt.Sprintf("status:%s", ks.id)
}

func (ks realMatchKeys) fleet() string {
	return fmt.Sprintf("fleet:%s", ks.id)
}

func (ks realMatchKeys) player1() string {
	return fmt.Sprintf("player1:%s", ks.id)
}

func (ks realMatchKeys) player2() string {
	return fmt.Sprintf("player2:%s", ks.id)
}

func (ks realMatchKeys) keys() []string {
	return []string{ks.status(), ks.fleet(), ks.player1(), ks.player2()}
}

type absentMatchKeys struct{}

func (ks absentMatchKeys) status() string {
	return ""
}

func (ks absentMatchKeys) fleet() string {
	return ""
}

func (ks absentMatchKeys) player1() string {
	return ""
}

func (ks absentMatchKeys) player2() string {
	return ""
}
