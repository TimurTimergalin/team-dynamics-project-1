cd $ROOT/deployment/jobs/rotate_auth_keys
docker build -t rotate-auth-keys:latest -f Dockerfile $ROOT
docker tag rotate-auth-keys:latest 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/rotate-auth-keys:latest
docker push 4274eb52-inventive-macaw.registry.twcstorage.ru/tag-duels/rotate-auth-keys:latest
kubectl apply -f serviceaccount.yaml
kubectl apply -f role.yaml
kubectl apply -f rolebinding.yaml
kubectl apply -f cronjob.yaml
