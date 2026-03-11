#!/bin/bash

docker build -f SimpleTestClient.Dockerfile -t simple-test-client:latest .
minikube image load simple-test-client:latest:latest
kubectl apply -f simple-test-client-pod.yaml