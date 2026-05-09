cd $ROOT/deployment/units/auth_service
docker build --build-arg PROJECT_DIR=services/auth_service -t auth-service-server:latest -f Dockerfile $ROOT
docker tag auth-service-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/auth-service-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/auth-service-server:latest
kubectl apply -f env.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
