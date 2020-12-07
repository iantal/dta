#!/bin/bash

kubectl apply -f dta-configmap.yml
kubectl apply -f dta-service.yml
kubectl apply -f dta-deployment.yml