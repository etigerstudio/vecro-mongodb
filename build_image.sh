#!/usr/bin/env bash

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ben-mongodb . &&
# docker build -t ben-mongodb:v1 .
minikube -p l1 image build -t ben-mongodb:v1 .