package services

import (
	"context"
	auth "team_dynamics/auth_sdk"
	pb "team_dynamics/api/proto/auth_sidecar"
	fleetPb "team_dynamics/api/proto/fleet_manager"
)

const roleMatchService = "match-service"

var matchServiceAuthorityMap = []*pb.AuthorityMapEntry{
	{
		Key:   &pb.AuthorityMapEntry_ServiceAccount{ServiceAccount: "system:serviceaccount:default:match-service-sa"},
		Roles: []string{roleMatchService},
	},
}

type AuthService interface {
	AuthorizeAllocate(ctx context.Context, req *fleetPb.AllocateRequest) bool
	AuthorizeGetServer(ctx context.Context, req *fleetPb.GetServerRequest) bool
}

type authServiceImpl struct {
	client *auth.AuthSidecarClient
}

func MakeAuthService(client *auth.AuthSidecarClient) AuthService {
	return &authServiceImpl{client}
}

func (s *authServiceImpl) authorize(ctx context.Context) bool {
	incoming := auth.ExtractAuthFromContext(ctx)
	if incoming == nil || incoming.ServiceAccount == nil {
		return false
	}
	result, err := s.client.AuthorizeService(ctx, incoming.Token, matchServiceAuthorityMap)
	if err != nil {
		return false
	}
	return result.HasRole(roleMatchService)
}

func (s *authServiceImpl) AuthorizeAllocate(ctx context.Context, _ *fleetPb.AllocateRequest) bool {
	return s.authorize(ctx)
}

func (s *authServiceImpl) AuthorizeGetServer(ctx context.Context, _ *fleetPb.GetServerRequest) bool {
	return s.authorize(ctx)
}
