#!/bin/bash

docker build -f Dockerfile -t ragserver:local .
docker tag ragserver:local localhost:32000/ragserver:registry
docker push localhost:32000/ragserver
