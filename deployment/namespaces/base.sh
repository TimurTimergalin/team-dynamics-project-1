cd $ROOT/deployment/namespaces
kubectl apply -f rating.yaml
kubectl apply -f match_history.yaml
kubectl apply -f fleet.yaml
kubectl apply -f match.yaml
kubectl apply -f matchmaking.yaml
kubectl apply -f user.yaml
