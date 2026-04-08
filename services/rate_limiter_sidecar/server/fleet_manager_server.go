package main

import (
	"agones.dev/agones/pkg/client/clientset/versioned"
	"fleet_manager/config"
	"fleet_manager/controllers"
	pb "fleet_manager/fleet_manager/proto"
	"fleet_manager/services"
	"google.golang.org/grpc"
	"k8s.io/client-go/rest"
	"log"
	"net"
)

func getController(fmConfig *config.FleetManagerConfig) (*controllers.FleetManagerController, error) {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	agonesClient, err := versioned.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}

	controller := controllers.FleetManagerController{
		UnimplementedFleetManagerServer: pb.UnimplementedFleetManagerServer{},
		Service: &services.FleetManagerService{
			K8sClient: agonesClient,
			FmConfig:  fmConfig,
		},
	}
	return &controller, nil
}

func main() {
	fmConfig, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error while getting fleet manager config, %s", err.Error())
	}

	controller, err := getController(fmConfig)
	if err != nil {
		log.Fatalf("Error while getting controllers, %s", err.Error())
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
