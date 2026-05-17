cd $ROOT/deployment/units/user_event
docker build --build-arg PROJECT_DIR=services/user_event -t user-event-server:latest -f Dockerfile $ROOT
docker tag user-event-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/user-event-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/user-event-server:latest
kubectl apply -f auth_sidecar_rolebinding.yaml
kubectl apply -f env.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
