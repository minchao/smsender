#!/bin/bash

docker build -t minchao/smsender-preview:latest .
echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
docker images
docker push "${DOCKER_USERNAME}/smsender-preview"
docker logout
