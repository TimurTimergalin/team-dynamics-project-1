package downstream

import (
	"context"
	"team_dynamics/auth_sdk"
)

func attachCredentials(ctx context.Context, authSidecar *auth_sdk.AuthSidecarClient) (context.Context, error) {
	return auth_sdk.AttachCredentialsFromSidecar(ctx, authSidecar)
}
