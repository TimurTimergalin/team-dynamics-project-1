package models

type Player struct {
	Id     int64
	Name   string
	Rating int64
	RegId  string
}

type MatchStatus int32

const (
	Requested MatchStatus = 1 + iota
	Ongoing
	Finished
	Initialising
)

type PlayerFailResponse int32

const (
	REENTER PlayerFailResponse = 1 + iota
	REMOVE
)

type Match struct {
	MatchId   string
	Player1Id int64
	Player2Id int64
	Status    MatchStatus
	Fleet     string
}
