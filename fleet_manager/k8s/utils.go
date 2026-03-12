package k8s

import (
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"fmt"
)

func GetAddress(host string, ports []v1.GameServerStatusPort) *string {
	for _, port := range ports {
		switch port.Name {
		case "default":
			res := fmt.Sprintf("%s:%d", host, port.Port)
			return &res
		}
	}
	return nil
}
