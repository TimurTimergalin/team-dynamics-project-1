cd $ROOT/deployment/units/rating_service/simple_testing_client
docker build --build-arg PROJECT_DIR=services/rating_service -t rating-service-simple-test-client:latest -f Dockerfile $ROOT
docker tag rating-service-simple-test-client:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/rating-service-simple-test-client:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/rating-service-simple-test-client:latest
kubectl apply -f pod.yaml
