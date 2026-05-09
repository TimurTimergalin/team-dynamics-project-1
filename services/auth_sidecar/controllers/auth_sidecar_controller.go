package controllers

import (
	"context"
	pb "team_dynamics/api/proto/auth_sidecar"
	"team_dynamics/auth_sidecar/services"
)

type AuthSidecarController struct {
	pb.UnimplementedAuthSidecarServer
	Service services.AuthSidecarService
}

func (c *AuthSidecarController) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {
	return c.Service.Authorize(ctx, req)
}
