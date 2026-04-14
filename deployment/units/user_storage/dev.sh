cd $ROOT/deployment/units/user_storage
kubectl apply -f pvc.yaml
kubectl apply -f init_sql.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
