cd $ROOT/deployment/units/match_history_service
docker build --build-arg PROJECT_DIR=services/match_history_service -t match-history-service-server:latest -f Dockerfile $ROOT
minikube image load match-history-service-server:latest
kubectl apply -f env.yaml
kubectl apply -f pod.yaml
kubectl apply -f service.yaml