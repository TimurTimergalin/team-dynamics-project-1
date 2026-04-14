cd $ROOT/deployment/units/mm_exec
docker build --build-arg PROJECT_DIR=services/match_history_service -t mm-exec:latest -f Dockerfile $ROOT
minikube image load  mm-exec:latest
kubectl apply -f env.yaml
kubectl apply -f pod.yaml
