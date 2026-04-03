cd $ROOT/deployment/units/rating_service
docker build --build-arg PROJECT_DIR=services/rating_service -t rating-service-server:latest -f Dockerfile $ROOT
kubectl apply -f env.yaml
kubectl apply -f pod.yaml
kubectl apply -f service.yaml