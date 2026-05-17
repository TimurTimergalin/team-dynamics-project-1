cd $ROOT/deployment/sidecars/auth_sidecar
docker build --build-arg PROJECT_DIR=services/auth_sidecar -t auth-sidecar-server:latest -f Dockerfile $ROOT
docker tag auth-sidecar-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/auth-sidecar-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/auth-sidecar-server:latest
kubectl apply -f env.yaml
kubectl apply -f role.yaml
