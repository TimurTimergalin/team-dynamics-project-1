cd $ROOT/deployment/units/fleet_manager
docker build --build-arg PROJECT_DIR=services/fleet_manager -t fleet-manager-server:latest -f Dockerfile $ROOT
docker tag fleet-manager-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/fleet-manager-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/fleet-manager-server:latest
kubectl apply -f auth_sidecar_rolebinding.yaml
kubectl apply -f env.yaml
kubectl apply -f role.yaml
kubectl apply -f serviceaccount.yaml
kubectl apply -f rolebinding.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml