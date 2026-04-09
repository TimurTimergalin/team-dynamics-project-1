cd $ROOT/deployment/units/match_history_service
docker build --build-arg PROJECT_DIR=services/fleet_manager -t fleet-manager-server:latest -f Dockerfile $ROOT
minikube image load fleet-manager-server:latest
kubectl apply -f env.yaml
kubectl apply -f role.yaml
kubectl apply -f serviceaccount.yaml
kubectl apply -f rolebinding.yaml
kubectl apply -f pod.yaml
kubectl apply -f service.yaml