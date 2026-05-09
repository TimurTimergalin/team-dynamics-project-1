package services

import (
	"context"
	"errors"
	authPb "team_dynamics/api/proto/auth_service"
	pbCommon "team_dynamics/api/proto/user_common"
	pb "team_dynamics/api/proto/auth_sidecar"
	"team_dynamics/auth_sidecar/downstream"
	"team_dynamics/auth_sidecar/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthSidecarService interface {
	Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error)
}

type authSidecarServiceImpl struct {
	jwtService    JwtService
	roleService   RoleService
	authClient    downstream.AuthServiceClient
	cachedVersion string
}

func NewAuthSidecarService(jwtService JwtService, roleService RoleService, authClient downstream.AuthServiceClient) AuthSidecarService {
	return &authSidecarServiceImpl{
		jwtService:  jwtService,
		roleService: roleService,
		authClient:  authClient,
	}
}

func (s *authSidecarServiceImpl) fetchAndUpdateKeys(ctx context.Context) (bool, error) {
	resp, err := s.authClient.GetPublicKey(ctx, &authPb.GetPublicKeyRequest{})
	if err != nil {
		return false, err
	}
	if resp.Version == nil || *resp.Version == s.cachedVersion {
		return false, nil
	}
	if err := s.jwtService.UpdateKeys(resp); err != nil {
		return false, err
	}
	s.cachedVersion = *resp.Version
	return true, nil
}

func protoToAuthorityMap(entries []*pb.AuthorityMapEntry) models.AuthorityMap {
	am := make(models.AuthorityMap)
	for _, entry := range entries {
		var key models.UserId
		switch k := entry.Key.(type) {
		case *pb.AuthorityMapEntry_AnyUser:
			key = models.UserId{}
		case *pb.AuthorityMapEntry_ExternalId:
			if steamKey, ok := k.ExternalId.Key.(*pbCommon.ExternalKey_SteamId); ok {
				steamId := steamKey.SteamId
				key = models.UserId{SteamId: &steamId}
			}
		case *pb.AuthorityMapEntry_UserId:
			playerId := int64(k.UserId)
			key = models.UserId{PlayerId: &playerId}
		}
		am[key] = append(am[key], entry.Roles...)
	}
	return am
}

func (s *authSidecarServiceImpl) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	if req.Token == nil {
		return nil, status.Error(codes.Unauthenticated, "token is required")
	}

	userId, err := s.jwtService.Validate(*req.Token)
	if err != nil {
		if !errors.Is(err, ErrKeysOutdated) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}
		updated, fetchErr := s.fetchAndUpdateKeys(ctx)
		if fetchErr != nil {
			return nil, status.Errorf(codes.Internal, "failed to fetch keys: %v", fetchErr)
		}
		if !updated {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		userId, err = s.jwtService.Validate(*req.Token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}
	}

	return &pb.AuthorizeResponse{
		Roles: s.roleService.ResolveRoles(protoToAuthorityMap(req.AuthorityMap), *userId),
	}, nil
}
