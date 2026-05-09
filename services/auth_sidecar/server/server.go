package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	pb "team_dynamics/api/proto/auth_sidecar"
	"team_dynamics/auth_sidecar/controllers"
	"team_dynamics/auth_sidecar/downstream"
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

	authClient := downstream.NewAuthServiceClient(authServiceAddress)
	jwtService := services.NewJwtService(issuer)
	roleService := services.NewRoleService()
	sidecarService := services.NewAuthSidecarService(jwtService, roleService, authClient)

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
