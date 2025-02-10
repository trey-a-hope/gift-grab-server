#!/bin/bash
# to run - ./docker-cleanup.sh

# Stop all containers and remove volumes
echo "Stopping containers and removing volumes..."
docker compose down -v

# Remove all images
echo "Removing all Docker images..."
docker rmi $(docker images -q)

# Remove any remaining containers
echo "Removing any remaining containers..."
docker rm $(docker ps -a -q)

# Pull latest images
echo "Pulling latest images..."
docker compose pull

# Start containers
echo "Starting containers..."
docker compose up -d

echo "Docker cleanup and restart complete!"