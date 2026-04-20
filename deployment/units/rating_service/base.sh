cd $ROOT/deployment/units/rating_service
docker build --build-arg PROJECT_DIR=services/rating_service -t rating-service-server:latest -f Dockerfile $ROOT
docker tag rating-service-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/rating-service-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/rating-service-server:latest
kubectl apply -f env.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml