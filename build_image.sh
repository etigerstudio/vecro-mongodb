#!/usr/bin/env bash

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ben-base . &&
docker build -t ben-base:v1 .