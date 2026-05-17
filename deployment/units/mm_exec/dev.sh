cd $ROOT/deployment/units/mm_exec
docker build --no-cache --build-arg PROJECT_DIR=services/mm_exec -t mm-exec:latest -f Dockerfile $ROOT
docker tag mm-exec:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/mm-exec:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/mm-exec:latest
kubectl apply -f auth_sidecar_rolebinding.yaml
kubectl apply -f env.yaml
kubectl apply -f deployment.yaml
