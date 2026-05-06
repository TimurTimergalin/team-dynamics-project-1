source base.sh
kubectl apply -f dev_rating_storage_credentials_secret.yaml
kubectl apply -f dev_match_history_storage_credentials_secret.yaml
kubectl apply -f dev_match_kv_credentials_secret.yaml
kubectl apply -f dev_mm_pool_credentials.yaml
kubectl apply -f dev_user_storage_credentials.yaml
kubectl apply -f dev_user_kv_credentials.yaml
