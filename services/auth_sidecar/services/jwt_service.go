package services

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	authPb "team_dynamics/api/proto/auth_service"
	"team_dynamics/auth_sidecar/models"
)

var ErrKeysOutdated = errors.New("token validation failed: keys may be outdated")

type JwtService interface {
	UpdateKeys(resp *authPb.GetPublicKeyResponse) error
	Validate(token string) (*models.UserId, error)
}

type jwtServiceImpl struct {
	mu      sync.RWMutex
	version string
	primary *rsa.PublicKey
	secondary *rsa.PublicKey
	issuer  string
}

func NewJwtService(issuer string) JwtService {
	return &jwtServiceImpl{issuer: issuer}
}

func parsePublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not RSA")
	}
	return rsaKey, nil
}

func (s *jwtServiceImpl) UpdateKeys(resp *authPb.GetPublicKeyResponse) error {
	if resp.PrimaryPublicKey == nil || resp.SecondaryPublicKey == nil || resp.Version == nil {
		return errors.New("incomplete key response")
	}
	primary, err := parsePublicKey(*resp.PrimaryPublicKey)
	if err != nil {
		return fmt.Errorf("primary key: %w", err)
	}
	secondary, err := parsePublicKey(*resp.SecondaryPublicKey)
	if err != nil {
		return fmt.Errorf("secondary key: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.primary = primary
	s.secondary = secondary
	s.version = *resp.Version
	return nil
}

func (s *jwtServiceImpl) tryValidate(token string, key *rsa.PublicKey) (*models.UserId, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	}, jwt.WithIssuer(s.issuer), jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, ErrKeysOutdated
		}
		return nil, err
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	userId := &models.UserId{}
	if v, ok := claims["player_id"]; ok {
		if f, ok := v.(float64); ok {
			playerId := int64(f)
			userId.PlayerId = &playerId
		}
	}
	return userId, nil
}

func (s *jwtServiceImpl) Validate(token string) (*models.UserId, error) {
	s.mu.RLock()
	primary := s.primary
	secondary := s.secondary
	s.mu.RUnlock()

	if primary == nil {
		return nil, ErrKeysOutdated
	}

	for _, key := range []*rsa.PublicKey{primary, secondary} {
		if key == nil {
			continue
		}
		userId, err := s.tryValidate(token, key)
		if err == nil {
			return userId, nil
		}
	}

	return nil, ErrKeysOutdated
}
