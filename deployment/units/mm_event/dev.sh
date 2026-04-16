cd $ROOT/deployment/units/mm_event
docker build --build-arg PROJECT_DIR=services/mm_event -t mm-event-server:latest -f Dockerfile $ROOT
minikube image load mm-event-server:latest
kubectl apply -f env.yaml
kubectl apply -f pod.yaml
kubectl apply -f service.yaml