package models

import "time"

type GameResult = int32

const (
	Winner GameResult = 2
	Draw              = 1
	Loser             = 0
)

type RatingInfo struct {
	UserId     int64
	Value      float64
	Deviation  float64
	Volatility float64
	LastUpdate time.Time
}

type MatchInfo struct {
	MatchId string
}
