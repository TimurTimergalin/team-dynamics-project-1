cd $ROOT/deployment/config_maps
kubectl apply -f rating_storage.yaml
kubectl apply -f match_history_storage.yaml
kubectl apply -f user_storage_config.yaml
