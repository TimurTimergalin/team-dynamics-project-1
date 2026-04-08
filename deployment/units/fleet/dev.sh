cd $ROOT/deployment/units/fleet
kubectl apply -f fleet.yaml
kubectl apply -f fleet_autoscaler.yaml
