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
	"k8s.io/apimachinery/pkg/runtime/schema"
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

func makeAllocationName(matchId string) string {
	return fmt.Sprintf("game-server-allocation-of-%s", matchId)
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

func withTimeout[T any](ctx context.Context, cfg *OpsConfig, op func(context.Context) (T, *OpsError), defaultValue T) (T, *OpsError) {
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
			Name:      makeAllocationName(matchId),
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
	opErr := o.clearAllocation(ctx, matchId)
	if opErr != nil && opErr.Type != NotFoundError {
		return "", opErr
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

func (o *opsImpl) clearAllocation(ctx context.Context, matchId string) *OpsError {
	tCtx, cancel := context.WithTimeout(ctx, o.config.WriteTimeout)
	defer cancel()
	name := makeAllocationName(matchId)
	err := o.dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "allocation.agones.dev",
		Version:  "v1",
		Resource: "gameserverallocations",
	}).Namespace(o.config.Namespace).Delete(tCtx, name, metav1.DeleteOptions{})

	if err == nil {
		return nil
	}
	if k8sErrors.IsNotFound(err) {
		return nil
	}
	if k8sErrors.IsTimeout(err) || k8sErrors.IsServerTimeout(err) {
		return ConnectionError.Wrap(err)
	}
	return UnknownError.Wrap(err)
}

func (o *opsImpl) getServerByMatchIdOnce(ctx context.Context, matchId string) (string, *OpsError) {
	tCtx, cancel := context.WithTimeout(ctx, o.config.ReadTimeout)
	defer cancel()
	servers, err := o.client.AgonesV1().GameServers(o.config.Namespace).List(
		tCtx,
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", matchIdLabel, matchId),
		},
	)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return "", NotFoundError.Wrap(err)
		}
		if k8sErrors.IsTimeout(err) || k8sErrors.IsServerTimeout(err) {
			return "", ConnectionError.Wrap(err)
		}
		return "", UnknownError.Wrap(err)
	}
	for _, server := range servers.Items {
		return makeAddress(server.Status.Address, server.Status.Ports), nil
	}
	return "", NotFoundError.MakeF("list of servers was empty")
}

func (o *opsImpl) Allocate(ctx context.Context, matchId string, fleet *string, annotations map[string]string) (string, *OpsError) {
	return withTimeout(ctx, o.config, func(ctx context.Context) (string, *OpsError) {
		return o.allocateOnce(ctx, matchId, fleet, annotations)
	}, "")
}

func (o *opsImpl) GetServerByMatchId(ctx context.Context, matchId string) (string, *OpsError) {
	return withTimeout(ctx, o.config, func(ctx context.Context) (string, *OpsError) {
		return o.getServerByMatchIdOnce(ctx, matchId)
	}, "")
}
