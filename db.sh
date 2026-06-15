#!/bin/sh

podman run --rm -d --name cosmosim-db --network host -e ALLOW_EXTERNAL_CONNECTIONS=true -e CREATE_USER=true -e USERNAME=cosmosim -e PASSWORD=cosmosim -e LOG_LEVEL=info -e SKIP_INIT_DATA=true -v cosmosim-db:/data  ghcr.io/documentdb/documentdb/documentdb-local:latest

podman run --rm -d --name cosmosim-webui --net host -e ME_CONFIG_MONGODB_SSL=true -e ME_CONFIG_MONGODB_URL="mongodb://cosmosim:cosmosim@localhost:10260/cosmosim?authSource=authdb&authMechanism=SCRAM-SHA-256&tls=true&tlsInsecure=true" -e PORT=8087 -e ME_CONFIG_MONGODB_AWS_DOCUMENTDB=true -e ME_CONFIG_MONGODB_ENABLE_ADMIN=false docker.io/mongo-express
