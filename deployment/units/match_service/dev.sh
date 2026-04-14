cd $ROOT/deployment/units/match_service
docker build --build-arg PROJECT_DIR=services/match_history_service -t match-service-server:latest -f Dockerfile $ROOT
minikube image load match-service-server:latest
kubectl apply -f env.yaml
kubectl apply -f pod.yaml
kubectl apply -f service.yaml