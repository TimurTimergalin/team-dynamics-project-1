package main

import (
	"agones.dev/agones/pkg/apis"
	"agones.dev/agones/pkg/client/clientset/versioned"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"log"
	"net"
	"os"
	"strconv"
	pb "team_dynamics/api/proto/fleet_manager"
	"team_dynamics/fleet_manager/config"
	"team_dynamics/fleet_manager/controllers"
	"team_dynamics/fleet_manager/k8s"
	"team_dynamics/fleet_manager/services"
	"time"
)

func getFleetManagerConfig() (*config.FleetManagerConfig, error) {
	cfg := &config.FleetManagerConfig{}
	if v := os.Getenv("LISTEN_ADDRESS"); v != "" {
		cfg.ListenAddress = v
	} else {
		return nil, errors.New("LISTEN_ADDRESS environment variable not set")
	}
	return cfg, nil
}

func getOpsConfig() (*k8s.OpsConfig, error) {
	cfg := &k8s.OpsConfig{}

	// WriteTimeout
	if v := os.Getenv("OPS_WRITE_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid OPS_WRITE_TIMEOUT: %w", err)
		}
		cfg.WriteTimeout = d
	} else {
		return nil, errors.New("OPS_WRITE_TIMEOUT environment variable not set")
	}

	// ReadTimeout
	if v := os.Getenv("OPS_READ_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid OPS_READ_TIMEOUT: %w", err)
		}
		cfg.ReadTimeout = d
	} else {
		return nil, errors.New("OPS_READ_TIMEOUT environment variable not set")
	}

	// Namespace
	if v := os.Getenv("OPS_NAMESPACE"); v != "" {
		cfg.Namespace = v
	} else {
		return nil, errors.New("OPS_NAMESPACE environment variable not set")
	}

	// CommonFleetLabelKey
	if v := os.Getenv("OPS_COMMON_FLEET_LABEL_KEY"); v != "" {
		cfg.CommonFleetLabelKey = v
	} else {
		return nil, errors.New("OPS_COMMON_FLEET_LABEL_KEY environment variable not set")
	}

	// CommonFleetLabelValue
	if v := os.Getenv("OPS_COMMON_FLEET_LABEL_VALUE"); v != "" {
		cfg.CommonFleetLabelValue = v
	} else {
		return nil, errors.New("OPS_COMMON_FLEET_LABEL_VALUE environment variable not set")
	}

	// AllocStrategy
	if v := os.Getenv("OPS_ALLOC_STRATEGY"); v != "" {
		switch v {
		case "Packed":
			cfg.AllocStrategy = apis.Packed
		case "Distributed":
			cfg.AllocStrategy = apis.Distributed
		default:
			return nil, fmt.Errorf("invalid OPS_ALLOC_STRATEGY: %q, must be 'Packed' or 'Distributed'", v)
		}
	} else {
		return nil, errors.New("OPS_ALLOC_STRATEGY environment variable not set")
	}

	// ConnectionRetries (int)
	if v := os.Getenv("OPS_CONNECTION_RETRIES"); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid OPS_CONNECTION_RETRIES: %w", err)
		}
		cfg.ConnectionRetries = i
	} else {
		return nil, errors.New("OPS_CONNECTION_RETRIES environment variable not set")
	}

	// ContentionRetries (int)
	if v := os.Getenv("OPS_CONTENTION_RETRIES"); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid OPS_CONTENTION_RETRIES: %w", err)
		}
		cfg.ContentionRetries = i
	} else {
		return nil, errors.New("OPS_CONTENTION_RETRIES environment variable not set")
	}

	// UnknownRetries (int)
	if v := os.Getenv("OPS_UNKNOWN_RETRIES"); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid OPS_UNKNOWN_RETRIES: %w", err)
		}
		cfg.UnknownRetries = i
	} else {
		return nil, errors.New("OPS_UNKNOWN_RETRIES environment variable not set")
	}

	return cfg, nil
}

func main() {
	fmConfig, err := getFleetManagerConfig()
	if err != nil {
		log.Fatalf("Error while getting fleet manager config, %v", err)
	}
	opsConfig, err := getOpsConfig()
	if err != nil {
		log.Fatalf("Error while gettign k8s ops config: %v", err)
	}
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error while getting k8s client config: %v", err)
	}
	versionedClient, err := versioned.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatalf("Error while getting agones client: %v", versionedClient)
	}
	dynamicClient, err := dynamic.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatalf("Error while getting dynamic client: %v", err)
	}

	controller := &controllers.FleetManagerController{
		Service: services.MakeFleetManagerService(
			k8s.MakeOps(
				versionedClient,
				dynamicClient,
				opsConfig,
			),
		),
	}

	lis, err := net.Listen("tcp", fmConfig.ListenAddress)
	if err != nil {
		log.Fatalf("Error while trying to listen on address: %s", err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterFleetManagerServer(s, controller)

	log.Printf("Strating to listen on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while trying to serve: %s", err.Error())
	}
}
