package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"team_dynamics/auth_service/models"
	"time"
)

var ErrNoUserId = errors.New("no player_id set")

type JwtService interface {
	MakeTokenPair(userId models.UserId) (models.TokenWithExp, models.TokenWithExp, error)
	Validate(token string) (*models.TokenPayload, error)
}

type jwtServiceImpl struct {
	AccessTokenExpirationTime  time.Duration
	RefreshTokenExpirationTime time.Duration
	Issuer                     string
	PrimaryKeyPair             models.KeyPair
	SecondaryKeyPair           models.KeyPair
}

func NewJwtService(
	accessExpiration time.Duration,
	refreshExpiration time.Duration,
	issuer string,
	primaryKeyPair models.KeyPair,
	secondaryKeyPair models.KeyPair,
) JwtService {
	return &jwtServiceImpl{
		AccessTokenExpirationTime:  accessExpiration,
		RefreshTokenExpirationTime: refreshExpiration,
		Issuer:                     issuer,
		PrimaryKeyPair:             primaryKeyPair,
		SecondaryKeyPair:           secondaryKeyPair,
	}
}

func (j *jwtServiceImpl) makeToken(claims jwt.MapClaims, keyPair models.KeyPair) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(keyPair.PrivateKey)
}

func (j *jwtServiceImpl) MakeTokenPair(userId models.UserId) (models.TokenWithExp, models.TokenWithExp, error) {
	now := time.Now()

	baseClaims := jwt.MapClaims{
		"iss": j.Issuer,
		"iat": now.Unix(),
	}
	if userId.PlayerId == nil {
		return models.TokenWithExp{}, models.TokenWithExp{}, ErrNoUserId
	}
	baseClaims["player_id"] = *userId.PlayerId

	accessExp := now.Add(j.AccessTokenExpirationTime)
	accessClaims := jwt.MapClaims{}
	for k, v := range baseClaims {
		accessClaims[k] = v
	}
	accessClaims["exp"] = accessExp.Unix()

	refreshExp := now.Add(j.RefreshTokenExpirationTime)
	refreshClaims := jwt.MapClaims{}
	for k, v := range baseClaims {
		refreshClaims[k] = v
	}
	refreshClaims["exp"] = refreshExp.Unix()
	refreshClaims["token_id"] = uuid.New().String()

	accessToken, err := j.makeToken(accessClaims, j.PrimaryKeyPair)
	if err != nil {
		return models.TokenWithExp{}, models.TokenWithExp{}, err
	}
	refreshToken, err := j.makeToken(refreshClaims, j.PrimaryKeyPair)
	if err != nil {
		return models.TokenWithExp{}, models.TokenWithExp{}, err
	}
	return models.TokenWithExp{Token: accessToken, ExpMs: accessExp.UnixMilli()},
		models.TokenWithExp{Token: refreshToken, ExpMs: refreshExp.UnixMilli()},
		nil
}

func (j *jwtServiceImpl) Validate(token string) (*models.TokenPayload, error) {
	var parseErr error
	var parsed *jwt.Token
	for _, keyPair := range []models.KeyPair{j.PrimaryKeyPair, j.SecondaryKeyPair} {
		parsed, parseErr = jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return keyPair.PublicKey, nil
		}, jwt.WithIssuer(j.Issuer), jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
		if parseErr == nil {
			break
		}
	}
	if parseErr != nil {
		return nil, parseErr
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	payload := &models.TokenPayload{}

	if v, ok := claims["player_id"]; ok {
		if f, ok := v.(float64); ok {
			playerId := int64(f)
			payload.UserId.PlayerId = &playerId
		}
	}
	if v, ok := claims["token_id"]; ok {
		if s, ok := v.(string); ok {
			payload.TokenId = &s
		}
	}

	return payload, nil
}
