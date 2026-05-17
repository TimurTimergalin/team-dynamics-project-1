package auth_sdk

import (
	"context"
	"google.golang.org/grpc/metadata"
)

const (
	headerToken          = "x-auth-token"
	headerServiceAccount = "x-service-account"
)

type IncomingAuth struct {
	Token          string
	ServiceAccount *string // nil if this is a user token
}

func AttachServiceAccountToContext(ctx context.Context, token string, serviceAccount string) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		headerToken, token,
		headerServiceAccount, serviceAccount,
	)
}

func AttachCredentialsFromSidecar(ctx context.Context, client *AuthSidecarClient) (context.Context, error) {
	info, err := client.GetServiceAccount(ctx)
	if err != nil {
		return nil, err
	}
	return AttachServiceAccountToContext(ctx, info.Token, info.ServiceAccount), nil
}

func ExtractAuthFromContext(ctx context.Context) *IncomingAuth {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	tokens := md.Get(headerToken)
	if len(tokens) == 0 {
		return nil
	}
	auth := &IncomingAuth{Token: tokens[0]}
	if saValues := md.Get(headerServiceAccount); len(saValues) > 0 {
		sa := saValues[0]
		auth.ServiceAccount = &sa
	}
	return auth
}
