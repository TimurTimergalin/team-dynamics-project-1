cd $ROOT/deployment/units/match_service
docker build --build-arg PROJECT_DIR=services/match_service -t match-service-server:latest -f Dockerfile $ROOT
docker tag match-service-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/match-service-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/match-service-server:latest
kubectl apply -f env.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml