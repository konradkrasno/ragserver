#!/bin/bash

docker build -f Dockerfile -t ragserver .
docker tag ragserver localhost:32000/ragserver
docker push localhost:32000/ragserver
