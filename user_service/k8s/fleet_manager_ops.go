package k8s

import (
	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	allocv1 "agones.dev/agones/pkg/apis/allocation/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	"context"
	"fleet_manager/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Allocate(client *versioned.Clientset, fmConfig *config.FleetManagerConfig, annotations map[string]string, ctx context.Context) (*allocv1.GameServerAllocationStatus, error) {
	timeoutCtx, cancelTimeoutCtx := context.WithTimeout(ctx, fmConfig.Timeout)
	defer cancelTimeoutCtx()

	allocation := &allocv1.GameServerAllocation{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: fmConfig.Namespace,
		},
		Spec: allocv1.GameServerAllocationSpec{
			Selectors: []allocv1.GameServerSelector{
				{
					LabelSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							"agones.dev/fleet": fmConfig.FleetName,
						},
					},
				},
			},
			Scheduling: fmConfig.Strategy,
			MetaPatch: allocv1.MetaPatch{
				Annotations: annotations,
			},
		},
	}

	result, err := client.AllocationV1().GameServerAllocations(fmConfig.Namespace).Create(timeoutCtx, allocation, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	} else {
		return &result.Status, nil
	}
}

func GetServer(client *versioned.Clientset, fmConfig *config.FleetManagerConfig, name string, ctx context.Context) (*agonesv1.GameServer, error) {
	timeoutCtx, cancelTimeoutCtx := context.WithTimeout(ctx, fmConfig.Timeout)
	defer cancelTimeoutCtx()

	return client.AgonesV1().GameServers(fmConfig.Namespace).Get(timeoutCtx, name, metav1.GetOptions{})
}
