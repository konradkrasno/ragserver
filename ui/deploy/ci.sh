#!/bin/bash

docker build -f Dockerfile -t ragserver-frontend .
docker tag ragserver-frontend localhost:32000/ragserver-frontend
docker push localhost:32000/ragserver-frontend
