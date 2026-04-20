cd $ROOT/deployment/units/mm_event
docker build --build-arg PROJECT_DIR=services/mm_event -t mm-event-server:latest -f Dockerfile $ROOT
docker tag mm-event-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/mm-event-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/mm-event-server:latest
kubectl apply -f env.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml