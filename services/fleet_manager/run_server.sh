#!/bin/bash

docker build -f ServerDev.Dockerfile -t fleet-manager-server:latest .
minikube image load fleet-manager-server:latest
kubectl apply -f fleet-manager-serviceaccount.yaml
kubectl apply -f fleet-manager-role.yaml
kubectl apply -f fleet-manager-rolebinding-sa.yaml
kubectl apply -f fleet-manager-pod.yaml
kubectl apply -f fleet-manager-service.yaml