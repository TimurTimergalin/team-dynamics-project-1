package models

import "time"

type Round struct {
	IsPlayer1Killer bool
	Length          time.Duration
}

type MatchResult int32

const (
	Win  MatchResult = iota
	Draw
	Loss
)

type Match struct {
	MatchId       string
	Player1Id     int64
	Player2Id     int64
	Player1Name   string
	Player2Name   string
	Player1Rating int64
	Player2Rating int64
	End           time.Time
	Result        MatchResult
}

type AggregatedMatch struct {
	MatchObj *Match
	Rounds   []*Round
}

type PageKey struct {
	Before time.Time `json:"before"`
}

type GameResult = int32

const (
	Winner GameResult = 2
	DrawResult GameResult = 1
	Loser  GameResult = 0
)

type RatingInfo struct {
	UserId     int64
	Value      float64
	Deviation  float64
	Volatility float64
	LastUpdate time.Time
}
