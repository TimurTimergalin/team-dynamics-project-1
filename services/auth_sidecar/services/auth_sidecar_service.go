package services

import (
	"context"
	"errors"
	authPb "team_dynamics/api/proto/auth_service"
	pbCommon "team_dynamics/api/proto/user_common"
	pb "team_dynamics/api/proto/auth_sidecar"
	"team_dynamics/auth_sidecar/downstream"
	"team_dynamics/auth_sidecar/k8s"
	"team_dynamics/auth_sidecar/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthSidecarService interface {
	Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error)
	GetServiceAccount(ctx context.Context, req *pb.GetServiceAccountRequest) (*pb.GetServiceAccountResponse, error)
}

type authSidecarServiceImpl struct {
	jwtService    JwtService
	roleService   RoleService
	authClient    downstream.AuthServiceClient
	k8sOps        k8s.Ops
	cachedVersion string
}

func NewAuthSidecarService(jwtService JwtService, roleService RoleService, authClient downstream.AuthServiceClient, k8sOps k8s.Ops) AuthSidecarService {
	return &authSidecarServiceImpl{
		jwtService:  jwtService,
		roleService: roleService,
		authClient:  authClient,
		k8sOps:      k8sOps,
	}
}

func (s *authSidecarServiceImpl) GetServiceAccount(ctx context.Context, _ *pb.GetServiceAccountRequest) (*pb.GetServiceAccountResponse, error) {
	info, err := s.k8sOps.GetServiceAccount(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get service account: %v", err)
	}
	token := info.Token
	sa := info.ServiceAccount
	return &pb.GetServiceAccountResponse{
		Token:          &token,
		ServiceAccount: &sa,
	}, nil
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
		var key models.Principal
		switch k := entry.Key.(type) {
		case *pb.AuthorityMapEntry_AnyUser:
			key = models.AnyUserPrincipal()
		case *pb.AuthorityMapEntry_ExternalId:
			if steamKey, ok := k.ExternalId.Key.(*pbCommon.ExternalKey_SteamId); ok {
				steamId := steamKey.SteamId
				key = models.UserPrincipal(models.UserId{SteamId: &steamId})
			}
		case *pb.AuthorityMapEntry_UserId:
			playerId := int64(k.UserId)
			key = models.UserPrincipal(models.UserId{PlayerId: &playerId})
		case *pb.AuthorityMapEntry_ServiceAccount:
			key = models.ServiceAccountPrincipal(k.ServiceAccount)
		}
		am[key] = append(am[key], entry.Roles...)
	}
	return am
}

func (s *authSidecarServiceImpl) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	if req.Token == nil {
		return nil, status.Error(codes.Unauthenticated, "token is required")
	}

	am := protoToAuthorityMap(req.AuthorityMap)

	if req.TokenType == pb.TokenType_TOKEN_TYPE_SERVICE_ACCOUNT {
		sa, err := s.k8sOps.ValidateServiceAccountToken(ctx, *req.Token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid service account token: %v", err)
		}
		return &pb.AuthorizeResponse{
			Roles: s.roleService.ResolveRolesForServiceAccount(am, sa),
		}, nil
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
		Roles: s.roleService.ResolveRolesForUser(am, *userId),
	}, nil
}
