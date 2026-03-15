package config

import (
	"agones.dev/agones/pkg/apis"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type FleetManagerConfig struct {
	Namespace     string
	FleetName     string
	Strategy      apis.SchedulingStrategy
	Timeout       time.Duration
	ListenAddress string
}

func GetConfig() (*FleetManagerConfig, error) {
	var res FleetManagerConfig

	res.Namespace = os.Getenv("CLUSTER_NAMESPACE")
	if res.Namespace == "" {
		return nil, errors.New("CLUSTER_NAMESPACE not set")
	}

	res.FleetName = os.Getenv("FLEET_NAME")
	if res.FleetName == "" {
		return nil, errors.New("FLEET_NAME not set")
	}

	if strategy, exists := map[string]apis.SchedulingStrategy{
		"packed":      apis.Packed,
		"distributed": apis.Distributed,
	}[strings.ToLower(os.Getenv("STRATEGY"))]; exists {
		res.Strategy = strategy
	} else {
		return nil, errors.New("STRATEGY is invalid")
	}

	if timeout, err := strconv.ParseUint(os.Getenv("K8S_TIMEOUT_MILLIS"), 10, 64); err == nil {
		res.Timeout = time.Duration(timeout) * time.Millisecond
	} else {
		return nil, err
	}

	res.ListenAddress = os.Getenv("LISTEN_ON")
	if res.ListenAddress == "" {
		return nil, errors.New("LISTEN_ON not set")
	}

	return &res, nil
}
