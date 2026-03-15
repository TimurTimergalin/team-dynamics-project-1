package auth_sidecar

import (
	"context"
	"crypto"
	"io/ioutil"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "path/to/your/generated/proto" // replace with actual import
)

// AuthSidecarService implements the gRPC service.
type AuthSidecarService struct {
	JwksFetcher      JwksFetcher      // for Kubernetes service account tokens
	PublicKeyFetcher PublicKeyFetcher // for user tokens (staggered keys)
}

// JwksFetcher abstracts fetching a public key by key ID (kid) from the cluster's JWKS.
type JwksFetcher interface {
	GetKeyFromJwks(kid string) crypto.PublicKey
}

// PublicKeyFetcher abstracts fetching all currently valid user‑auth public keys.
type PublicKeyFetcher interface {
	GetPublicKeys() []crypto.PublicKey
}

// GetToken returns the pod's own Kubernetes service account token.
func (s *AuthSidecarService) GetToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	tokenBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read token: %v", err)
	}
	token := string(tokenBytes)
	return &pb.GetTokenResponse{Token: &token}, nil
}

// Authorize validates a JWT and returns the role associated with its identity,
// if present in the authority map.
func (s *AuthSidecarService) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	// If no token provided, return empty response (no role).
	if req.Token == nil || *req.Token == "" {
		return &pb.AuthorizeResponse{}, nil
	}
	tokenStr := *req.Token

	// 1. Parse token without validation to inspect issuer and header.
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	unverified, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse token: %v", err)
	}
	claims, ok := unverified.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid claims format")
	}

	issuer, ok := claims["iss"].(string)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing 'iss' claim")
	}

	var validatedClaims jwt.MapClaims

	// 2. Validate the token with the appropriate key source.
	switch issuer {
	case "kubernetes/serviceaccount": // Kubernetes service account token
		kid, ok := unverified.Header["kid"].(string)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing 'kid' in token header")
		}
		key := s.JwksFetcher.GetKeyFromJwks(kid)
		if key == nil {
			return nil, status.Errorf(codes.Unauthenticated, "key not found for kid: %s", kid)
		}
		validatedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return key, nil
		})
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid service account token: %v", err)
		}
		validatedClaims, ok = validatedToken.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Errorf(codes.Internal, "invalid claims after validation")
		}

	case "TagDuels/UserService": // User token issued by your auth service
		keys := s.PublicKeyFetcher.GetPublicKeys()
		var validatedToken *jwt.Token
		for _, key := range keys {
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return key, nil
			})
			if err == nil {
				validatedToken = token
				break
			}
		}
		if validatedToken == nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid user token")
		}
		validatedClaims, ok = validatedToken.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Errorf(codes.Internal, "invalid claims after validation")
		}

	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown issuer: %s", issuer)
	}

	// 3. Extract identity from the validated claims.
	var matchedRole *string
	for _, entry := range req.AuthorityMap {
		switch key := entry.Key.(type) {
		case *pb.AuthorityMapEntry_ServiceAccount:
			// Service account token path
			if issuer != "kubernetes/serviceaccount" {
				continue
			}
			sub, ok := validatedClaims["sub"].(string)
			if !ok {
				return nil, status.Errorf(codes.InvalidArgument, "missing 'sub' in token")
			}
			// sub format: system:serviceaccount:<namespace>:<name>
			parts := strings.Split(sub, ":")
			if len(parts) != 4 {
				return nil, status.Errorf(codes.InvalidArgument, "malformed 'sub' claim")
			}
			saName := parts[3]
			if saName == key.ServiceAccount {
				matchedRole = entry.Role // entry.Role is a *string (proto optional)
				break
			}

		case *pb.AuthorityMapEntry_UserId:
			// User token path
			if issuer != "TagDuels/UserService" {
				continue
			}
			userIDClaim, ok := validatedClaims["user_id"]
			if !ok {
				return nil, status.Errorf(codes.InvalidArgument, "missing 'user_id' in token")
			}
			var userID uint64
			switch v := userIDClaim.(type) {
			case float64:
				userID = uint64(v)
			case int64:
				userID = uint64(v)
			case uint64:
				userID = v
			default:
				return nil, status.Errorf(codes.InvalidArgument, "invalid type for 'user_id'")
			}
			if userID == key.UserId {
				matchedRole = entry.Role
				break
			}
		}
	}

	// 4. Return the role (if any) in the response.
	resp := &pb.AuthorizeResponse{}
	if matchedRole != nil {
		resp.Role = matchedRole
	}
	return resp, nil
}
