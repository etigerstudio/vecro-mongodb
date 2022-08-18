#!/usr/bin/env bash

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o vecro-mongodb . &&
# docker build -t vecro-mongodb:v1 .
minikube -p l1 image build -t vecro-mongodb:v1 .