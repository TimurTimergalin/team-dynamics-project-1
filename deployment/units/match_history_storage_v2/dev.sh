kubectl apply -f ../../../deployment/config_maps/match_history_storage_v2.yaml
kubectl apply -f init_sql.yaml
kubectl apply -f dev_persistent_volume_claim.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
