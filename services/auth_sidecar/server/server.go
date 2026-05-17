package main

import (
	"fmt"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net"
	"os"
	pb "team_dynamics/api/proto/auth_sidecar"
	"team_dynamics/auth_sdk"
	"team_dynamics/auth_sidecar/controllers"
	"team_dynamics/auth_sidecar/downstream"
	"team_dynamics/auth_sidecar/k8s"
	"team_dynamics/auth_sidecar/services"
)

func requireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("%s environment variable not set", key)
	}
	return val, nil
}

func main() {
	listenAddress, err := requireEnv("LISTEN_ADDRESS")
	if err != nil {
		log.Fatalf("%v", err)
	}
	authServiceAddress, err := requireEnv("AUTH_SERVICE_ADDRESS")
	if err != nil {
		log.Fatalf("%v", err)
	}
	issuer, err := requireEnv("JWT_ISSUER")
	if err != nil {
		log.Fatalf("%v", err)
	}
	tokenPath, err := requireEnv("TOKEN_PATH")
	if err != nil {
		log.Fatalf("%v", err)
	}
	namespace, err := requireEnv("NAMESPACE")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Create Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("failed to create in-cluster config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create kubernetes client: %v", err)
	}

	// Create k8s operations
	k8sOps := k8s.MakeOps(clientset, &k8s.OpsConfig{
		TokenPath: tokenPath,
		Namespace: namespace,
	})

	authSidecar := auth_sdk.NewAuthSidecarClient(listenAddress)
	authClient := downstream.NewAuthServiceClient(authServiceAddress, authSidecar)
	jwtService := services.NewJwtService(issuer)
	roleService := services.NewRoleService()
	sidecarService := services.NewAuthSidecarService(jwtService, roleService, authClient, k8sOps)

	controller := &controllers.AuthSidecarController{Service: sidecarService}

	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("cannot listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthSidecarServer(s, controller)

	log.Printf("auth-sidecar listening on %s", listenAddress)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("cannot serve: %v", err)
	}
}
