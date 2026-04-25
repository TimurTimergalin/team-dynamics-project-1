cd $ROOT/deployment/units/fleet
docker build -t tag-duels-linux-server:latest -f Dockerfile .
docker tag tag-duels-linux-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/tag-duels-linux-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/tag-duels-linux-server:latest
kubectl apply -f fleet.yaml
kubectl apply -f fleet_autoscaler.yaml
