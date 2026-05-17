cd $ROOT/deployment/units/user_service
docker build --build-arg PROJECT_DIR=services/user_service -t user-service-server:latest -f Dockerfile $ROOT
docker tag user-service-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/user-service-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/user-service-server:latest
kubectl apply -f auth_sidecar_rolebinding.yaml
kubectl apply -f env.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml