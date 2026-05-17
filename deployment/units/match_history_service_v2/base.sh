cd $ROOT/deployment/units/match_history_service_v2
docker build --build-arg PROJECT_DIR=services/match_history_service_v2 -t match-history-service-v2-server:latest -f Dockerfile $ROOT
docker tag match-history-service-v2-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/match-history-service-v2-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/match-history-service-v2-server:latest
kubectl apply -f serviceaccount.yaml
kubectl apply -f auth_sidecar_rolebinding.yaml
kubectl apply -f env.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
