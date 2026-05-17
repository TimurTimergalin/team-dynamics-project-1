package services

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "team_dynamics/api/proto/auth_service"
	pbCommon "team_dynamics/api/proto/user_common"
	userPb "team_dynamics/api/proto/user_service"
	"team_dynamics/auth_service/downstream"
	"team_dynamics/auth_service/models"
	"team_dynamics/auth_service/repos"
	"team_dynamics/logging"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	jwtService            JwtService
	steamService          SteamService
	eosService            EosService
	userServiceClient     downstream.UserServiceClientFactory
	authKvRepo            repos.AuthKvRepo
	primaryPublicKeyPem   string
	secondaryPublicKeyPem string
	version               string
}

func marshalPublicKeyPem(keyPair models.KeyPair) (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(keyPair.PublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})), nil
}

func NewAuthService(jwtService JwtService, steamService SteamService, eosService EosService, userServiceClient downstream.UserServiceClientFactory, authKvRepo repos.AuthKvRepo, primaryKeyPair models.KeyPair, secondaryKeyPair models.KeyPair, version string) (*AuthService, error) {
	primaryPem, err := marshalPublicKeyPem(primaryKeyPair)
	if err != nil {
		return nil, err
	}
	secondaryPem, err := marshalPublicKeyPem(secondaryKeyPair)
	if err != nil {
		return nil, err
	}
	return &AuthService{
		jwtService:            jwtService,
		steamService:          steamService,
		eosService:            eosService,
		userServiceClient:     userServiceClient,
		authKvRepo:            authKvRepo,
		primaryPublicKeyPem:   primaryPem,
		secondaryPublicKeyPem: secondaryPem,
		version:               version,
	}, nil
}

func (s *AuthService) AuthExternal(ctx context.Context, req *pb.AuthExternalRequest) (*pb.AuthExternalResponse, error) {
	logger := logging.GetLogger(ctx)
	if req.ExternalKey == nil || req.AuthToken == nil {
		logger.Debug("invalid request: external_key and auth_token are required")
		return nil, status.Error(codes.InvalidArgument, "external_key and auth_token are required")
	}

	switch key := req.ExternalKey.Key.(type) {
	case *pbCommon.ExternalKey_SteamId:
		logger.Debug("AuthExternal steam started", "steam_id", key.SteamId)
		steamId, err := s.steamService.Validate(ctx, *req.AuthToken)
		if err != nil {
			logger.Error("steam validation failed", "steam_id", key.SteamId, "error", err)
			return nil, status.Errorf(codes.Unauthenticated, "steam validation failed: %v", err)
		}
		if steamId != key.SteamId {
			logger.Warn("steam id mismatch", "expected", key.SteamId, "actual", steamId)
			return nil, status.Error(codes.Unauthenticated, "steam id mismatch")
		}
		userResp, err := s.userServiceClient.GetSelfData(ctx, &userPb.GetSelfDataRequest{
			Key: &pbCommon.ExternalKey{Key: &pbCommon.ExternalKey_SteamId{SteamId: steamId}},
		})
		if err != nil {
			logger.Error("failed to get user data", "steam_id", steamId, "error", err)
			return nil, status.Errorf(codes.Internal, "failed to get user data: %v", err)
		}
		if userResp == nil || userResp.UserData == nil || userResp.UserData.Id == nil {
			logger.Error("user service returned invalid response", "steam_id", steamId)
			return nil, status.Error(codes.Internal, "user service returned invalid response")
		}
		userId := models.UserId{PlayerId: userResp.UserData.Id}
		access, refresh, err := s.jwtService.MakeTokenPair(userId)
		if err != nil {
			logger.Error("failed to create token pair", "user_id", *userResp.UserData.Id, "error", err)
			return nil, status.Errorf(codes.Internal, "failed to create token pair: %v", err)
		}
		logger.Debug("AuthExternal steam: ok", "user_id", *userResp.UserData.Id, "steam_id", steamId)
		return &pb.AuthExternalResponse{
			Access:      &pb.TokenWithExp{Token: &access.Token, ExpMs: &access.ExpMs},
			Refresh:     &pb.TokenWithExp{Token: &refresh.Token, ExpMs: &refresh.ExpMs},
			UserData:    userResp.UserData,
			ExternalKey: req.ExternalKey,
		}, nil
	case *pbCommon.ExternalKey_EosId:
		logger.Debug("AuthExternal eos started", "eos_id", key.EosId)
		accountId, err := s.eosService.Validate(ctx, *req.AuthToken)
		if err != nil {
			logger.Error("eos validation failed", "eos_id", key.EosId, "error", err)
			return nil, status.Errorf(codes.Unauthenticated, "eos validation failed: %v", err)
		}
		if accountId != key.EosId {
			logger.Warn("eos id mismatch", "expected", key.EosId, "actual", accountId)
			return nil, status.Error(codes.Unauthenticated, "eos id mismatch")
		}
		userResp, err := s.userServiceClient.GetSelfData(ctx, &userPb.GetSelfDataRequest{
			Key: &pbCommon.ExternalKey{Key: &pbCommon.ExternalKey_EosId{EosId: accountId}},
		})
		if err != nil {
			logger.Error("failed to get user data", "eos_id", accountId, "error", err)
			return nil, status.Errorf(codes.Internal, "failed to get user data: %v", err)
		}
		if userResp == nil || userResp.UserData == nil || userResp.UserData.Id == nil {
			logger.Error("user service returned invalid response", "eos_id", accountId)
			return nil, status.Error(codes.Internal, "user service returned invalid response")
		}
		userId := models.UserId{PlayerId: userResp.UserData.Id}
		access, refresh, err := s.jwtService.MakeTokenPair(userId)
		if err != nil {
			logger.Error("failed to create token pair", "user_id", *userResp.UserData.Id, "error", err)
			return nil, status.Errorf(codes.Internal, "failed to create token pair: %v", err)
		}
		logger.Debug("AuthExternal eos: ok", "user_id", *userResp.UserData.Id, "eos_id", accountId)
		return &pb.AuthExternalResponse{
			Access:      &pb.TokenWithExp{Token: &access.Token, ExpMs: &access.ExpMs},
			Refresh:     &pb.TokenWithExp{Token: &refresh.Token, ExpMs: &refresh.ExpMs},
			UserData:    userResp.UserData,
			ExternalKey: req.ExternalKey,
		}, nil
	default:
		logger.Debug("unsupported external key type")
		return nil, status.Error(codes.InvalidArgument, "unsupported external key type")
	}
}

func (s *AuthService) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	logger := logging.GetLogger(ctx)
	if req.Refresh == nil {
		logger.Debug("invalid request: refresh token is required")
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	payload, err := s.jwtService.Validate(*req.Refresh)
	if err != nil {
		logger.Error("invalid refresh token", "error", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token: %v", err)
	}

	if payload.TokenId == nil {
		logger.Warn("token is not a refresh token")
		return nil, status.Error(codes.Unauthenticated, "not a refresh token")
	}

	used, err := s.authKvRepo.MarkUsed(ctx, *payload.TokenId)
	if err != nil {
		logger.Error("failed to check token", "token_id", *payload.TokenId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to check token: %v", err)
	}
	if !used {
		logger.Warn("refresh token already used", "token_id", *payload.TokenId)
		return nil, status.Error(codes.Unauthenticated, "refresh token already used")
	}

	if payload.UserId.PlayerId == nil {
		logger.Warn("refresh token has no user id", "token_id", *payload.TokenId)
		return nil, status.Error(codes.Unauthenticated, "token has no user id")
	}

	access, refresh, err := s.jwtService.MakeTokenPair(payload.UserId)
	if err != nil {
		logger.Error("failed to create token pair", "user_id", *payload.UserId.PlayerId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create token pair: %v", err)
	}
	logger.Debug("Refresh: ok", "user_id", *payload.UserId.PlayerId, "token_id", *payload.TokenId)
	return &pb.RefreshResponse{
		Access:  &pb.TokenWithExp{Token: &access.Token, ExpMs: &access.ExpMs},
		Refresh: &pb.TokenWithExp{Token: &refresh.Token, ExpMs: &refresh.ExpMs},
	}, nil
}

func (s *AuthService) GetPublicKey(ctx context.Context, _ *pb.GetPublicKeyRequest) (*pb.GetPublicKeyResponse, error) {
	logger := logging.GetLogger(ctx)
	logger.Debug("GetPublicKey: ok", "version", s.version)
	return &pb.GetPublicKeyResponse{
		PrimaryPublicKey:   &s.primaryPublicKeyPem,
		SecondaryPublicKey: &s.secondaryPublicKeyPem,
		Version:            &s.version,
	}, nil
}

var _ pb.AuthServiceServer = (*AuthService)(nil)
