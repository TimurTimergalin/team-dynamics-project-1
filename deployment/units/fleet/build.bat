docker build -t tag-duels-linux-server:latest -f Dockerfile .
docker tag tag-duels-linux-server:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/tag-duels-linux-server:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/tag-duels-linux-server:latest