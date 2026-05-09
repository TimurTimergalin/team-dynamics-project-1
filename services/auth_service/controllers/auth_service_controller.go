package controllers

import (
	"context"
	pb "team_dynamics/api/proto/auth_service"
	"team_dynamics/auth_service/services"
	"team_dynamics/grpc_lib"
)

type AuthServiceController struct {
	pb.UnimplementedAuthServiceServer
	Service *services.AuthService
}

func (s *AuthServiceController) AuthExternal(ctx context.Context, req *pb.AuthExternalRequest) (resp *pb.AuthExternalResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "AuthExternal")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.AuthExternal(ctx, req)
}

func (s *AuthServiceController) Refresh(ctx context.Context, req *pb.RefreshRequest) (resp *pb.RefreshResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "Refresh")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.Refresh(ctx, req)
}

func (s *AuthServiceController) GetPublicKey(ctx context.Context, req *pb.GetPublicKeyRequest) (resp *pb.GetPublicKeyResponse, err error) {
	ctx = grpc_lib.WithLogger(ctx, "GetPublicKey")
	defer grpc_lib.HandlePanic(ctx, &err)
	return s.Service.GetPublicKey(ctx, req)
}
