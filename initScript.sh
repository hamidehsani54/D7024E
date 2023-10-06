#!/bin/bash

# Remove the stack
echo "Removing kadstack..."
sudo docker stack rm kadstack
# Compile modified files
echo "Building Docker image..."
sudo docker build -t kadlab:latest .

# Run the new code
echo "Deploying kadstack..."
sudo docker stack deploy -c docker-compose.yml kadstack