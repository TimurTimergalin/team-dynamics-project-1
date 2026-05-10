package k8s

import (
	"context"
	"fmt"
	"os"
	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type OpsConfig struct {
	TokenPath string
	Namespace string
}

type ServiceAccountInfo struct {
	Token          string
	ServiceAccount string
}

type Ops interface {
	GetServiceAccount(ctx context.Context) (*ServiceAccountInfo, error)
	ValidateServiceAccountToken(ctx context.Context, token string) (string, error)
}

type opsImpl struct {
	client *kubernetes.Clientset
	config *OpsConfig
}

func MakeOps(client *kubernetes.Clientset, config *OpsConfig) Ops {
	return &opsImpl{client, config}
}

func (o *opsImpl) GetServiceAccount(ctx context.Context) (*ServiceAccountInfo, error) {
	tokenBytes, err := os.ReadFile(o.config.TokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account token: %w", err)
	}
	token := string(tokenBytes)

	sa, err := o.ValidateServiceAccountToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate own token: %w", err)
	}

	return &ServiceAccountInfo{
		Token:          token,
		ServiceAccount: sa,
	}, nil
}

func (o *opsImpl) ValidateServiceAccountToken(ctx context.Context, token string) (string, error) {
	review := &authv1.TokenReview{
		Spec: authv1.TokenReviewSpec{
			Token: token,
		},
	}

	result, err := o.client.AuthenticationV1().TokenReviews().Create(ctx, review, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("token review failed: %w", err)
	}
	if !result.Status.Authenticated {
		return "", fmt.Errorf("token is not authenticated: %s", result.Status.Error)
	}

	return result.Status.User.Username, nil
}
