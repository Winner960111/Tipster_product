#!/bin/bash

echo "Enter parameter: $1"

if [ "$1" == "-p" ]; then
    echo "Shutdown mongodb ..."
    docker-compose -p "database-transaction" -f "./docker-compose.db.yml" -f "./docker-compose.yml" down -v
fi

echo "Starting mongodb ..."
docker-compose -p "database-transaction" -f "./docker-compose.db.yml" -f "./docker-compose.yml" up -d --build

# docker-compose up -d --build

