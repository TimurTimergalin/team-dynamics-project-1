package k8s

import (
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"fmt"
	"log"
)

func GetAddress(host string, ports []v1.GameServerStatusPort) *string {
	log.Printf("Trying to get address, host: %s, ports: %v", host, ports)
	for _, port := range ports {
		switch port.Name {
		case "game":
			res := fmt.Sprintf("%s:%d", host, port.Port)
			return &res
		}
	}
	return nil
}
