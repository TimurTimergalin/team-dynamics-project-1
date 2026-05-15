package auth_sdk

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "team_dynamics/api/proto/auth_sidecar"
)

type ServiceAccountInfo struct {
	Token          string
	ServiceAccount string
}

type AuthorizeResult struct {
	Roles map[string]struct{}
}

// HasRole checks if the specified role exists in the authorization result
func (r *AuthorizeResult) HasRole(role string) bool {
	_, exists := r.Roles[role]
	return exists
}

type AuthSidecarClient struct {
	address string
}

func NewAuthSidecarClient(address string) *AuthSidecarClient {
	return &AuthSidecarClient{address: address}
}

func (c *AuthSidecarClient) dial() (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(c.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth sidecar: %w", err)
	}
	return conn, nil
}

func (c *AuthSidecarClient) GetServiceAccount(ctx context.Context) (*ServiceAccountInfo, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	resp, err := pb.NewAuthSidecarClient(conn).GetServiceAccount(ctx, &pb.GetServiceAccountRequest{})
	if err != nil {
		return nil, fmt.Errorf("GetServiceAccount failed: %w", err)
	}
	if resp.Token == nil || resp.ServiceAccount == nil {
		return nil, fmt.Errorf("incomplete response from auth sidecar")
	}
	return &ServiceAccountInfo{
		Token:          *resp.Token,
		ServiceAccount: *resp.ServiceAccount,
	}, nil
}

func (c *AuthSidecarClient) AuthorizeService(ctx context.Context, token string, authorityMap []*pb.AuthorityMapEntry) (*AuthorizeResult, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	resp, err := pb.NewAuthSidecarClient(conn).Authorize(ctx, &pb.AuthorizeRequest{
		Token:        &token,
		TokenType:    pb.TokenType_TOKEN_TYPE_SERVICE_ACCOUNT,
		AuthorityMap: authorityMap,
	})
	if err != nil {
		return nil, fmt.Errorf("AuthorizeService failed: %w", err)
	}
	return toAuthorizeResult(resp), nil
}

func (c *AuthSidecarClient) AuthorizeUser(ctx context.Context, token string, authorityMap []*pb.AuthorityMapEntry) (*AuthorizeResult, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	resp, err := pb.NewAuthSidecarClient(conn).Authorize(ctx, &pb.AuthorizeRequest{
		Token:        &token,
		TokenType:    pb.TokenType_TOKEN_TYPE_USER,
		AuthorityMap: authorityMap,
	})
	if err != nil {
		return nil, fmt.Errorf("AuthorizeUser failed: %w", err)
	}
	return toAuthorizeResult(resp), nil
}

func toAuthorizeResult(resp *pb.AuthorizeResponse) *AuthorizeResult {
	rolesSet := make(map[string]struct{})
	for _, role := range resp.Roles {
		rolesSet[role] = struct{}{}
	}
	return &AuthorizeResult{Roles: rolesSet}
}
