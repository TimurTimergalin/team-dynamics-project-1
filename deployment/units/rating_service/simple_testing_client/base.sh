cd $ROOT/deployment/units/rating_service/simple_testing_client
docker build --build-arg PROJECT_DIR=services/rating_service -t rating-service-simple-test-client:latest:latest -f Dockerfile $ROOT
minikube image load rating-service-server:latest
kubectl apply -f pod.yaml
