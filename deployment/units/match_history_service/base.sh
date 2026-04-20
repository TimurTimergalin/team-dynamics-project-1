cd $ROOT/deployment/units/match_history_service
docker build --build-arg PROJECT_DIR=services/match_history_service -t match-history-service-server:latest -f Dockerfile $ROOT
docker tag match-history-service-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/match-history-service-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/match-history-service-server:latest
kubectl apply -f env.yaml
kubectl apply -f pod.yaml
kubectl apply -f service.yaml