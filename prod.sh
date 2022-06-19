#!/bin/bash

env=${1:-env}
docker-compose pull
docker-compose --env-file $env -f docker-compose.yml up -d
