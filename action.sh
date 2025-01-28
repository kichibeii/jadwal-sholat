#!/bin/bash

# Define container and image names
CONTAINER_NAME="jadwal_sholat_container"
IMAGE_NAME="jadwal_sholat_service"

# Kill the container if it's running
echo "Stopping container if running..."
docker kill $CONTAINER_NAME 2>/dev/null || echo "Container not running."

# Remove the container if it exists
echo "Removing container if exists..."
docker rm $CONTAINER_NAME 2>/dev/null || echo "Container not found."

# Build the Docker image
echo "Building Docker image..."
docker build -t $IMAGE_NAME .

# Run the Docker container
echo "Running Docker container..."
docker run -d --name $CONTAINER_NAME $IMAGE_NAME

echo "Done!"
