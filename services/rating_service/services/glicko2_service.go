package services

import (
	"time"

	glicko2 "github.com/zelenin/go-glicko2"
	"team_dynamics/rating_service/models"
)

type Glicko2Service interface {
	GetInitialRating(userId int64) *models.RatingInfo
	UpdateRatings(a, b *models.RatingInfo, result models.GameResult) (*models.RatingInfo, *models.RatingInfo)
}

type glicko2ServiceImpl struct {
	Tau     float64
	Epsilon float64
}

func MakeGlicko2Service() Glicko2Service {
	return &glicko2ServiceImpl{
		Tau:     0.5,
		Epsilon: 1e-6,
	}
}

func (s *glicko2ServiceImpl) GetInitialRating(userId int64) *models.RatingInfo {
	return &models.RatingInfo{
		UserId:     userId,
		Value:      1200.0,
		Deviation:  300.0,
		Volatility: 0.08,
		LastUpdate: time.Now().UTC(),
	}
}

func convertModelMatchResult(mod models.GameResult) glicko2.MatchResult {
	switch mod {
	case models.Winner:
		return glicko2.MATCH_RESULT_WIN
	case models.Draw:
		return glicko2.MATCH_RESULT_DRAW
	default:
		return glicko2.MATCH_RESULT_LOSS
	}
}

func (s *glicko2ServiceImpl) UpdateRatings(a, b *models.RatingInfo, result models.GameResult) (*models.RatingInfo, *models.RatingInfo) {
	player1 := glicko2.NewPlayer(glicko2.NewRating(a.Value, a.Deviation, a.Volatility))
	player2 := glicko2.NewPlayer(glicko2.NewRating(b.Value, b.Deviation, b.Volatility))

	period := glicko2.NewRatingPeriod()
	period.AddMatch(player1, player2, convertModelMatchResult(result))

	period.Calculate()

	updateTime := time.Now().UTC()
	return &models.RatingInfo{
			UserId:     a.UserId,
			Value:      player1.Rating().R(),
			Deviation:  player1.Rating().Rd(),
			Volatility: player1.Rating().Sigma(),
			LastUpdate: updateTime,
		}, &models.RatingInfo{
			UserId:     b.UserId,
			Value:      player2.Rating().R(),
			Deviation:  player2.Rating().Rd(),
			Volatility: player2.Rating().Sigma(),
			LastUpdate: updateTime,
		}
}
