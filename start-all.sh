#!/bin/bash
docker-compose -f ./django/docker-compose.yaml \
               -f ./golang/docker-compose.yaml \
               -f ./nextjs/docker-compose.yaml \
               -f ./nginx/docker-compose.yaml \
               -f ./postgres/docker-compose.yaml \
               -f ./rabbitmq/docker-compose.yaml \
               up --build
