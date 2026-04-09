package k8s

import (
	"agones.dev/agones/pkg/apis"
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	allocv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	"context"
	"fmt"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"time"
)

const matchIdLabel = "match-id"

type OpsConfig struct {
	WriteTimeout          time.Duration
	ReadTimeout           time.Duration
	Namespace             string
	CommonFleetLabelKey   string
	CommonFleetLabelValue string
	AllocStrategy         apis.SchedulingStrategy
	ConnectionRetries     int
	ContentionRetries     int
	UnknownRetries        int
}

func ptr[T any](v T) *T {
	return &v
}

type Ops interface {
	Allocate(ctx context.Context, matchId string, fleet *string, annotations map[string]string) (string, *OpsError)
	GetServerByMatchId(ctx context.Context, matchId string) (string, *OpsError)
}

type opsImpl struct {
	client        *versioned.Clientset
	dynamicClient dynamic.Interface
	config        *OpsConfig
}

func MakeOps(client *versioned.Clientset, dynamicClient dynamic.Interface, config *OpsConfig) Ops {
	return &opsImpl{client, dynamicClient, config}
}

func makeAddress(address string, ports []v1.GameServerStatusPort) string {
	for _, port := range ports {
		switch port.Name {
		case "default":
			return fmt.Sprintf("%s:%d", address, port.Port)
		}
	}
	return fmt.Sprintf("%s:%d", address, ports[0].Port)
}

func makeFleetLabelSelector(commonFleetLabelKey, commonFleetLabelValue string, fleet *string) map[string]string {
	res := map[string]string{
		commonFleetLabelKey: commonFleetLabelValue,
	}
	if fleet != nil {
		res["agones.dev/fleet"] = *fleet
	}
	return res
}

func withRetry[T any](ctx context.Context, cfg *OpsConfig, op func(context.Context) (T, *OpsError), defaultValue T) (T, *OpsError) {
	connectionRetries := cfg.ConnectionRetries
	contentionRetries := cfg.ContentionRetries
	unknownRetries := cfg.UnknownRetries
	for {
		select {
		case <-ctx.Done():
			return defaultValue, ConnectionError.MakeF("Context cancelled")
		default:
		}
		res, err := op(ctx)
		if err == nil {
			return res, nil
		}
		switch err.Type {
		case ConnectionError:
			connectionRetries -= 1
			if connectionRetries < 0 {
				return res, err
			}
		case ContentionError:
			contentionRetries -= 1
			if contentionRetries < 0 {
				return res, err
			}
		case UnknownError:
			unknownRetries -= 1
			if unknownRetries < 0 {
				return res, err
			}
		default:
			return res, err
		}
	}

}

func (o *opsImpl) allocateOnce(ctx context.Context, matchId string, fleet *string, annotations map[string]string) (string, *OpsError) {
	tCtx, cancel := context.WithTimeout(ctx, o.config.WriteTimeout)
	defer cancel()

	allocation := &allocv1.GameServerAllocation{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: o.config.Namespace,
		},
		Spec: allocv1.GameServerAllocationSpec{
			Selectors: []allocv1.GameServerSelector{
				{
					LabelSelector: metav1.LabelSelector{
						MatchLabels: makeFleetLabelSelector(o.config.CommonFleetLabelKey, o.config.CommonFleetLabelValue, fleet),
					},
				},
			},
			Scheduling: o.config.AllocStrategy,
			MetaPatch: allocv1.MetaPatch{
				Annotations: annotations,
				Labels: map[string]string{
					matchIdLabel: matchId,
				},
			},
		},
	}

	result, err := o.client.AllocationV1().GameServerAllocations(o.config.Namespace).Create(tCtx, allocation, metav1.CreateOptions{})
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return "", UniquenessViolationError.Make()
		}
		if k8sErrors.IsTimeout(err) || k8sErrors.IsServerTimeout(err) {
			return "", ConnectionError.Wrap(err)
		}
		return "", UnknownError.Wrap(err)
	}

	if result.Status.State == allocv1.GameServerAllocationAllocated {
		return makeAddress(result.Status.Address, result.Status.Ports), nil
	}

	switch result.Status.State {
	case allocv1.GameServerAllocationContention:
		return "", ContentionError.Make()
	case allocv1.GameServerAllocationUnAllocated:
		fallthrough
	default:
		return "", FleetFullError.Make()
	}
}

func (o *opsImpl) getServerByMatchIdOnce(ctx context.Context, matchId string) (string, *OpsError) {
	tCtx, cancel := context.WithTimeout(ctx, o.config.ReadTimeout)
	defer cancel()

	allocation := &allocv1.GameServerAllocation{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: o.config.Namespace,
		},
		Spec: allocv1.GameServerAllocationSpec{
			Selectors: []allocv1.GameServerSelector{
				{
					LabelSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							matchIdLabel: matchId,
						},
					},
					GameServerState: ptr(v1.GameServerStateAllocated),
				},
			},
		},
	}
	result, err := o.client.AllocationV1().GameServerAllocations(o.config.Namespace).Create(tCtx, allocation, metav1.CreateOptions{})
	if err != nil {
		if k8sErrors.IsTimeout(err) || k8sErrors.IsServerTimeout(err) {
			return "", ConnectionError.Wrap(err)
		}
		return "", UnknownError.Wrap(err)
	}
	switch result.Status.State {
	case allocv1.GameServerAllocationContention:
		return "", ContentionError.MakeF("contention while getting server")
	case allocv1.GameServerAllocationUnAllocated:
		return "", NotFoundError.MakeF("couldn't find game server")
	default:
	}
	return makeAddress(result.Status.Address, result.Status.Ports), nil
}

func (o *opsImpl) Allocate(ctx context.Context, matchId string, fleet *string, annotations map[string]string) (string, *OpsError) {
	return withRetry(ctx, o.config, func(ctx context.Context) (string, *OpsError) {
		return o.allocateOnce(ctx, matchId, fleet, annotations)
	}, "")
}

func (o *opsImpl) GetServerByMatchId(ctx context.Context, matchId string) (string, *OpsError) {
	return withRetry(ctx, o.config, func(ctx context.Context) (string, *OpsError) {
		return o.getServerByMatchIdOnce(ctx, matchId)
	}, "")
}
