package auth_sdk

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "team_dynamics/api/proto/auth_sidecar"
	pbCommon "team_dynamics/api/proto/user_common"
)

var ErrPrincipalMismatch = errors.New("token principal does not match expected principal")

type ServiceAccountInfo struct {
	Token          string
	ServiceAccount string
}

type AuthorizeResult struct {
	Roles     []string
	Principal isPrincipal
}

type isPrincipal interface{ isPrincipal() }

type UserIdPrincipal struct{ UserId uint64 }
type SteamIdPrincipal struct{ SteamId int64 }
type ServiceAccountPrincipal struct{ ServiceAccount string }

func (UserIdPrincipal) isPrincipal()         {}
func (SteamIdPrincipal) isPrincipal()        {}
func (ServiceAccountPrincipal) isPrincipal() {}

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

func (c *AuthSidecarClient) AuthorizeService(ctx context.Context, token string, expectedSA string, authorityMap []*pb.AuthorityMapEntry) (*AuthorizeResult, error) {
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
	result := toAuthorizeResult(resp)
	sa, ok := result.Principal.(ServiceAccountPrincipal)
	if !ok || sa.ServiceAccount != expectedSA {
		return nil, ErrPrincipalMismatch
	}
	return result, nil
}

func (c *AuthSidecarClient) AuthorizeUser(ctx context.Context, token string, expectedPrincipal isPrincipal, authorityMap []*pb.AuthorityMapEntry) (*AuthorizeResult, error) {
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
	result := toAuthorizeResult(resp)
	if result.Principal != expectedPrincipal {
		return nil, ErrPrincipalMismatch
	}
	return result, nil
}

func toAuthorizeResult(resp *pb.AuthorizeResponse) *AuthorizeResult {
	result := &AuthorizeResult{Roles: resp.Roles}
	switch p := resp.Principal.(type) {
	case *pb.AuthorizeResponse_UserId:
		result.Principal = UserIdPrincipal{UserId: p.UserId}
	case *pb.AuthorizeResponse_ExternalId:
		if p.ExternalId != nil {
			if steamKey, ok := p.ExternalId.Key.(*pbCommon.ExternalKey_SteamId); ok {
				result.Principal = SteamIdPrincipal{SteamId: steamKey.SteamId}
			}
		}
	case *pb.AuthorizeResponse_ServiceAccount:
		result.Principal = ServiceAccountPrincipal{ServiceAccount: p.ServiceAccount}
	}
	return result
}
