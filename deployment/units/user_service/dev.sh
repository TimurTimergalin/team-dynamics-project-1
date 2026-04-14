cd $ROOT/deployment/units/user_service
docker build --build-arg PROJECT_DIR=services/user_service -t user-service-server:latest -f Dockerfile $ROOT
minikube image load user-service-server:latest
kubectl apply -f env.yaml
kubectl apply -f pod.yaml
kubectl apply -f service.yaml