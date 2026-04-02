package services

import (
	"math"
	"time"

	"team_dynamics/rating_service/models" // Replace with actual import path
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

func (s *glicko2ServiceImpl) UpdateRatings(a, b *models.RatingInfo, result models.GameResult) (*models.RatingInfo, *models.RatingInfo) {
	var scoreA float64
	switch result {
	case models.Winner:
		scoreA = 1.0
	case models.Draw:
		scoreA = 0.5
	case models.Loser:
		scoreA = 0.0
	default:
		scoreA = 0.0
	}
	scoreB := 1.0 - scoreA

	now := time.Now().UTC()

	newA := s.updateSingle(a, b, scoreA)
	newB := s.updateSingle(b, a, scoreB)

	newA.UserId = a.UserId
	newA.LastUpdate = now
	newB.UserId = b.UserId
	newB.LastUpdate = now

	return newA, newB
}

func (s *glicko2ServiceImpl) updateSingle(player, opponent *models.RatingInfo, score float64) *models.RatingInfo {
	gOpp := g(opponent.Deviation)

	E := expected(player.Value, opponent.Value, gOpp)

	v := 1.0 / (math.Pow(gOpp, 2) * E * (1.0 - E))

	delta := v * gOpp * (score - E)

	newVolatility := s.updateVolatility(player.Volatility, delta, v)

	phiStar := math.Sqrt(math.Pow(player.Deviation, 2) + math.Pow(newVolatility, 2))

	newRating := player.Value + math.Pow(phiStar, 2)*delta
	newDeviation := 1.0 / math.Sqrt(1.0/math.Pow(phiStar, 2)+1.0/v)

	if newDeviation > 350.0 {
		newDeviation = 350.0
	}

	return &models.RatingInfo{
		Value:      newRating,
		Deviation:  newDeviation,
		Volatility: newVolatility,
	}
}

func g(deviation float64) float64 {
	return 1.0 / math.Sqrt(1.0+3.0*math.Pow(deviation, 2)/math.Pow(math.Pi, 2))
}

func expected(rating, oppRating, gOpp float64) float64 {
	return 1.0 / (1.0 + math.Exp(-gOpp*(rating-oppRating)))
}

func (s *glicko2ServiceImpl) updateVolatility(volatility, delta, v float64) float64 {
	a := math.Log(math.Pow(volatility, 2))
	B := 0.0

	f := func(x float64) float64 {
		expX := math.Exp(x)
		num := expX * (math.Pow(delta, 2) - math.Pow(volatility, 2) - v - expX)
		den := 2.0 * math.Pow(math.Pow(volatility, 2)+v+expX, 2)
		return num/den - (x-a)/math.Pow(s.Tau, 2)
	}

	if math.Pow(delta, 2) > math.Pow(volatility, 2)+v {
		B = math.Log(math.Pow(delta, 2) - math.Pow(volatility, 2) - v)
	} else {
		k := 1.0
		for f(a-k) < 0 {
			k += 1.0
		}
		B = a - k
	}

	for math.Abs(B-a) > s.Epsilon {
		C := a + (B-a)/2.0
		fC := f(C)
		if fC*f(B) <= 0 {
			a = C
		} else {
			B = C
		}
	}

	newVolatility := math.Exp(a / 2.0)
	if newVolatility > 0.1 {
		newVolatility = 0.1
	}
	return newVolatility
}
